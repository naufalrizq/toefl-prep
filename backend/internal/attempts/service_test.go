package attempts

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"toefl-prep/backend/internal/exams"
	"toefl-prep/backend/internal/questions"
)

func i(v int) *int { return &v }

type fakeStore struct {
	attempts map[int64]*Attempt
	items    map[int64][]*AttemptItem
	active   map[string]*Attempt // "userID:templateID" -> attempt
	nextID   int64
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		attempts: map[int64]*Attempt{},
		items:    map[int64][]*AttemptItem{},
		active:   map[string]*Attempt{},
		nextID:   1,
	}
}

func (f *fakeStore) InsertAttempt(ctx context.Context, a *Attempt, items []*AttemptItem) (int64, error) {
	a.ID = f.nextID
	f.nextID++
	f.attempts[a.ID] = a
	f.items[a.ID] = items
	f.active[actKey(a.UserID, a.ExamTemplateID)] = a
	return a.ID, nil
}
func (f *fakeStore) GetByID(ctx context.Context, userID, id int64) (*Attempt, error) {
	a, ok := f.attempts[id]
	if !ok || a.UserID != userID {
		return nil, errNotFound()
	}
	return a, nil
}
func (f *fakeStore) ActiveByTemplate(ctx context.Context, userID, templateID int64) (*Attempt, error) {
	return f.active[actKey(userID, templateID)], nil
}
func (f *fakeStore) Items(ctx context.Context, attemptID int64) ([]*AttemptItem, error) {
	return f.items[attemptID], nil
}
func (f *fakeStore) SetAnswer(ctx context.Context, attemptID, itemID int64, chosen *int, timeMs *int) error {
	for _, it := range f.items[attemptID] {
		if it.ID == itemID {
			it.ChosenIndex = chosen
			it.TimeTakenMs = timeMs
			return nil
		}
	}
	return errNotFound()
}
func (f *fakeStore) SetFlagged(ctx context.Context, attemptID, itemID int64, flagged bool) error {
	for _, it := range f.items[attemptID] {
		if it.ID == itemID {
			it.Flagged = flagged
			return nil
		}
	}
	return errNotFound()
}
func (f *fakeStore) Finalize(ctx context.Context, a *Attempt) error {
	f.attempts[a.ID] = a
	return nil
}
func (f *fakeStore) List(ctx context.Context, userID int64, limit int) ([]*Attempt, error) {
	var out []*Attempt
	for _, a := range f.attempts {
		if a.UserID == userID {
			out = append(out, a)
		}
	}
	return out, nil
}

func actKey(userID, templateID int64) string {
	return strconv.FormatInt(userID, 10) + ":" + strconv.FormatInt(templateID, 10)
}

func errNotFound() error { return &errNotFoundT{} }

type errNotFoundT struct{}

func (*errNotFoundT) Error() string { return "not_found" }

type fakeExamSource struct {
	byID map[int64]*exams.ExamTemplate
}

func (f *fakeExamSource) GetByID(ctx context.Context, id int64) (*exams.ExamTemplate, error) {
	e, ok := f.byID[id]
	if !ok {
		return nil, errNotFound()
	}
	return e, nil
}

type fakeBank struct {
	qs map[string][]*questions.Question
}

func (f *fakeBank) AllActiveBySectionType(ctx context.Context, section, typ string) ([]*questions.Question, error) {
	return f.qs[section+"/"+typ], nil
}

func question(id int64, section string) *questions.Question {
	return &questions.Question{
		ID:               id,
		Section:          section,
		Type:             "sentence-completion",
		QuestionText:     "q" + strconv.FormatInt(id, 10),
		Options:          []string{"a", "b", "c", "d"},
		CorrectIndex:     0,
		Explanation:      "Penjelasan.",
		HighlightRegions: []questions.HighlightRegion{{Start: 0, End: 1, Pos: "noun"}},
	}
}

func publishedExam(id int64, mode string) *exams.ExamTemplate {
	return &exams.ExamTemplate{
		ID:              id,
		Title:           "Exam " + strconv.FormatInt(id, 10),
		SectionFilters: map[string]exams.SectionConfig{
			"structure": {Parts: []exams.PartConfig{
				{Title: "Part 1", Type: "sentence-completion", Count: 2},
			}},
		},
		Shuffle:            true,
		Mode:               mode,
		TotalMinutes:       i(1),
		SecondsPerQuestion: i(20),
		Published:          true,
		Active:             true,
	}
}

func serviceWith(store *fakeStore, examSource *fakeExamSource, bank *fakeBank) *Service {
	if store == nil {
		store = newFakeStore()
	}
	if examSource == nil {
		examSource = &fakeExamSource{byID: map[int64]*exams.ExamTemplate{}}
	}
	if bank == nil {
		bank = &fakeBank{qs: map[string][]*questions.Question{}}
	}
	return &Service{store: store, exams: examSource, bank: bank}
}

func TestStartSnapshotsQuestions(t *testing.T) {
	store := newFakeStore()
	bank := &fakeBank{qs: map[string][]*questions.Question{
		"structure/sentence-completion": {question(1, "structure"), question(2, "structure")},
	}}
	svc := serviceWith(store, &fakeExamSource{byID: map[int64]*exams.ExamTemplate{
		9: publishedExam(9, "per_question"),
	}}, bank)

	res, err := svc.Start(context.Background(), 5, 9, "per_question")
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if len(res.Items) != 2 {
		t.Fatalf("got %d items, want 2", len(res.Items))
	}
	// The snapshot must strip the answer: no explanation leaked to the quiz.
	if res.Items[0].Snapshot.Explanation != "" {
		t.Error("quiz snapshot must strip explanation")
	}
	if res.Items[0].Snapshot.HighlightRegions == nil {
		t.Error("quiz snapshot must keep highlight regions")
	}
	if res.Attempt.Mode != "per_question" {
		t.Errorf("mode = %q", res.Attempt.Mode)
	}
}

func TestStartRejectsUnpublishedOrInactive(t *testing.T) {
	bank := &fakeBank{qs: map[string][]*questions.Question{
		"structure/sentence-completion": {question(1, "structure"), question(2, "structure")},
	}}
	for _, name := range []string{"unpublished", "inactive"} {
		t.Run(name, func(t *testing.T) {
			exam := publishedExam(9, "per_question")
			if name == "unpublished" {
				exam.Published = false
			} else {
				exam.Active = false
			}
			svc := serviceWith(nil, &fakeExamSource{byID: map[int64]*exams.ExamTemplate{9: exam}}, bank)
			if _, err := svc.Start(context.Background(), 5, 9, "per_question"); err == nil {
				t.Fatal("Start() = nil, want not_found")
			}
		})
	}
}

func TestStartRejectsUnsupportedMode(t *testing.T) {
	svc := serviceWith(nil, &fakeExamSource{byID: map[int64]*exams.ExamTemplate{
		9: publishedExam(9, "overall"),
	}}, nil)
	if _, err := svc.Start(context.Background(), 5, 9, "per_question"); err == nil {
		t.Fatal("Start() should reject a mode the exam does not support")
	}
}

func TestStartConflictWithActiveAttempt(t *testing.T) {
	store := newFakeStore()
	store.active["5:9"] = &Attempt{ID: 1, UserID: 5, ExamTemplateID: 9, Status: "in_progress"}
	svc := serviceWith(store, &fakeExamSource{byID: map[int64]*exams.ExamTemplate{
		9: publishedExam(9, "per_question"),
	}}, nil)
	_, err := svc.Start(context.Background(), 5, 9, "per_question")
	if err == nil || !strings.Contains(err.Error(), "already in progress") {
		t.Fatalf("Start() error = %v, want conflict", err)
	}
}

func TestAnswerChecksState(t *testing.T) {
	store := newFakeStore()
	attempt := &Attempt{ID: 1, UserID: 5, ExamTemplateID: 9, Mode: "per_question", Status: "in_progress", StartedAt: time.Now()}
	store.attempts[1] = attempt
	store.items[1] = []*AttemptItem{{ID: 11, AttemptID: 1, QuestionID: 1}}

	svc := serviceWith(store, nil, nil)

	if err := svc.Answer(context.Background(), 5, 1, 11, i(2), i(3000)); err != nil {
		t.Fatalf("Answer() error = %v", err)
	}
	if err := svc.Answer(context.Background(), 5, 1, 11, i(9), nil); err == nil {
		t.Fatal("Answer() should reject chosen index 9")
	}

	store.attempts[1].Status = "submitted"
	if err := svc.Answer(context.Background(), 5, 1, 11, i(1), nil); err == nil {
		t.Fatal("Answer() should reject writes to a submitted attempt")
	}
}

func TestSubmitComputesReport(t *testing.T) {
	store := newFakeStore()
	attempt := &Attempt{ID: 1, UserID: 5, ExamTemplateID: 9, Mode: "overall", Status: "in_progress", StartedAt: time.Now()}
	store.attempts[1] = attempt
	store.items[1] = []*AttemptItem{
		{ID: 11, AttemptID: 1, CorrectIndex: 0, ChosenIndex: i(0), Snapshot: QuestionSnapshot{Section: "structure"}},
		{ID: 12, AttemptID: 1, CorrectIndex: 2, ChosenIndex: i(1), Snapshot: QuestionSnapshot{Section: "structure"}},
	}
	svc := serviceWith(store, &fakeExamSource{byID: map[int64]*exams.ExamTemplate{
		9: publishedExam(9, "overall"),
	}}, nil)

	view, err := svc.Submit(context.Background(), 5, 1)
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if view.Report.ScorePct != 50 {
		t.Errorf("score = %d, want 50", view.Report.ScorePct)
	}
	if store.attempts[1].Status != "submitted" {
		t.Error("attempt must be marked submitted")
	}
	if store.attempts[1].ScorePct == nil || *store.attempts[1].ScorePct != 50 {
		t.Errorf("persisted score wrong: %+v", store.attempts[1])
	}
}

func TestSubmitLateOverallClearsAnswers(t *testing.T) {
	store := newFakeStore()
	attempt := &Attempt{ID: 1, UserID: 5, ExamTemplateID: 9, Mode: "overall", Status: "in_progress",
		StartedAt: time.Now().Add(-2 * time.Minute)} // 1-minute exam started 2 minutes ago
	store.attempts[1] = attempt
	store.items[1] = []*AttemptItem{
		{ID: 11, AttemptID: 1, CorrectIndex: 0, ChosenIndex: i(0), Snapshot: QuestionSnapshot{Section: "structure"}},
	}
	svc := serviceWith(store, &fakeExamSource{byID: map[int64]*exams.ExamTemplate{
		9: publishedExam(9, "overall"),
	}}, nil)

	view, err := svc.Submit(context.Background(), 5, 1)
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if view.Report.Unanswered != 1 {
		t.Errorf("unanswered = %d, want 1 (answer past deadline is discarded)", view.Report.Unanswered)
	}
	if store.attempts[1].FinishedAt == nil || store.attempts[1].FinishedAt.After(attempt.StartedAt.Add(time.Minute)) {
		t.Error("finished_at must be capped at the deadline")
	}
}

func TestReviewRequiresSubmitted(t *testing.T) {
	store := newFakeStore()
	store.attempts[1] = &Attempt{ID: 1, UserID: 5, Status: "in_progress"}
	store.items[1] = []*AttemptItem{{ID: 11, AttemptID: 1, CorrectIndex: 0, ChosenIndex: i(0)}}
	svc := serviceWith(store, nil, nil)

	if _, err := svc.Review(context.Background(), 5, 1); err == nil {
		t.Fatal("Review() should reject an in-progress attempt")
	}
}

func TestDeadlineFor(t *testing.T) {
	start := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	a := &Attempt{Mode: "overall", StartedAt: start}
	exam := &exams.ExamTemplate{TotalMinutes: i(15)}
	d := deadlineFor(a, exam)
	if d == nil {
		t.Fatal("deadlineFor() = nil")
	}
	want := start.Add(15 * time.Minute)
	if !d.Equal(want) {
		t.Errorf("deadline = %v, want %v", d, want)
	}
	if deadlineFor(&Attempt{Mode: "per_question", StartedAt: start}, exam) != nil {
		t.Error("per_question attempts must not get a deadline")
	}
}