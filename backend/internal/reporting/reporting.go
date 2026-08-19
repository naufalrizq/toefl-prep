// Package reporting computes the dashboard aggregates (PRD "Progress
// dashboard") from stored attempt summaries.
package reporting

import (
	"context"
	"time"

	"toefl-prep/backend/internal/attempts"
)

// Store is the minimal persistence surface reporting needs. attempts.Repo
// satisfies it; tests substitute a fake.
type Store interface {
	List(ctx context.Context, userID int64, limit int) ([]*attempts.Attempt, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service { return &Service{store: store} }

type SeriesPoint struct {
	ID        int64     `json:"id"`
	ScorePct  int       `json:"score_pct"`
	StartedAt time.Time `json:"started_at"`
}

type SectionStat struct {
	Correct  int  `json:"correct"`
	Total    int  `json:"total"`
	Accuracy *int `json:"accuracy,omitempty"`
}

type Stats struct {
	TotalAttempts int                   `json:"total_attempts"`
	AverageScore  *int                  `json:"average_score,omitempty"`
	BestScore     *int                  `json:"best_score,omitempty"`
	Trend         *int                  `json:"trend,omitempty"`
	Series        []SeriesPoint         `json:"series"`
	Sections      map[string]SectionStat `json:"sections"`
	WorstPOS      map[string]int        `json:"worst_pos"`
}

func (s *Service) Dashboard(ctx context.Context, userID int64) (*Stats, error) {
	items, err := s.store.List(ctx, userID, 1000)
	if err != nil {
		return nil, err
	}

	stats := &Stats{
		Series:   []SeriesPoint{},
		Sections: map[string]SectionStat{},
		WorstPOS: map[string]int{},
	}

	var submitted []*attempts.Attempt
	for _, a := range items {
		if a.Status != "submitted" || a.ScorePct == nil {
			continue
		}
		submitted = append(submitted, a)
		stats.Series = append(stats.Series, SeriesPoint{ID: a.ID, ScorePct: *a.ScorePct, StartedAt: a.StartedAt})
		if sum, ok := a.Summary.(map[string]any); ok {
			// summary.sections is stored as an array of
			// {"section","correct","total"} objects (see attempts service).
			if secs, ok := sum["sections"].([]any); ok {
				for _, v := range secs {
					sec, ok := v.(map[string]any)
					if !ok {
						continue
					}
					name, _ := sec["section"].(string)
					if name == "" {
						continue
					}
					st := stats.Sections[name]
					st.Correct += intOf(sec["correct"])
					st.Total += intOf(sec["total"])
					stats.Sections[name] = st
				}
			}
			if pos, ok := sum["worst_pos"].(string); ok && pos != "" {
				stats.WorstPOS[pos]++
			}
		}
	}

	n := len(submitted)
	stats.TotalAttempts = n
	if n == 0 {
		return stats, nil
	}

	total := 0
	best := 0
	for _, a := range submitted {
		total += *a.ScorePct
		if *a.ScorePct > best {
			best = *a.ScorePct
		}
	}
	avg := total / n
	stats.AverageScore = &avg
	stats.BestScore = &best

	// Trend: average of the most recent 5 attempts minus the 5 before that,
	// in percentage points (positive = improving). List returns DESC order.
	// Temporarily disabled; re-enable when the dashboard shows Trend again.
	// if n >= 2 {
	// 	recent := submitted[:min(n, 5)]
	// 	prev := submitted[min(n, 5):min(n, 10)]
	// 	if len(prev) > 0 {
	// 		trend := avgScore(recent) - avgScore(prev)
	// 		stats.Trend = &trend
	// 	}
	// }

	for name, st := range stats.Sections {
		if st.Total > 0 {
			acc := (st.Correct * 100) / st.Total
			st.Accuracy = &acc
		}
		stats.Sections[name] = st
	}

	return stats, nil
}

func avgScore(as []*attempts.Attempt) int {
	total := 0
	for _, a := range as {
		total += *a.ScorePct
	}
	return total / len(as)
}

func intOf(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case float64:
		return int(n)
	default:
		return 0
	}
}