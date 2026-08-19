package questions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	httppkg "toefl-prep/backend/internal/http"
)

// Draft is the AI output contract (prompts/question-generator.md §3).
// Highlights are POS -> phrases, normalized into regions on import.
type Draft struct {
	Section      string              `json:"section"`
	Type         string              `json:"type"`
	Difficulty   string              `json:"difficulty"`
	QuestionText string              `json:"question_text"`
	Passage      string              `json:"passage,omitempty"`
	Options      []string            `json:"options"`
	CorrectIndex int                 `json:"correct_index"`
	Explanation  string              `json:"explanation"`
	Highlights   map[string][]string `json:"highlights"`
}

type ImportResult struct {
	Index    int       `json:"index"`
	Valid    bool      `json:"valid"`
	Error    string    `json:"error,omitempty"`
	Question *Question `json:"question,omitempty"`
}

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service { return &Service{repo: repo} }

func (s *Service) Repo() *Repo { return s.repo }

// Import parses pasted AI JSON (single object or array) into validated draft
// questions. Valid items are returned ready for the admin to review and save;
// invalid items carry a reason and never block the others (FR-2.8).
func (s *Service) Import(ctx context.Context, body io.Reader) ([]ImportResult, error) {
	raw, err := io.ReadAll(io.LimitReader(body, 1<<20))
	if err != nil {
		return nil, err
	}

	var drafts []Draft
	var single Draft
	if json.Unmarshal(raw, &drafts) != nil {
		if err := json.Unmarshal(raw, &single); err != nil {
			return nil, httppkg.NewError(422, "validation_failed", "expected a question object or an array of questions")
		}
		drafts = []Draft{single}
	}

	results := make([]ImportResult, 0, len(drafts))
	for i, d := range drafts {
		q := draftToQuestion(d)
		res := ImportResult{Index: i, Question: q}
		if err := Validate(q); err != nil {
			res.Valid = false
			res.Error = err.Error()
			res.Question = nil
		} else {
			res.Valid = true
		}
		results = append(results, res)
	}
	return results, nil
}

func draftToQuestion(d Draft) *Question {
	section := d.Section
	typ := d.Type
	diff := d.Difficulty
	if section == "" {
		section = "structure"
	}
	if typ == "" {
		typ = DefaultTypeForSection(section)
	}
	if typ == "" {
		typ = "sentence-completion"
	}
	if diff == "" {
		diff = "medium"
	}
	if len(d.Options) == 4 {
		for i := range d.Options {
			if d.Options[i] == "" && d.CorrectIndex == i {
				d.Options[i] = "_____"
			}
		}
	}
	// Avoid nil slices -> JSON null.
	options := d.Options
	if options == nil {
		options = []string{}
	}
	return &Question{
		Section:          section,
		Type:             typ,
		QuestionText:     d.QuestionText,
		Passage:          d.Passage,
		Options:          options,
		CorrectIndex:     d.CorrectIndex,
		Explanation:      d.Explanation,
		HighlightRegions: NormalizeHighlights(d.QuestionText, d.Highlights),
		Difficulty:       diff,
		Active:           true,
	}
}

func (s *Service) Seed(ctx context.Context, qs []*Question) (int, error) {
	n := 0
	for _, q := range qs {
		if err := Validate(q); err != nil {
			return n, fmt.Errorf("seed question invalid: %w", err)
		}
		if err := s.repo.SeedUpsert(ctx, q); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// FromDraft converts an AI draft (prompts/question-generator.md) into a
// Question with normalized highlight regions. Used by the seed loader.
func FromDraft(d Draft) *Question {
	return draftToQuestion(d)
}