package exams

import (
	"context"
	"strings"
	"testing"

	httppkg "toefl-prep/backend/internal/http"
)

func i(v int) *int { return &v }

func examFixture() *ExamTemplate {
	return &ExamTemplate{
		Title: "Structure Basics",
		SectionFilters: map[string]SectionConfig{
			"structure": {Parts: []PartConfig{
				{Title: "Part 1", Type: "sentence-completion", Count: 5},
			}},
			"vocabulary": {Parts: []PartConfig{
				{Title: "Part 1", Type: "vocab-multiple-choice", Count: 3},
			}},
		},
		Shuffle:            true,
		Mode:               "both",
		SecondsPerQuestion: i(20),
		TotalMinutes:       i(15),
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*ExamTemplate)
		wantErr string
	}{
		{name: "valid both-mode passes", mutate: func(e *ExamTemplate) {}},
		{name: "title required", mutate: func(e *ExamTemplate) { e.Title = "" }, wantErr: "title is required"},
		{name: "at least one section", mutate: func(e *ExamTemplate) {
			e.SectionFilters = map[string]SectionConfig{}
		}, wantErr: "at least one section"},
		{name: "unknown section", mutate: func(e *ExamTemplate) {
			e.SectionFilters = map[string]SectionConfig{"listening": {Parts: []PartConfig{{Title: "Part 1", Type: "reading-comprehension", Count: 2}}}}
		}, wantErr: "unknown section"},
		{name: "section without parts", mutate: func(e *ExamTemplate) {
			e.SectionFilters = map[string]SectionConfig{"structure": {}}
		}, wantErr: "must have at least one part"},
		{name: "part without title", mutate: func(e *ExamTemplate) {
			e.SectionFilters = map[string]SectionConfig{"structure": {Parts: []PartConfig{{Type: "sentence-completion", Count: 5}}}}
		}, wantErr: "without a title"},
		{name: "part zero count", mutate: func(e *ExamTemplate) {
			e.SectionFilters = map[string]SectionConfig{"structure": {Parts: []PartConfig{{Title: "Part 1", Type: "sentence-completion", Count: 0}}}}
		}, wantErr: "count must be >= 1"},
		{name: "invalid mode", mutate: func(e *ExamTemplate) { e.Mode = "timed" }, wantErr: "mode must be"},
		{name: "per_question needs seconds", mutate: func(e *ExamTemplate) {
			e.Mode = "per_question"
			e.SecondsPerQuestion = nil
		}, wantErr: "seconds_per_question"},
		{name: "per_question minimum", mutate: func(e *ExamTemplate) {
			e.Mode = "per_question"
			e.SecondsPerQuestion = i(5)
		}, wantErr: "seconds_per_question must be >= 10"},
		{name: "overall needs minutes", mutate: func(e *ExamTemplate) {
			e.Mode = "overall"
			e.TotalMinutes = nil
		}, wantErr: "total_minutes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := examFixture()
			tt.mutate(e)
			err := (&Service{}).Validate(e)
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
				t.Errorf("error = %q, want containing %q", err.Error(), tt.wantErr)
			}
		})
	}
}

type fakeBank struct {
	counts map[string]int // key "section/type"
}

func (f *fakeBank) CountActiveBySectionType(ctx context.Context, section, typ string) (int, error) {
	return f.counts[section+"/"+typ], nil
}

type fakeStore struct {
	byID       map[int64]*ExamTemplate
	published  map[int64]bool
	createdIDs int64
}

func (f *fakeStore) Create(ctx context.Context, e *ExamTemplate) (int64, error) {
	f.createdIDs++
	e.ID = f.createdIDs
	f.byID[e.ID] = e
	return e.ID, nil
}
func (f *fakeStore) Update(ctx context.Context, e *ExamTemplate) error {
	if _, ok := f.byID[e.ID]; !ok {
		return httppkg.ErrNotFound
	}
	f.byID[e.ID] = e
	return nil
}
func (f *fakeStore) SoftDelete(ctx context.Context, id int64) error {
	if _, ok := f.byID[id]; !ok {
		return httppkg.ErrNotFound
	}
	f.byID[id].Active = false
	return nil
}
func (f *fakeStore) GetByID(ctx context.Context, id int64) (*ExamTemplate, error) {
	e, ok := f.byID[id]
	if !ok {
		return nil, httppkg.ErrNotFound
	}
	return e, nil
}
func (f *fakeStore) ListAll(ctx context.Context) ([]*ExamTemplate, error) {
	out := make([]*ExamTemplate, 0, len(f.byID))
	for _, e := range f.byID {
		out = append(out, e)
	}
	return out, nil
}
func (f *fakeStore) ListPublished(ctx context.Context) ([]*ExamTemplate, error) {
	var out []*ExamTemplate
	for _, e := range f.byID {
		if e.Published && e.Active {
			out = append(out, e)
		}
	}
	return out, nil
}
func (f *fakeStore) SetPublished(ctx context.Context, id int64, published bool) error {
	e, ok := f.byID[id]
	if !ok {
		return httppkg.ErrNotFound
	}
	f.published[id] = published
	e.Published = published
	return nil
}

func newFakeStore() *fakeStore {
	return &fakeStore{byID: map[int64]*ExamTemplate{}, published: map[int64]bool{}}
}

func TestCanPublish(t *testing.T) {
	tests := []struct {
		name    string
		bank    *fakeBank
		wantErr bool
	}{
		{
			name: "enough active questions",
			bank: &fakeBank{counts: map[string]int{"structure/sentence-completion": 5, "vocabulary/vocab-multiple-choice": 3}},
		},
		{
			name:    "not enough vocabulary",
			bank:    &fakeBank{counts: map[string]int{"structure/sentence-completion": 5, "vocabulary/vocab-multiple-choice": 2}},
			wantErr: true,
		},
		{
			name:    "no structure questions at all",
			bank:    &fakeBank{counts: map[string]int{"vocabulary/vocab-multiple-choice": 3}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{bank: tt.bank}
			err := svc.CanPublish(context.Background(), examFixture())
			if tt.wantErr && err == nil {
				t.Fatal("CanPublish() = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("CanPublish() error = %v, want nil", err)
			}
		})
	}
}

func TestPublish(t *testing.T) {
	t.Run("publish sets published when bank is sufficient", func(t *testing.T) {
		store := newFakeStore()
		exam := examFixture()
		store.byID[7] = exam
		svc := &Service{store: store, bank: &fakeBank{counts: map[string]int{"structure/sentence-completion": 5, "vocabulary/vocab-multiple-choice": 3}}}

		if err := svc.Publish(context.Background(), 7, true); err != nil {
			t.Fatalf("Publish() error = %v", err)
		}
		if !exam.Published {
			t.Error("exam should be published")
		}
	})

	t.Run("publish blocked by empty bank", func(t *testing.T) {
		store := newFakeStore()
		exam := examFixture()
		store.byID[7] = exam
		svc := &Service{store: store, bank: &fakeBank{counts: map[string]int{}}}

		err := svc.Publish(context.Background(), 7, true)
		if err == nil {
			t.Fatal("Publish() = nil, want bank error")
		}
		if exam.Published {
			t.Error("exam must stay unpublished")
		}
	})

	t.Run("unpublish does not need the bank", func(t *testing.T) {
		store := newFakeStore()
		exam := examFixture()
		exam.Published = true
		store.byID[7] = exam
		svc := &Service{store: store, bank: &fakeBank{counts: map[string]int{}}}

		if err := svc.Publish(context.Background(), 7, false); err != nil {
			t.Fatalf("Publish() error = %v", err)
		}
		if exam.Published {
			t.Error("exam should be unpublished")
		}
	})

	t.Run("unknown exam", func(t *testing.T) {
		svc := &Service{store: newFakeStore(), bank: &fakeBank{}}
		err := svc.Publish(context.Background(), 99, true)
		if err == nil || !strings.Contains(err.Error(), "not_found") {
			t.Fatalf("Publish() error = %v, want not_found", err)
		}
	})
}