package questions

import "time"

type HighlightRegion struct {
	Start int    `json:"start"`
	End   int    `json:"end"`
	Pos   string `json:"pos"`
	Label string `json:"label,omitempty"`
}

type Question struct {
	ID               int64             `json:"id"`
	Section          string            `json:"section"`
	Type             string            `json:"type"`
	QuestionText     string            `json:"question_text"`
	Passage          string            `json:"passage,omitempty"`
	Options          []string          `json:"options"`
	CorrectIndex     int               `json:"correct_index"`
	Explanation      string            `json:"explanation"`
	HighlightRegions []HighlightRegion `json:"highlight_regions"`
	Difficulty       string            `json:"difficulty"`
	Active           bool              `json:"active"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

var AllowedPOS = map[string]bool{
	"verb": true, "noun": true, "pronoun": true, "adjective": true,
	"adverb": true, "preposition": true, "conjunction": true,
	"determiner": true, "other": true,
}

func POSLabel(pos string) string {
	if AllowedPOS[pos] {
		return pos
	}
	return "other"
}