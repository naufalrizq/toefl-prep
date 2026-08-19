package questions

import (
	"context"
	"strings"
	"testing"
)

func TestImportSingleDraft(t *testing.T) {
	svc := NewService(nil)
	body := `{
		"section": "structure",
		"type": "sentence-completion",
		"difficulty": "medium",
		"question_text": "She _____ to school. She goes every day.",
		"options": ["goes", "go", "going", "went"],
		"correct_index": 0,
		"explanation": "She is singular, jadi pakai goes.",
		"highlights": {"verb": ["goes"]}
	}`
	results, err := svc.Import(context.Background(), strings.NewReader(body))
	if err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if !results[0].Valid {
		t.Fatalf("draft should be valid: %s", results[0].Error)
	}
	q := results[0].Question
	if q.Section != "structure" || q.Type != "sentence-completion" {
		t.Errorf("defaults applied wrong: %+v", q)
	}
	if len(q.HighlightRegions) != 1 || q.HighlightRegions[0].Pos != "verb" {
		t.Errorf("highlights not normalized: %+v", q.HighlightRegions)
	}
}

func TestImportArrayDraft(t *testing.T) {
	svc := NewService(nil)
	body := `[
		{"section":"structure","type":"sentence-completion","question_text":"A _____.",
		 "options":["x","y","z","w"],"correct_index":1,"explanation":"e","difficulty":"easy"},
		{"section":"vocabulary","type":"vocab-multiple-choice","question_text":"B?",
		 "options":["1","2","3","4"],"correct_index":2,"explanation":"e","difficulty":"hard"}
	]`
	results, err := svc.Import(context.Background(), strings.NewReader(body))
	if err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Index != 0 || results[1].Index != 1 {
		t.Errorf("indexes wrong: %+v", results)
	}
	for i, r := range results {
		if !r.Valid {
			t.Errorf("result %d should be valid: %s", i, r.Error)
		}
	}
}

func TestImportInvalidItemsDoNotBlock(t *testing.T) {
	svc := NewService(nil)
	body := `[
		{"section":"structure","type":"sentence-completion","question_text":"Bad blank missing.",
		 "options":["a","b","c","d"],"correct_index":0,"explanation":"e"},
		{"section":"vocabulary","type":"vocab-multiple-choice","question_text":"Good?",
		 "options":["a","b","c","d"],"correct_index":0,"explanation":"e","difficulty":"easy"}
	]`
	results, err := svc.Import(context.Background(), strings.NewReader(body))
	if err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Valid {
		t.Error("first draft should be invalid (missing blank)")
	}
	if results[0].Question != nil {
		t.Error("invalid draft must not carry a Question")
	}
	if !results[1].Valid {
		t.Errorf("second draft should be valid: %s", results[1].Error)
	}
}

func TestImportNonJSON(t *testing.T) {
	svc := NewService(nil)
	_, err := svc.Import(context.Background(), strings.NewReader("definitely not json"))
	if err == nil {
		t.Fatal("Import() should reject non-JSON body")
	}
}

func TestImportDefaultsBlankOptionForCorrectIndex(t *testing.T) {
	svc := NewService(nil)
	body := `{
		"section": "structure",
		"question_text": "_____ is the answer.",
		"options": ["", "b", "c", "d"],
		"correct_index": 0,
		"explanation": "e"
	}`
	results, err := svc.Import(context.Background(), strings.NewReader(body))
	if err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	if !results[0].Valid {
		t.Fatalf("draft should be valid after blank fill: %s", results[0].Error)
	}
	if got := results[0].Question.Options[0]; got != "_____" {
		t.Errorf("option[0] = %q, want _____", got)
	}
}

func TestDraftDefaults(t *testing.T) {
	d := Draft{QuestionText: "X _____ y.", Options: []string{"a", "b", "c", "d"}}
	q := FromDraft(d)
	if q.Section != "structure" || q.Type != "sentence-completion" || q.Difficulty != "medium" {
		t.Errorf("defaults wrong: %+v", q)
	}
	if !q.Active {
		t.Error("drafted questions must be active by default")
	}
}