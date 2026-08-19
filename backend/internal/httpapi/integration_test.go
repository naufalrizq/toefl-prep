//go:build integration

package httpapi_test

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"io"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/santhosh-tekuri/jsonschema/v5"

	"toefl-prep/backend/internal/attempts"
	"toefl-prep/backend/internal/auth"
	"toefl-prep/backend/internal/database"
	"toefl-prep/backend/internal/exams"
	"toefl-prep/backend/internal/httpapi"
	"toefl-prep/backend/internal/questions"
	"toefl-prep/backend/internal/reporting"
	"toefl-prep/backend/internal/seed"
)

//go:embed schemas/*.json
var schemaFS embed.FS

var schemas = map[string]*jsonschema.Schema{}

func init() {
	dir, err := schemaFS.ReadDir("schemas")
	if err != nil {
		panic(err)
	}
	compiler := jsonschema.NewCompiler()
	for _, entry := range dir {
		data, err := schemaFS.ReadFile("schemas/" + entry.Name())
		if err != nil {
			panic(err)
		}
		if err := compiler.AddResource("https://toefl-prep/api/schemas/"+entry.Name(), bytes.NewReader(data)); err != nil {
			panic(err)
		}
	}
	for _, entry := range dir {
		schema, err := compiler.Compile("https://toefl-prep/api/schemas/" + entry.Name())
		if err != nil {
			panic(err)
		}
		schemas[strings.TrimSuffix(entry.Name(), ".json")] = schema
	}
}

func envOrSkip(t *testing.T) string {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set; run scripts/db-test-setup.sh first")
	}
	return url
}

// boot migrates, seeds, publishes one exam, and returns a live test server.
func boot(t *testing.T) (string, int64) {
	t.Helper()
	ctx := context.Background()
	if err := database.Migrate(ctx, envOrSkip(t)); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	pool, err := database.New(ctx, envOrSkip(t))
	if err != nil {
		t.Fatalf("database: %v", err)
	}
	t.Cleanup(pool.Close)
	if err := seed.EnsureUsers(ctx, pool); err != nil {
		t.Fatalf("seed users: %v", err)
	}

	authSvc := auth.New(pool, 24*time.Hour)
	qRepo := questions.NewRepo(pool)
	qSvc := questions.NewService(qRepo)
	examRepo := exams.NewRepo(pool)
	examSvc := exams.NewService(examRepo, qRepo)
	attemptRepo := attempts.NewRepo(pool)
	attemptSvc := attempts.NewService(attemptRepo, examRepo, qRepo)
	reportSvc := reporting.NewService(attemptRepo)

	qs, err := seed.Load()
	if err != nil {
		t.Fatalf("seed load: %v", err)
	}
	if n, err := qSvc.Seed(ctx, qs); err != nil || n == 0 {
		t.Fatalf("seed questions: n=%d err=%v", n, err)
	}
	exam := &exams.ExamTemplate{
		Title: "Structure Basics",
		SectionFilters: map[string]exams.SectionConfig{
			"structure": {Parts: []exams.PartConfig{{Title: "Part 1", Type: "sentence-completion", Count: 2}}},
			"vocabulary": {Parts: []exams.PartConfig{{Title: "Part 1", Type: "vocab-multiple-choice", Count: 2}}},
		},
		Shuffle:            true,
		Mode:               "both",
		SecondsPerQuestion: intp(20),
		TotalMinutes:       intp(15),
	}
	id, err := examSvc.Create(ctx, exam)
	if err != nil {
		t.Fatalf("create exam: %v", err)
	}
	if err := examSvc.Publish(ctx, id, true); err != nil {
		t.Fatalf("publish exam: %v", err)
	}

	r := httpapi.New(httpapi.Deps{
		Auth:      auth.NewHandler(authSvc, 10),
		Questions: questions.NewHandler(qSvc),
		Seed:      seed.NewHandler(qSvc),
		Exams:     exams.NewHandler(examSvc),
		Attempts:  attempts.NewHandler(attemptSvc),
		Reporting: reporting.NewHandler(reportSvc),
		CORS:      []string{},
	})
	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)
	return srv.URL, id
}

func intp(v int) *int { return &v }

func validate(t *testing.T, name string, body []byte) {
	t.Helper()
	var doc any
	if err := json.Unmarshal(body, &doc); err != nil {
		t.Fatalf("%s: invalid JSON: %v", name, err)
	}
	if err := schemas[name].Validate(doc); err != nil {
		t.Errorf("%s: schema violation: %v\nbody: %s", name, err, body)
	}
}

type apiClient struct {
	t      *testing.T
	base   *url.URL
	client *http.Client
}

func newAPIClient(t *testing.T, base string) *apiClient {
	t.Helper()
	u, err := url.Parse(base)
	if err != nil {
		t.Fatalf("parse base: %v", err)
	}
	jar, _ := cookiejar.New(nil)
	return &apiClient{t: t, base: u, client: &http.Client{Jar: jar}}
}

func (c *apiClient) csrf() string {
	for _, cookie := range c.client.Jar.Cookies(c.base) {
		if cookie.Name == "toefl_csrf" {
			return cookie.Value
		}
	}
	return ""
}

func (c *apiClient) do(method, path string, body any, token string) (*http.Response, []byte) {
	c.t.Helper()
	var rd io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rd = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, c.base.String()+path, rd)
	if err != nil {
		c.t.Fatalf("new request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("X-CSRF-Token", token)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		c.t.Fatalf("%s %s: %v", method, path, err)
	}
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return resp, data
}

func TestContract(t *testing.T) {
	base, examID := boot(t)
	c := newAPIClient(t, base)

	// Health (public)
	{
		resp, body := c.do("GET", "/api/v1/health", nil, "")
		if resp.StatusCode != 200 {
			t.Fatalf("health status = %d", resp.StatusCode)
		}
		validate(t, "health", body)
	}

	// Unauthenticated protected call -> auth_required envelope
	{
		resp, body := c.do("GET", "/api/v1/attempts", nil, "")
		if resp.StatusCode != 401 {
			t.Fatalf("attempts without auth = %d, want 401", resp.StatusCode)
		}
		validate(t, "error", body)
		if !strings.Contains(string(body), "auth_required") {
			t.Errorf("expected auth_required code, got %s", body)
		}
	}

	// Login as student
	{
		resp, body := c.do("POST", "/api/v1/auth/login", map[string]string{
			"email": "student", "password": "123",
		}, "")
		if resp.StatusCode != 200 {
			t.Fatalf("login status = %d body = %s", resp.StatusCode, body)
		}
		validate(t, "user", body)
	}

	// Wrong password -> 401
	{
		resp, _ := c.do("POST", "/api/v1/auth/login", map[string]string{
			"email": "student", "password": "wrong",
		}, "")
		if resp.StatusCode != 401 {
			t.Errorf("bad login status = %d, want 401", resp.StatusCode)
		}
	}

	// CSRF missing on mutation -> 403
	{
		resp, body := c.do("POST", "/api/v1/attempts", map[string]any{
			"exam_template_id": examID, "mode": "per_question",
		}, "")
		if resp.StatusCode != 403 {
			t.Fatalf("mutation without CSRF = %d, want 403 (body %s)", resp.StatusCode, body)
		}
	}

	token := c.csrf()
	if token == "" {
		t.Fatal("no csrf cookie after login")
	}

	// Start an attempt
	startResp, startBody := c.do("POST", "/api/v1/attempts", map[string]any{
		"exam_template_id": examID, "mode": "per_question",
	}, token)
	if startResp.StatusCode != 201 {
		t.Fatalf("start status = %d body = %s", startResp.StatusCode, startBody)
	}
	validate(t, "start-result", startBody)
	var start struct {
		Data struct {
			Attempt struct {
				ID int64 `json:"id"`
			} `json:"attempt"`
			Items []struct {
				ID int64 `json:"id"`
			} `json:"items"`
		} `json:"data"`
	}
	if err := json.Unmarshal(startBody, &start); err != nil {
		t.Fatalf("parse start: %v", err)
	}
	if len(start.Data.Items) != 4 {
		t.Errorf("expected 4 quiz items, got %d", len(start.Data.Items))
	}

	// Duplicate start -> 409
	{
		resp, body := c.do("POST", "/api/v1/attempts", map[string]any{
			"exam_template_id": examID, "mode": "per_question",
		}, token)
		if resp.StatusCode != 409 {
			t.Errorf("duplicate start = %d, want 409 (body %s)", resp.StatusCode, body)
		}
		validate(t, "error", body)
	}

	// Resume questions
	{
		resp, body := c.do("GET", fmt.Sprintf("/api/v1/attempts/%d/questions", start.Data.Attempt.ID), nil, "")
		if resp.StatusCode != 200 {
			t.Fatalf("questions status = %d", resp.StatusCode)
		}
		if !strings.Contains(string(body), "question_snapshot") {
			t.Errorf("questions payload missing snapshots: %s", body)
		}
	}

	// Answer one question correctly
	{
		resp, _ := c.do("PUT", fmt.Sprintf("/api/v1/attempts/%d/answers/%d", start.Data.Attempt.ID, start.Data.Items[0].ID),
			map[string]any{"chosen_index": 0, "time_taken_ms": 1500}, token)
		if resp.StatusCode != 204 {
			t.Errorf("answer status = %d", resp.StatusCode)
		}
	}

	// Submit and validate the report contract
	reportResp, reportBody := c.do("POST", fmt.Sprintf("/api/v1/attempts/%d/submit", start.Data.Attempt.ID), nil, token)
	if reportResp.StatusCode != 200 {
		t.Fatalf("submit status = %d body = %s", reportResp.StatusCode, reportBody)
	}
	validate(t, "report", reportBody)

	// Review
	{
		resp, body := c.do("GET", fmt.Sprintf("/api/v1/attempts/%d/review", start.Data.Attempt.ID), nil, "")
		if resp.StatusCode != 200 {
			t.Fatalf("review status = %d", resp.StatusCode)
		}
		validate(t, "report", body)
	}

	// Dashboard stats
	{
		resp, body := c.do("GET", "/api/v1/dashboard/stats", nil, "")
		if resp.StatusCode != 200 {
			t.Fatalf("dashboard status = %d", resp.StatusCode)
		}
		validate(t, "dashboard", body)
	}

	// Exams list (student sees published)
	{
		resp, body := c.do("GET", "/api/v1/exams", nil, "")
		if resp.StatusCode != 200 {
			t.Fatalf("exams status = %d", resp.StatusCode)
		}
		validate(t, "exam-list", body)
	}

	// Student cannot access admin endpoints -> 403
	{
		resp, body := c.do("GET", "/api/v1/questions", nil, "")
		if resp.StatusCode != 403 {
			t.Errorf("student questions access = %d, want 403 (body %s)", resp.StatusCode, body)
		}
		validate(t, "error", body)
	}
}

// TestAdminContract exercises the admin surface with the same envelope rules.
func TestAdminContract(t *testing.T) {
	base, _ := boot(t)
	c := newAPIClient(t, base)

	// Login as admin
	{
		resp, body := c.do("POST", "/api/v1/auth/login", map[string]string{
			"email": "admin", "password": "123",
		}, "")
		if resp.StatusCode != 200 {
			t.Fatalf("admin login = %d body = %s", resp.StatusCode, body)
		}
	}
	token := c.csrf()

	// Seed endpoint (idempotent upsert)
	{
		resp, body := c.do("POST", "/api/v1/seed", nil, token)
		if resp.StatusCode != 200 {
			t.Fatalf("seed status = %d body = %s", resp.StatusCode, body)
		}
	}

	// Question import: one valid + one invalid draft
	{
		resp, body := c.do("POST", "/api/v1/questions/import", []map[string]any{
			{
				"section": "structure", "type": "sentence-completion",
				"question_text": "The student _____ the exam every day.",
				"options":       []string{"takes", "take", "taking", "taken"},
				"correct_index": 0, "explanation": "Subjek tunggal, jadi takes.",
				"highlights":    map[string]any{"verb": []string{"takes"}},
			},
			{
				"section": "structure", "type": "sentence-completion",
				"question_text": "missing blank",
				"options":       []string{"a", "b", "c", "d"},
				"correct_index": 0, "explanation": "x",
			},
		}, token)
		if resp.StatusCode != 200 {
			t.Fatalf("import status = %d body = %s", resp.StatusCode, body)
		}
		var wrapped struct {
			Data []struct {
				Index int    `json:"index"`
				Valid bool   `json:"valid"`
				Error string `json:"error"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &wrapped); err != nil {
			t.Fatalf("parse import: %v", err)
		}
		results := wrapped.Data
		if !results[0].Valid || results[1].Valid {
			t.Errorf("import validation wrong: %+v", results)
		}
	}

	// Admin lists questions (paginated)
	{
		resp, body := c.do("GET", "/api/v1/questions", nil, "")
		if resp.StatusCode != 200 {
			t.Fatalf("questions list = %d", resp.StatusCode)
		}
		if !strings.Contains(string(body), "total") {
			t.Errorf("questions list missing pagination: %s", body)
		}
	}

	// Admin exams list
	{
		resp, body := c.do("GET", "/api/v1/exams", nil, "")
		if resp.StatusCode != 200 {
			t.Fatalf("admin exams = %d", resp.StatusCode)
		}
		validate(t, "exam-list", body)
	}

	// Logout clears the session
	{
		resp, _ := c.do("POST", "/api/v1/auth/logout", nil, token)
		if resp.StatusCode != 204 {
			t.Errorf("logout = %d, want 204", resp.StatusCode)
		}
		resp2, _ := c.do("GET", "/api/v1/dashboard/stats", nil, "")
		if resp2.StatusCode != 401 {
			t.Errorf("post-logout dashboard = %d, want 401", resp2.StatusCode)
		}
	}
}