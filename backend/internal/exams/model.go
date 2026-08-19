package exams

import "time"

// PartConfig defines one part inside a section: a labelled group of items of
// a single type. Parts are drawn without replacement, so a question can never
// appear twice within one attempt.
type PartConfig struct {
	Title string `json:"title"`
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// SectionConfig holds the ordered parts of one section.
type SectionConfig struct {
	Parts []PartConfig `json:"parts"`
}

type ExamTemplate struct {
	ID                 int64                     `json:"id"`
	Title              string                    `json:"title"`
	SectionFilters     map[string]SectionConfig  `json:"section_filters"`
	Shuffle            bool                      `json:"shuffle"`
	Mode               string                    `json:"mode"`
	SecondsPerQuestion *int                      `json:"seconds_per_question,omitempty"`
	TotalMinutes       *int                      `json:"total_minutes,omitempty"`
	Published          bool                      `json:"published"`
	Active             bool                      `json:"active"`
	CreatedAt          time.Time                 `json:"created_at"`
	UpdatedAt          time.Time                 `json:"updated_at"`
}

func validMode(m string) bool {
	return m == "per_question" || m == "overall" || m == "both"
}

// TotalForSection returns the sum of part counts for a section.
func TotalForSection(cfg SectionConfig) int {
	n := 0
	for _, p := range cfg.Parts {
		n += p.Count
	}
	return n
}

// TotalQuestions returns the grand total across all sections.
func (e *ExamTemplate) TotalQuestions() int {
	n := 0
	for _, cfg := range e.SectionFilters {
		n += TotalForSection(cfg)
	}
	return n
}