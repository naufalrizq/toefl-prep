package reporting

import (
	"context"
	"testing"
	"time"

	"toefl-prep/backend/internal/attempts"
)

func i(v int) *int { return &v }

type sectionStat struct {
	Section string `json:"section"`
	Correct int    `json:"correct"`
	Total   int    `json:"total"`
}

func submittedAttempt(id int64, score int, at time.Time, sections []sectionStat, worst string) *attempts.Attempt {
	secs := make([]any, 0, len(sections))
	for _, s := range sections {
		secs = append(secs, map[string]any{
			"section": s.Section,
			"correct": s.Correct,
			"total":   s.Total,
		})
	}
	return &attempts.Attempt{
		ID:        id,
		Status:    "submitted",
		ScorePct:  i(score),
		StartedAt: at,
		Summary: map[string]any{
			"score_pct": score,
			"sections":  secs,
			"worst_pos": worst,
		},
	}
}

type fakeStore struct {
	items []*attempts.Attempt
}

func (f *fakeStore) List(ctx context.Context, userID int64, limit int) ([]*attempts.Attempt, error) {
	return f.items, nil
}

func TestDashboardEmpty(t *testing.T) {
	svc := NewService(&fakeStore{})
	stats, err := svc.Dashboard(context.Background(), 1)
	if err != nil {
		t.Fatalf("Dashboard() error = %v", err)
	}
	if stats.TotalAttempts != 0 {
		t.Errorf("total_attempts = %d, want 0", stats.TotalAttempts)
	}
	if stats.AverageScore != nil || stats.BestScore != nil || stats.Trend != nil {
		t.Errorf("nil aggregates expected for empty history: %+v", stats)
	}
	if len(stats.Series) != 0 {
		t.Errorf("series = %d, want 0", len(stats.Series))
	}
}

func TestDashboardSkipsInProgress(t *testing.T) {
	svc := NewService(&fakeStore{items: []*attempts.Attempt{
		{ID: 1, Status: "in_progress"},
		submittedAttempt(2, 80, time.Now(), nil, ""),
	}})
	stats, err := svc.Dashboard(context.Background(), 1)
	if err != nil {
		t.Fatalf("Dashboard() error = %v", err)
	}
	if stats.TotalAttempts != 1 {
		t.Errorf("total_attempts = %d, want 1", stats.TotalAttempts)
	}
}

func TestDashboardAggregates(t *testing.T) {
	now := time.Now()
	items := []*attempts.Attempt{
		submittedAttempt(1, 50, now.Add(-2*time.Hour), []sectionStat{
			{Section: "structure", Correct: 2, Total: 4},
		}, "verb"),
		submittedAttempt(2, 70, now.Add(-1*time.Hour), []sectionStat{
			{Section: "structure", Correct: 3, Total: 4},
		}, "verb"),
		submittedAttempt(3, 90, now, []sectionStat{
			{Section: "structure", Correct: 4, Total: 4},
			{Section: "vocabulary", Correct: 1, Total: 2},
		}, "noun"),
	}
	svc := NewService(&fakeStore{items: items})

	stats, err := svc.Dashboard(context.Background(), 1)
	if err != nil {
		t.Fatalf("Dashboard() error = %v", err)
	}

	if stats.TotalAttempts != 3 {
		t.Errorf("total_attempts = %d, want 3", stats.TotalAttempts)
	}
	if stats.AverageScore == nil || *stats.AverageScore != 70 {
		t.Errorf("average = %v, want 70", stats.AverageScore)
	}
	if stats.BestScore == nil || *stats.BestScore != 90 {
		t.Errorf("best = %v, want 90", stats.BestScore)
	}
	// Series must be ordered by attempt (already DESC) and keep scores.
	if len(stats.Series) != 3 || stats.Series[0].ScorePct != 50 {
		t.Errorf("series wrong: %+v", stats.Series)
	}
	// Section aggregation: structure 9/12 = 75, vocabulary 1/2 = 50.
	sec := stats.Sections["structure"]
	if sec.Correct != 9 || sec.Total != 12 || sec.Accuracy == nil || *sec.Accuracy != 75 {
		t.Errorf("structure stat wrong: %+v", sec)
	}
	vocab := stats.Sections["vocabulary"]
	if vocab.Total != 2 || vocab.Accuracy == nil || *vocab.Accuracy != 50 {
		t.Errorf("vocabulary stat wrong: %+v", vocab)
	}
	// Worst POS frequency.
	if stats.WorstPOS["verb"] != 2 || stats.WorstPOS["noun"] != 1 {
		t.Errorf("worst_pos counts wrong: %+v", stats.WorstPOS)
	}
}

func TestDashboardAggregatesRealJSONShape(t *testing.T) {
	// Regression: the summary JSONB column decodes `sections` as an array of
	// objects (see /attempts response), not a map keyed by section name.
	now := time.Now()
	items := []*attempts.Attempt{
		submittedAttempt(1, 25, now.Add(-1*time.Hour), []sectionStat{
			{Section: "structure", Correct: 2, Total: 4},
			{Section: "vocabulary", Correct: 0, Total: 4},
		}, "noun"),
	}
	svc := NewService(&fakeStore{items: items})

	stats, err := svc.Dashboard(context.Background(), 1)
	if err != nil {
		t.Fatalf("Dashboard() error = %v", err)
	}
	if stats.Sections["structure"].Correct != 2 || stats.Sections["structure"].Total != 4 {
		t.Errorf("structure stat wrong: %+v", stats.Sections["structure"])
	}
	if stats.Sections["vocabulary"].Total != 4 {
		t.Errorf("vocabulary stat wrong: %+v", stats.Sections["vocabulary"])
	}
}

// Trend is temporarily disabled (see Dashboard), so this test is suspended too.
// func TestDashboardTrend(t *testing.T) {
// 	now := time.Now()
// 	// newest-first (matches repo DESC order): recent 5 avg 80, prev 5 avg 50.
// 	scores := []int{90, 85, 80, 75, 70, 60, 55, 50, 45, 40}
// 	items := make([]*attempts.Attempt, 0, len(scores))
// 	for idx, s := range scores {
// 		items = append(items, submittedAttempt(int64(idx+1), s, now.Add(-time.Duration(idx)*time.Minute), nil, ""))
// 	}
// 	svc := NewService(&fakeStore{items: items})
// 	stats, err := svc.Dashboard(context.Background(), 1)
// 	if err != nil {
// 		t.Fatalf("Dashboard() error = %v", err)
// 	}
// 	if stats.Trend == nil || *stats.Trend != 30 {
// 		t.Errorf("trend = %v, want 30", stats.Trend)
// 	}
// }