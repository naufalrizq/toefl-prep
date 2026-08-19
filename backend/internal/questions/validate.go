package questions

import (
	"fmt"
	"strings"
	"unicode"
)

func validSection(s string) bool {
	switch s {
	case "structure", "vocabulary", "reading", "grammar_adv":
		return true
	}
	return false
}

func validType(s string) bool {
	switch s {
	case "sentence-completion", "vocab-multiple-choice", "reading-comprehension", "error-identification":
		return true
	}
	return false
}

// DefaultTypeForSection returns the canonical item type for a section.
func DefaultTypeForSection(section string) string {
	switch section {
	case "structure":
		return "sentence-completion"
	case "vocabulary":
		return "vocab-multiple-choice"
	case "reading":
		return "reading-comprehension"
	case "grammar_adv":
		return "error-identification"
	}
	return ""
}

func validDifficulty(s string) bool {
	return s == "easy" || s == "medium" || s == "hard"
}

// Validate checks a question against the SRS FR-2.5 / FR-2.6 rules.
func Validate(q *Question) error {
	if !validSection(q.Section) {
		return fmt.Errorf("section must be one of: structure, vocabulary, reading, grammar_adv")
	}
	if !validType(q.Type) {
		return fmt.Errorf("type must be one of: sentence-completion, vocab-multiple-choice, reading-comprehension, error-identification")
	}
	text := strings.TrimSpace(q.QuestionText)
	if text == "" {
		return fmt.Errorf("question_text is required")
	}
	if len([]rune(text)) > 1000 {
		return fmt.Errorf("question_text exceeds 1000 characters")
	}
	if q.Type == "sentence-completion" && !strings.Contains(text, "_____") {
		return fmt.Errorf("sentence-completion question_text must contain a blank '_____'")
	}
	if q.Type == "reading-comprehension" && strings.TrimSpace(q.Passage) == "" {
		return fmt.Errorf("reading-comprehension questions require a passage")
	}
	if len([]rune(q.Passage)) > 4000 {
		return fmt.Errorf("passage exceeds 4000 characters")
	}
	if len(q.Options) != 4 {
		return fmt.Errorf("exactly 4 options are required")
	}
	seen := map[string]bool{}
	for i, opt := range q.Options {
		opt = strings.TrimSpace(opt)
		if opt == "" {
			return fmt.Errorf("option %s is empty", optionLetter(i))
		}
		if len([]rune(opt)) > 200 {
			return fmt.Errorf("option %s exceeds 200 characters", optionLetter(i))
		}
		if seen[opt] {
			return fmt.Errorf("options must be distinct")
		}
		seen[opt] = true
	}
	if q.CorrectIndex < 0 || q.CorrectIndex > 3 {
		return fmt.Errorf("correct_index must be 0-3")
	}
	if strings.TrimSpace(q.Explanation) == "" {
		return fmt.Errorf("explanation is required")
	}
	if len([]rune(q.Explanation)) > 2000 {
		return fmt.Errorf("explanation exceeds 2000 characters")
	}
	if !validDifficulty(q.Difficulty) {
		return fmt.Errorf("difficulty must be easy, medium or hard")
	}
	if err := validateRegions(q.QuestionText, q.HighlightRegions); err != nil {
		return err
	}
	return nil
}

// validateRegions enforces bounds, valid POS, and no overlap between
// regions of different POS (FR-2.6). Identical spans with the same POS are
// allowed (they deduplicate during normalization).
func validateRegions(text string, regions []HighlightRegion) error {
	runes := []rune(text)
	for _, r := range regions {
		if !AllowedPOS[r.Pos] {
			return fmt.Errorf("invalid highlight pos %q", r.Pos)
		}
		if r.Start < 0 || r.End <= r.Start || r.End > len(runes) {
			return fmt.Errorf("highlight region out of bounds")
		}
	}
	for i := 0; i < len(regions); i++ {
		for j := i + 1; j < len(regions); j++ {
			a, b := regions[i], regions[j]
			if overlap(a, b) && a.Pos != b.Pos {
				return fmt.Errorf("overlapping highlight regions with different pos")
			}
		}
	}
	return nil
}

func overlap(a, b HighlightRegion) bool {
	return a.Start < b.End && b.Start < a.End
}

func optionLetter(i int) string { return string(rune('A' + i)) }

// NormalizeHighlights converts AI-provided POS -> phrases into validated
// character-offset regions (rune-safe). Phrases not found are skipped,
// never fatal (SRS §6.5).
func NormalizeHighlights(text string, phrases map[string][]string) []HighlightRegion {
	var regions []HighlightRegion
	rs := []rune(strings.ToLower(text))
	used := map[string]bool{}
	for pos, words := range phrases {
		if !AllowedPOS[pos] {
			continue
		}
		for _, w := range words {
			w = strings.TrimSpace(w)
			if w == "" {
				continue
			}
			start := findWord(rs, []rune(strings.ToLower(w)))
			if start < 0 {
				continue
			}
			end := start + len([]rune(w))
			key := fmt.Sprintf("%d:%d", start, end)
			if used[key] {
				continue
			}
			used[key] = true
			regions = append(regions, HighlightRegion{
				Start: start,
				End:   end,
				Pos:   pos,
				Label: POSLabel(pos),
			})
		}
	}
	return regions
}

// findWord returns the rune index of the first occurrence of needle in
// haystack at a word boundary, or -1.
func findWord(haystack, needle []rune) int {
	if len(needle) == 0 || len(needle) > len(haystack) {
		return -1
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if matchAt(haystack, i, needle) && isBoundary(haystack, i, i+len(needle)) {
			return i
		}
	}
	return -1
}

func matchAt(rs []rune, start int, needle []rune) bool {
	for j := range needle {
		if rs[start+j] != needle[j] {
			return false
		}
	}
	return true
}

func isBoundary(rs []rune, start, end int) bool {
	if start > 0 && isPartOfWord(rs, start-1) {
		return false
	}
	if end < len(rs) && isPartOfWord(rs, end) {
		return false
	}
	return true
}

// isPartOfWord reports whether the rune at i belongs to the same word token as
// an adjacent match, i.e. the match is NOT at a word boundary. Apostrophes
// count as part of a word only when flanked by word characters (contractions
// and possessives like "don't"), never when they quote a term ("'ubiquitous'")
// or trail a possessive ("students'").
func isPartOfWord(rs []rune, i int) bool {
	r := rs[i]
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
		return true
	}
	if r == '\'' {
		return isWordCharAt(rs, i-1) && isWordCharAt(rs, i+1)
	}
	return false
}

func isWordCharAt(rs []rune, i int) bool {
	if i < 0 || i >= len(rs) {
		return false
	}
	r := rs[i]
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}