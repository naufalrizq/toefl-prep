package questions

import (
	"fmt"
	"strings"
	"testing"
)

func validQuestion() *Question {
	return &Question{
		Section:      "structure",
		Type:         "sentence-completion",
		QuestionText: "The committee _____ responsible for the report.",
		Options:      []string{"is", "are", "were", "being"},
		CorrectIndex: 0,
		Explanation:  "Jawaban: is — the committee is singular.",
		Difficulty:   "medium",
		HighlightRegions: []HighlightRegion{
			{Start: 12, End: 15, Pos: "noun", Label: "noun"},
		},
		Active: true,
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*Question)
		wantErr string
	}{
		{name: "valid question passes", mutate: func(q *Question) {}},
		{name: "invalid section", mutate: func(q *Question) { q.Section = "listening" }, wantErr: "section must be"},
		{name: "invalid type", mutate: func(q *Question) { q.Type = "essay" }, wantErr: "type must be"},
		{name: "empty question text", mutate: func(q *Question) { q.QuestionText = "  " }, wantErr: "question_text is required"},
		{name: "sentence-completion without blank", mutate: func(q *Question) {
			q.QuestionText = "No blank here."
		}, wantErr: "must contain a blank"},
		{name: "vocabulary type does not need blank", mutate: func(q *Question) {
			q.Type = "vocab-multiple-choice"
			q.Section = "vocabulary"
			q.QuestionText = "A synonym of 'rapid' is _____."
		}},
		{name: "reading requires a passage", mutate: func(q *Question) {
			q.Type = "reading-comprehension"
			q.Section = "reading"
			q.QuestionText = "What is the main idea of the passage?"
			q.Passage = ""
		}, wantErr: "require a passage"},
		{name: "wrong option count", mutate: func(q *Question) {
			q.Options = []string{"a", "b", "c"}
		}, wantErr: "exactly 4 options"},
		{name: "empty option", mutate: func(q *Question) {
			q.Options = []string{"a", "", "c", "d"}
		}, wantErr: "option B is empty"},
		{name: "duplicate options", mutate: func(q *Question) {
			q.Options = []string{"a", "a", "c", "d"}
		}, wantErr: "must be distinct"},
		{name: "correct index out of range", mutate: func(q *Question) { q.CorrectIndex = 7 }, wantErr: "correct_index"},
		{name: "missing explanation", mutate: func(q *Question) { q.Explanation = "" }, wantErr: "explanation is required"},
		{name: "invalid difficulty", mutate: func(q *Question) { q.Difficulty = "impossible" }, wantErr: "difficulty"},
		{name: "invalid region pos", mutate: func(q *Question) {
			q.HighlightRegions = []HighlightRegion{{Start: 0, End: 2, Pos: "pizza"}}
		}, wantErr: "invalid highlight pos"},
		{name: "region out of bounds", mutate: func(q *Question) {
			q.HighlightRegions = []HighlightRegion{{Start: 0, End: 9999, Pos: "noun"}}
		}, wantErr: "out of bounds"},
		{name: "overlapping regions with different pos", mutate: func(q *Question) {
			q.HighlightRegions = []HighlightRegion{
				{Start: 12, End: 22, Pos: "noun"},
				{Start: 14, End: 20, Pos: "verb"},
			}
		}, wantErr: "overlapping"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := validQuestion()
			tt.mutate(q)
			err := Validate(q)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("Validate() = nil, want error containing %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Validate() error = %q, want containing %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestNormalizeHighlights(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog."
	phrases := map[string][]string{
		"verb":        {"jumps"},
		"adjective":   {"quick", "lazy"},
		"noun":        {"dog"},
		"not-a-pos":   {"fox"},              // invalid POS ignored
		"preposition": {"over", "over"},     // duplicate region skipped
	}
	phrases["noun"] = append(phrases["noun"], "does-not-exist") // phrase not found skipped
	regions := NormalizeHighlights(text, phrases)

	want := []HighlightRegion{
		{Start: 4, End: 9, Pos: "adjective", Label: "adjective"},
		{Start: 20, End: 25, Pos: "verb", Label: "verb"},
		{Start: 26, End: 30, Pos: "preposition", Label: "preposition"},
		{Start: 35, End: 39, Pos: "adjective", Label: "adjective"},
		{Start: 40, End: 43, Pos: "noun", Label: "noun"},
	}
	if len(regions) != len(want) {
		t.Fatalf("got %d regions, want %d: %+v", len(regions), len(want), regions)
	}
	// Order is map-dependent; compare as sets.
	gotSet := map[string]HighlightRegion{}
	for _, r := range regions {
		gotSet[regionKey(r)] = r
	}
	for _, r := range want {
		if _, ok := gotSet[regionKey(r)]; !ok {
			t.Errorf("missing region %+v in %+v", r, regions)
		}
	}
}

func regionKey(r HighlightRegion) string {
	return fmt.Sprintf("%d:%d:%s", r.Start, r.End, r.Pos)
}

func TestNormalizeHighlightsRuneSafe(t *testing.T) {
	// "hábito" has an accent; the phrase must match at the correct rune offset.
	regions := NormalizeHighlights("El hábito de estudiar.", map[string][]string{
		"noun": {"hábito"},
	})
	if len(regions) != 1 {
		t.Fatalf("got %d regions, want 1", len(regions))
	}
	if regions[0].Start != 3 || regions[0].End != 9 {
		t.Errorf("region = [%d,%d), want [3,9)", regions[0].Start, regions[0].End)
	}
}

func TestFindWordBoundary(t *testing.T) {
	tests := []struct {
		name     string
		haystack string
		needle   string
		want     int
	}{
		{"found at start", "running fast", "running", 0},
		{"found in middle", "the running man", "running", 4},
		{"not inside a word", "running", "run", -1},
		{"not inside apostrophe", "don't", "don", -1},
		{"quote-wrapped word", "the 'ubiquitous' term", "ubiquitous", 5},
		{"possessive before trailing apostrophe", "students' books", "students", 0},
		{"empty needle", "anything", "", -1},
		{"needle longer than haystack", "ab", "abc", -1},
		{"case sensitive by design (caller lowercases)", "Run now", "Run", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findWord([]rune(tt.haystack), []rune(tt.needle))
			if got != tt.want {
				t.Errorf("findWord(%q, %q) = %d, want %d", tt.haystack, tt.needle, got, tt.want)
			}
		})
	}
}

func TestPOSLabel(t *testing.T) {
	if got := POSLabel("noun"); got != "noun" {
		t.Errorf("POSLabel(noun) = %q", got)
	}
	if got := POSLabel("adverb"); got != "adverb" {
		t.Errorf("POSLabel(adverb) = %q", got)
	}
	if got := POSLabel("pizza"); got != "other" {
		t.Errorf("POSLabel(pizza) = %q, want other", got)
	}
}