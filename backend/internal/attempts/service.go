package attempts

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"toefl-prep/backend/internal/exams"
	"toefl-prep/backend/internal/grading"
	httppkg "toefl-prep/backend/internal/http"
	"toefl-prep/backend/internal/questions"
)

type ExamSource interface {
	GetByID(ctx context.Context, id int64) (*exams.ExamTemplate, error)
}

type Bank interface {
	AllActiveBySectionType(ctx context.Context, section, typ string) ([]*questions.Question, error)
}

// AttemptStore is the persistence surface the attempts service depends on.
// *Repo satisfies it; tests use a fake so the service logic is unit-testable.
type AttemptStore interface {
	InsertAttempt(ctx context.Context, a *Attempt, items []*AttemptItem) (int64, error)
	GetByID(ctx context.Context, userID, id int64) (*Attempt, error)
	ActiveByTemplate(ctx context.Context, userID, templateID int64) (*Attempt, error)
	Items(ctx context.Context, attemptID int64) ([]*AttemptItem, error)
	SetAnswer(ctx context.Context, attemptID, itemID int64, chosen *int, timeMs *int) error
	SetFlagged(ctx context.Context, attemptID, itemID int64, flagged bool) error
	Finalize(ctx context.Context, a *Attempt) error
	List(ctx context.Context, userID int64, limit int) ([]*Attempt, error)
}

type Service struct {
	store AttemptStore
	exams ExamSource
	bank  Bank
}

func NewService(store AttemptStore, exams ExamSource, bank Bank) *Service {
	return &Service{store: store, exams: exams, bank: bank}
}

type QuizItem struct {
	ID       int64            `json:"id"`
	Snapshot QuestionSnapshot `json:"question_snapshot"`
	Flagged  bool             `json:"flagged"`
}

type StartResult struct {
	Attempt  *Attempt    `json:"attempt"`
	Items    []QuizItem  `json:"items"`
}

// Start creates an attempt, snapshotting the selected questions (FR-4.1).
func (s *Service) Start(ctx context.Context, userID, examID int64, mode string) (*StartResult, error) {
	exam, err := s.exams.GetByID(ctx, examID)
	if err != nil {
		return nil, err
	}
	if !exam.Active || !exam.Published {
		return nil, httppkg.ErrNotFound
	}
	if exam.Mode != "both" && mode != exam.Mode {
		return nil, httppkg.NewError(422, "validation_failed",
			"exam does not support the requested mode")
	}
	if mode != "per_question" && mode != "overall" {
		return nil, httppkg.NewError(422, "validation_failed",
			"mode must be per_question or overall")
	}

	active, err := s.store.ActiveByTemplate(ctx, userID, examID)
	if err != nil {
		return nil, err
	}
	if active != nil {
		return nil, httppkg.NewError(409, "conflict",
			"an attempt for this exam is already in progress")
	}

	var items []*AttemptItem
	// Iterate sections and parts in deterministic order. Per (section,type) we
	// draw ONE shuffled pool and slice it across parts, so a question can never
	// be drawn twice anywhere in the attempt.
	for _, section := range exams.SectionOrder {
		cfg, ok := exam.SectionFilters[section]
		if !ok {
			continue
		}
		need := map[string]int{}
		for _, p := range cfg.Parts {
			need[p.Type] += p.Count
		}
		pools := map[string][]*questions.Question{}
		for typ, n := range need {
			qs, err := s.bank.AllActiveBySectionType(ctx, section, typ)
			if err != nil {
				return nil, err
			}
			if len(qs) < n {
				return nil, httppkg.NewError(422, "validation_failed",
					fmt.Sprintf("not enough active %s/%s questions (%d needed, %d available)",
						section, typ, n, len(qs)))
			}
			shuffled := make([]*questions.Question, len(qs))
			copy(shuffled, qs)
			rand.Shuffle(len(shuffled), func(i, j int) {
				shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
			})
			pools[typ] = shuffled
		}
		used := map[string]int{}
		for _, p := range cfg.Parts {
			pool := pools[p.Type]
			start := used[p.Type]
			for _, q := range pool[start : start+p.Count] {
				items = append(items, &AttemptItem{
					QuestionID:   q.ID,
					CorrectIndex: q.CorrectIndex,
					Snapshot: QuestionSnapshot{
						QuestionText:     q.QuestionText,
						Passage:          q.Passage,
						Options:          q.Options,
						Explanation:      q.Explanation,
						HighlightRegions: q.HighlightRegions,
						Section:          q.Section,
						Type:             q.Type,
						Part:             p.Title,
					},
				})
			}
			used[p.Type] += p.Count
		}
	}

	attempt := &Attempt{UserID: userID, ExamTemplateID: examID, Mode: mode, Status: "in_progress", StartedAt: time.Now()}
	attempt.ID, err = s.store.InsertAttempt(ctx, attempt, items)
	if err != nil {
		return nil, err
	}
	attempt.ExamTitle = exam.Title
	attempt.Deadline = deadlineFor(attempt, exam)

	quizItems := make([]QuizItem, 0, len(items))
	for _, it := range items {
		quizItems = append(quizItems, QuizItem{
			ID:      it.ID,
			Snapshot: stripAnswers(it.Snapshot),
			Flagged:  it.Flagged,
		})
	}
	return &StartResult{Attempt: attempt, Items: quizItems}, nil
}

func stripAnswers(s QuestionSnapshot) QuestionSnapshot {
	s.Explanation = ""
	if s.HighlightRegions == nil {
		s.HighlightRegions = []questions.HighlightRegion{}
	}
	return s
}

func (s *Service) QuizQuestions(ctx context.Context, userID, attemptID int64) ([]QuizItem, error) {
	a, err := s.store.GetByID(ctx, userID, attemptID)
	if err != nil {
		return nil, err
	}
	items, err := s.store.Items(ctx, a.ID)
	if err != nil {
		return nil, err
	}
	out := make([]QuizItem, 0, len(items))
	for _, it := range items {
		out = append(out, QuizItem{ID: it.ID, Snapshot: stripAnswers(it.Snapshot), Flagged: it.Flagged})
	}
	if a.Mode == "overall" && a.Deadline == nil {
		if exam, err := s.exams.GetByID(ctx, a.ExamTemplateID); err == nil {
			a.Deadline = deadlineFor(a, exam)
		}
	}
	return out, nil
}

// deadlineFor returns the absolute deadline for overall-mode attempts.
func deadlineFor(a *Attempt, exam *exams.ExamTemplate) *time.Time {
	if a.Mode != "overall" || exam.TotalMinutes == nil {
		return nil
	}
	d := a.StartedAt.Add(time.Duration(*exam.TotalMinutes) * time.Minute)
	return &d
}

func (s *Service) Answer(ctx context.Context, userID, attemptID, itemID int64, chosen *int, timeMs *int) error {
	a, err := s.store.GetByID(ctx, userID, attemptID)
	if err != nil {
		return err
	}
	if a.Status == "submitted" {
		return httppkg.NewError(409, "conflict", "attempt already submitted")
	}
	if chosen != nil && (*chosen < 0 || *chosen > 3) {
		return httppkg.NewError(422, "validation_failed", "chosen index must be 0-3")
	}
	return s.store.SetAnswer(ctx, attemptID, itemID, chosen, timeMs)
}

func (s *Service) Flag(ctx context.Context, userID, attemptID, itemID int64, flagged bool) error {
	a, err := s.store.GetByID(ctx, userID, attemptID)
	if err != nil {
		return err
	}
	if a.Status == "submitted" {
		return httppkg.NewError(409, "conflict", "attempt already submitted")
	}
	return s.store.SetFlagged(ctx, attemptID, itemID, flagged)
}

// Submit finalizes the attempt and computes the report server-side (FR-4.4).
// It is idempotent: a second call returns the same report.
func (s *Service) Submit(ctx context.Context, userID, attemptID int64) (*ReportView, error) {
	a, err := s.store.GetByID(ctx, userID, attemptID)
	if err != nil {
		return nil, err
	}
	items, err := s.store.Items(ctx, attemptID)
	if err != nil {
		return nil, err
	}

	// Overall-mode deadline enforcement (FR-4.9).
	finish := time.Now()
	if a.Mode == "overall" {
		exam, err := s.exams.GetByID(ctx, a.ExamTemplateID)
		if err != nil {
			return nil, err
		}
		if deadline := deadlineFor(a, exam); deadline != nil && finish.After(*deadline) {
			finish = *deadline
			for _, it := range items {
				it.ChosenIndex = nil // anything past the deadline is unanswered
			}
		}
	}

	report := gradeItems(items)
	summary := map[string]any{
		"score_pct":   report.ScorePct,
		"correct":     report.Correct,
		"wrong":       report.Wrong,
		"unanswered":  report.Unanswered,
		"sections":    report.Sections,
		"worst_pos":   report.WorstPOS,
	}

	a.FinishedAt = &finish
	a.Status = "submitted"
	a.ScorePct = &report.ScorePct
	a.Summary = summary
	if err := s.store.Finalize(ctx, a); err != nil {
		return nil, err
	}
	return buildReportView(a, items, report), nil
}

func (s *Service) Review(ctx context.Context, userID, attemptID int64) (*Review, error) {
	a, err := s.store.GetByID(ctx, userID, attemptID)
	if err != nil {
		return nil, err
	}
	if a.Status != "submitted" {
		return nil, httppkg.NewError(409, "conflict", "attempt not yet submitted")
	}
	items, err := s.store.Items(ctx, attemptID)
	if err != nil {
		return nil, err
	}
	report := gradeItems(items)
	view := buildReportView(a, items, report)
	return &Review{
		Attempt: a,
		Report:  report,
		Items:   view.Items,
	}, nil
}

func (s *Service) List(ctx context.Context, userID int64) ([]*Attempt, error) {
	return s.store.List(ctx, userID, 50)
}

// ReportView is the enriched payload returned on submit (result screen).
type ReportView struct {
	Attempt *Attempt         `json:"attempt"`
	Report  grading.Report   `json:"report"`
	Items   []ReviewItem     `json:"items"`
}

func gradeItems(items []*AttemptItem) grading.Report {
	g := make([]grading.Item, 0, len(items))
	for _, it := range items {
		pos := make([]string, 0, len(it.Snapshot.HighlightRegions))
		for _, r := range it.Snapshot.HighlightRegions {
			pos = append(pos, r.Pos)
		}
		g = append(g, grading.Item{
			ID:           it.ID,
			CorrectIdx:   it.CorrectIndex,
			ChosenIdx:    it.ChosenIndex,
			Section:      it.Snapshot.Section,
			HighlightPOS: pos,
		})
	}
	return grading.Grade(g)
}

func buildReportView(a *Attempt, items []*AttemptItem, report grading.Report) *ReportView {
	results := map[int64]grading.ItemResult{}
	for _, r := range report.Items {
		results[r.ID] = r
	}
	view := &ReportView{Attempt: a, Report: report}
	view.Items = make([]ReviewItem, 0, len(items))
	for _, it := range items {
		res := results[it.ID]
		regions := it.Snapshot.HighlightRegions
		if regions == nil {
			regions = []questions.HighlightRegion{}
		}
		view.Items = append(view.Items, ReviewItem{
			ID:               it.ID,
			Section:          it.Snapshot.Section,
			Type:             it.Snapshot.Type,
			Part:             it.Snapshot.Part,
			QuestionText:     it.Snapshot.QuestionText,
			Passage:          it.Snapshot.Passage,
			Options:          it.Snapshot.Options,
			CorrectIndex:     it.CorrectIndex,
			ChosenIndex:      it.ChosenIndex,
			IsCorrect:        res.IsCorrect,
			IsUnanswered:     res.IsUnanswered,
			Explanation:      it.Snapshot.Explanation,
			HighlightRegions: regions,
			TimeTakenMs:      it.TimeTakenMs,
			Flagged:          it.Flagged,
		})
	}
	return view
}

