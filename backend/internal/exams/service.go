package exams

import (
	"context"
	"fmt"

	httppkg "toefl-prep/backend/internal/http"
)

// QuestionBank lets the exams service validate publishability against the
// live question bank without coupling packages.
type QuestionBank interface {
	CountActiveBySectionType(ctx context.Context, section, typ string) (int, error)
}

// SectionOrder is the canonical, deterministic display order for sections.
var SectionOrder = []string{"structure", "vocabulary", "reading", "grammar_adv"}

// ValidSection reports whether s is a known section key.
func ValidSection(s string) bool {
	for _, k := range SectionOrder {
		if k == s {
			return true
		}
	}
	return false
}

// Store is the persistence surface the exams service depends on. *Repo
// satisfies it; tests use a fake so service logic is unit-testable.
type Store interface {
	Create(ctx context.Context, e *ExamTemplate) (int64, error)
	Update(ctx context.Context, e *ExamTemplate) error
	SoftDelete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*ExamTemplate, error)
	ListAll(ctx context.Context) ([]*ExamTemplate, error)
	ListPublished(ctx context.Context) ([]*ExamTemplate, error)
	SetPublished(ctx context.Context, id int64, published bool) error
}

type Service struct {
	store Store
	bank  QuestionBank
}

func NewService(store Store, bank QuestionBank) *Service {
	return &Service{store: store, bank: bank}
}

func (s *Service) Create(ctx context.Context, e *ExamTemplate) (int64, error) {
	return s.store.Create(ctx, e)
}

func (s *Service) Update(ctx context.Context, e *ExamTemplate) error {
	return s.store.Update(ctx, e)
}

func (s *Service) SoftDelete(ctx context.Context, id int64) error {
	return s.store.SoftDelete(ctx, id)
}

func (s *Service) ListAll(ctx context.Context) ([]*ExamTemplate, error) {
	return s.store.ListAll(ctx)
}

func (s *Service) ListPublished(ctx context.Context) ([]*ExamTemplate, error) {
	return s.store.ListPublished(ctx)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*ExamTemplate, error) {
	return s.store.GetByID(ctx, id)
}

func (s *Service) Validate(e *ExamTemplate) error {
	if e.Title == "" {
		return httppkg.NewError(422, "validation_failed", "title is required")
	}
	if len(e.SectionFilters) == 0 {
		return httppkg.NewError(422, "validation_failed", "section_filters must include at least one section")
	}
	for section, cfg := range e.SectionFilters {
		if !ValidSection(section) {
			return httppkg.NewError(422, "validation_failed", fmt.Sprintf("unknown section %q", section))
		}
		if len(cfg.Parts) == 0 {
			return httppkg.NewError(422, "validation_failed", fmt.Sprintf("section %q must have at least one part", section))
		}
		for _, p := range cfg.Parts {
			if p.Title == "" {
				return httppkg.NewError(422, "validation_failed", fmt.Sprintf("section %q has a part without a title", section))
			}
			if p.Type == "" {
				return httppkg.NewError(422, "validation_failed", fmt.Sprintf("part %q has no type", p.Title))
			}
			if p.Count < 1 {
				return httppkg.NewError(422, "validation_failed",
					fmt.Sprintf("part %q count must be >= 1", p.Title))
			}
		}
	}
	if !validMode(e.Mode) {
		return httppkg.NewError(422, "validation_failed", "mode must be per_question, overall or both")
	}
	if e.Mode == "per_question" || e.Mode == "both" {
		if e.SecondsPerQuestion == nil || *e.SecondsPerQuestion < 10 {
			return httppkg.NewError(422, "validation_failed", "seconds_per_question must be >= 10 for per-question mode")
		}
	}
	if e.Mode == "overall" || e.Mode == "both" {
		if e.TotalMinutes == nil || *e.TotalMinutes < 1 {
			return httppkg.NewError(422, "validation_failed", "total_minutes must be >= 1 for overall mode")
		}
	}
	return nil
}

// CanPublish checks the bank holds enough active questions for every
// section+type needed by the exam's parts (FR-3.4).
func (s *Service) CanPublish(ctx context.Context, e *ExamTemplate) error {
	for section, cfg := range e.SectionFilters {
		need := map[string]int{}
		for _, p := range cfg.Parts {
			need[p.Type] += p.Count
		}
		for typ, n := range need {
			have, err := s.bank.CountActiveBySectionType(ctx, section, typ)
			if err != nil {
				return err
			}
			if have < n {
				return httppkg.NewError(422, "validation_failed",
					fmt.Sprintf("only %d active %s/%s questions; %d required", have, section, typ, n))
			}
		}
	}
	return nil
}

func (s *Service) Publish(ctx context.Context, id int64, published bool) error {
	e, err := s.store.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if published {
		if err := s.CanPublish(ctx, e); err != nil {
			return err
		}
	}
	return s.store.SetPublished(ctx, id, published)
}