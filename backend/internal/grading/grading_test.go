package grading

import "testing"

func i(idx int) *int { return &idx }

func TestGradeAllCorrect(t *testing.T) {
	items := []Item{
		{ID: 1, CorrectIdx: 2, ChosenIdx: i(2), Section: "structure", HighlightPOS: []string{"verb"}},
		{ID: 2, CorrectIdx: 0, ChosenIdx: i(0), Section: "vocabulary", HighlightPOS: []string{"noun"}},
		{ID: 3, CorrectIdx: 3, ChosenIdx: i(3), Section: "structure"},
	}
	rep := Grade(items)

	if rep.ScorePct != 100 {
		t.Errorf("score = %d, want 100", rep.ScorePct)
	}
	if rep.Correct != 3 || rep.Wrong != 0 || rep.Unanswered != 0 {
		t.Errorf("counts wrong: %+v", rep)
	}
	if rep.WorstPOS != "" {
		t.Errorf("WorstPOS = %q, want empty", rep.WorstPOS)
	}
	if len(rep.Sections) != 2 {
		t.Errorf("sections = %d, want 2", len(rep.Sections))
	}
}

func TestGradeMixed(t *testing.T) {
	items := []Item{
		{ID: 1, CorrectIdx: 1, ChosenIdx: i(0), Section: "structure", HighlightPOS: []string{"verb"}},
		{ID: 2, CorrectIdx: 0, ChosenIdx: nil, Section: "vocabulary", HighlightPOS: []string{"noun"}},
		{ID: 3, CorrectIdx: 3, ChosenIdx: i(3), Section: "structure", HighlightPOS: []string{"pronoun"}},
		{ID: 4, CorrectIdx: 2, ChosenIdx: i(1), Section: "vocabulary", HighlightPOS: []string{"noun"}},
	}
	rep := Grade(items)

	if rep.ScorePct != 25 {
		t.Errorf("score = %d, want 25", rep.ScorePct)
	}
	if rep.Correct != 1 || rep.Wrong != 2 || rep.Unanswered != 1 {
		t.Errorf("counts wrong: %+v", rep)
	}
	if rep.WorstPOS != "noun" {
		t.Errorf("WorstPOS = %q, want noun", rep.WorstPOS)
	}
	if len(rep.Sections) != 2 {
		t.Errorf("sections = %d, want 2", len(rep.Sections))
	}
	for _, s := range rep.Sections {
		switch s.Section {
		case "structure":
			if s.ScorePct != 50 {
				t.Errorf("structure score = %d, want 50", s.ScorePct)
			}
		case "vocabulary":
			if s.ScorePct != 0 {
				t.Errorf("vocab score = %d, want 0", s.ScorePct)
			}
		}
	}
}

func TestGradeEmpty(t *testing.T) {
	rep := Grade(nil)
	if rep.Total != 0 || rep.ScorePct != 0 || len(rep.Items) != 0 {
		t.Errorf("empty report wrong: %+v", rep)
	}
}

func TestGradeAllUnanswered(t *testing.T) {
	items := []Item{
		{ID: 1, CorrectIdx: 0, ChosenIdx: nil, Section: "structure"},
		{ID: 2, CorrectIdx: 1, ChosenIdx: nil, Section: "structure"},
	}
	rep := Grade(items)
	if rep.ScorePct != 0 || rep.Unanswered != 2 || rep.Correct != 0 || rep.Wrong != 0 {
		t.Errorf("unanswered report wrong: %+v", rep)
	}
}

func TestWorstPOSTie(t *testing.T) {
	items := []Item{
		{ID: 1, CorrectIdx: 1, ChosenIdx: i(0), Section: "s", HighlightPOS: []string{"verb"}},
		{ID: 2, CorrectIdx: 1, ChosenIdx: i(0), Section: "s", HighlightPOS: []string{"noun"}},
	}
	if got := Grade(items).WorstPOS; got != "noun" {
		t.Errorf("WorstPOS = %q, want noun (alphabetical tie-break)", got)
	}
}