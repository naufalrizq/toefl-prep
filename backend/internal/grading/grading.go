// Package grading is the deep module of TOEFL Prep. It is strictly pure:
// no I/O, no time, no logging. Given the same inputs it always produces the
// same report, which makes it trivially testable.
package grading

import "sort"

type Item struct {
	ID           int64
	CorrectIdx   int
	ChosenIdx    *int // nil = unanswered
	Section      string
	HighlightPOS []string // POS categories covered by this item's highlights
}

type ItemResult struct {
	ID           int64 `json:"id"`
	IsCorrect    bool  `json:"is_correct"`
	IsUnanswered bool  `json:"is_unanswered"`
}

type SectionReport struct {
	Section  string `json:"section"`
	Correct  int    `json:"correct"`
	Total    int    `json:"total"`
	ScorePct int    `json:"score_pct"`
}

type Report struct {
	ScorePct   int            `json:"score_pct"`
	Correct    int            `json:"correct"`
	Wrong      int            `json:"wrong"`
	Unanswered int            `json:"unanswered"`
	Total      int            `json:"total"`
	Sections   []SectionReport `json:"sections"`
	Items      []ItemResult    `json:"items"`
	WorstPOS   string          `json:"worst_pos"`
}

// Grade computes the report for one attempt. It never panics and never fails:
// malformed items degrade to "unanswered".
func Grade(items []Item) Report {
	rep := Report{Total: len(items)}
	rep.Items = make([]ItemResult, 0, len(items))
	sectionCounts := map[string]*SectionReport{}
	var sectionOrder []string
	wrongPOS := map[string]int{}

	for _, it := range items {
		isUnanswered := it.ChosenIdx == nil
		isCorrect := !isUnanswered && *it.ChosenIdx == it.CorrectIdx

		rep.Items = append(rep.Items, ItemResult{
			ID:           it.ID,
			IsCorrect:    isCorrect,
			IsUnanswered: isUnanswered,
		})

		switch {
		case isUnanswered:
			rep.Unanswered++
		case isCorrect:
			rep.Correct++
		default:
			rep.Wrong++
			for _, pos := range it.HighlightPOS {
				wrongPOS[pos]++
			}
		}

		if it.Section != "" {
			sr, ok := sectionCounts[it.Section]
			if !ok {
				sr = &SectionReport{Section: it.Section}
				sectionCounts[it.Section] = sr
				sectionOrder = append(sectionOrder, it.Section)
			}
			sr.Total++
			if isCorrect {
				sr.Correct++
			}
		}
	}

	if rep.Total > 0 {
		rep.ScorePct = roundPct(rep.Correct, rep.Total)
	}
	rep.Sections = make([]SectionReport, 0, len(sectionOrder))
	sort.Strings(sectionOrder)
	for _, name := range sectionOrder {
		sr := sectionCounts[name]
		sr.ScorePct = roundPct(sr.Correct, sr.Total)
		rep.Sections = append(rep.Sections, *sr)
	}
	rep.WorstPOS = worstPOS(wrongPOS)
	return rep
}

func roundPct(num, den int) int {
	if den == 0 {
		return 0
	}
	return int((float64(num)/float64(den))*100 + 0.5)
}

func worstPOS(counts map[string]int) string {
	best, bestCount := "", 0
	for pos, n := range counts {
		if n > bestCount || (n == bestCount && best != "" && pos < best) {
			best, bestCount = pos, n
		}
	}
	return best
}