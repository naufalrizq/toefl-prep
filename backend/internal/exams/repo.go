package exams

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	httppkg "toefl-prep/backend/internal/http"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo { return &Repo{pool: pool} }

const examCols = `id, title, section_filters, shuffle, mode, seconds_per_question, total_minutes, published, active, created_at, updated_at`

func scanExam(row pgx.Row) (*ExamTemplate, error) {
	var (
		e         ExamTemplate
		filters   []byte
		spq, tmin *int
	)
	if err := row.Scan(&e.ID, &e.Title, &filters, &e.Shuffle, &e.Mode, &spq, &tmin,
		&e.Published, &e.Active, &e.CreatedAt, &e.UpdatedAt); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(filters, &e.SectionFilters); err != nil {
		return nil, err
	}
	e.SecondsPerQuestion = spq
	e.TotalMinutes = tmin
	return &e, nil
}

func (r *Repo) Create(ctx context.Context, e *ExamTemplate) (int64, error) {
	filters, _ := json.Marshal(e.SectionFilters)
	var id int64
	err := r.pool.QueryRow(ctx, `
		INSERT INTO exam_templates (title, section_filters, shuffle, mode, seconds_per_question, total_minutes)
		VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		e.Title, filters, e.Shuffle, e.Mode, e.SecondsPerQuestion, e.TotalMinutes,
	).Scan(&id)
	return id, err
}

func (r *Repo) Update(ctx context.Context, e *ExamTemplate) error {
	filters, _ := json.Marshal(e.SectionFilters)
	tag, err := r.pool.Exec(ctx, `
		UPDATE exam_templates SET title=$2, section_filters=$3, shuffle=$4, mode=$5,
			seconds_per_question=$6, total_minutes=$7, updated_at=now()
		WHERE id=$1`,
		e.ID, e.Title, filters, e.Shuffle, e.Mode, e.SecondsPerQuestion, e.TotalMinutes,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return httppkg.ErrNotFound
	}
	return nil
}

func (r *Repo) SoftDelete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, `UPDATE exam_templates SET active=false, updated_at=now() WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return httppkg.ErrNotFound
	}
	return nil
}

func (r *Repo) GetByID(ctx context.Context, id int64) (*ExamTemplate, error) {
	e, err := scanExam(r.pool.QueryRow(ctx,
		`SELECT `+examCols+` FROM exam_templates WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, httppkg.ErrNotFound
		}
		return nil, err
	}
	return e, nil
}

func (r *Repo) ListAll(ctx context.Context) ([]*ExamTemplate, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+examCols+` FROM exam_templates ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collect(rows)
}

func (r *Repo) ListPublished(ctx context.Context) ([]*ExamTemplate, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+examCols+` FROM exam_templates WHERE published=true AND active=true ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collect(rows)
}

func (r *Repo) SetPublished(ctx context.Context, id int64, published bool) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE exam_templates SET published=$2, updated_at=now() WHERE id=$1`, id, published)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return httppkg.ErrNotFound
	}
	return nil
}

func collect(rows pgx.Rows) ([]*ExamTemplate, error) {
	var out []*ExamTemplate
	for rows.Next() {
		e, err := scanExam(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}