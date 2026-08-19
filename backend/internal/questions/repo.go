package questions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	httppkg "toefl-prep/backend/internal/http"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo { return &Repo{pool: pool} }

const questionCols = `id, section, type, question_text, passage, options, correct_index, explanation, highlight_regions, difficulty, active, created_at, updated_at`

func scanQuestion(row pgx.Row) (*Question, error) {
	var (
		q         Question
		optsJSON  []byte
		regions   []byte
		passage   *string
	)
	if err := row.Scan(&q.ID, &q.Section, &q.Type, &q.QuestionText, &passage, &optsJSON,
		&q.CorrectIndex, &q.Explanation, &regions, &q.Difficulty, &q.Active,
		&q.CreatedAt, &q.UpdatedAt); err != nil {
		return nil, err
	}
	if passage != nil {
		q.Passage = *passage
	}
	if err := json.Unmarshal(optsJSON, &q.Options); err != nil {
		return nil, err
	}
	if len(regions) == 0 {
		regions = []byte("[]")
	}
	if err := json.Unmarshal(regions, &q.HighlightRegions); err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *Repo) Create(ctx context.Context, q *Question) (int64, error) {
	opts, _ := json.Marshal(q.Options)
	regions, _ := json.Marshal(q.HighlightRegions)
	var id int64
	err := r.pool.QueryRow(ctx, `
		INSERT INTO questions (section, type, question_text, passage, options, correct_index, explanation, highlight_regions, difficulty)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`,
		q.Section, q.Type, q.QuestionText, q.Passage, opts, q.CorrectIndex, q.Explanation, regions, q.Difficulty,
	).Scan(&id)
	return id, err
}

func (r *Repo) Update(ctx context.Context, q *Question) error {
	opts, _ := json.Marshal(q.Options)
	regions, _ := json.Marshal(q.HighlightRegions)
	_, err := r.pool.Exec(ctx, `
		UPDATE questions SET section=$2, type=$3, question_text=$4, passage=$5, options=$6, correct_index=$7,
			explanation=$8, highlight_regions=$9, difficulty=$10, updated_at=now()
		WHERE id=$1`,
		q.ID, q.Section, q.Type, q.QuestionText, q.Passage, opts, q.CorrectIndex, q.Explanation, regions, q.Difficulty,
	)
	return err
}

func (r *Repo) SoftDelete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, `UPDATE questions SET active=false, updated_at=now() WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return httppkg.ErrNotFound
	}
	return nil
}

func (r *Repo) GetByID(ctx context.Context, id int64) (*Question, error) {
	q, err := scanQuestion(r.pool.QueryRow(ctx,
		`SELECT `+questionCols+` FROM questions WHERE id=$1`, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, httppkg.ErrNotFound
		}
		return nil, err
	}
	return q, nil
}

type Filter struct {
	Section    string
	Type       string
	Difficulty string
	Search     string
	Page       int
	Limit      int
}

func (r *Repo) List(ctx context.Context, f Filter) ([]*Question, int64, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 {
		f.Limit = 20
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
	var (
		where  []string
		params []any
	)
	add := func(cond string, val any) {
		params = append(params, val)
		where = append(where, fmt.Sprintf(cond, len(params)))
	}
	if f.Section != "" {
		add("section = $%d", f.Section)
	}
	if f.Type != "" {
		add("type = $%d", f.Type)
	}
	if f.Difficulty != "" {
		add("difficulty = $%d", f.Difficulty)
	}
	if f.Search != "" {
		add("question_text ILIKE '%%' || $%d || '%%'", f.Search)
	}
	whereClause := ""
	if len(where) > 0 {
		whereClause = " WHERE " + strings.Join(where, " AND ")
	}

	var total int64
	if err := r.pool.QueryRow(ctx,
		`SELECT count(*) FROM questions`+whereClause, params...).Scan(&total); err != nil {
		return nil, 0, err
	}

	params = append(params, f.Limit, (f.Page-1)*f.Limit)
	rows, err := r.pool.Query(ctx, `
		SELECT `+questionCols+` FROM questions`+whereClause+`
		ORDER BY id DESC LIMIT $`+fmt.Sprint(len(params)-1)+` OFFSET $`+fmt.Sprint(len(params)),
		params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*Question
	for rows.Next() {
		q, err := scanQuestion(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, q)
	}
	return out, total, rows.Err()
}

// CountActiveBySectionType counts active questions of a section+type (bank check).
func (r *Repo) CountActiveBySectionType(ctx context.Context, section, typ string) (int, error) {
	var n int
	err := r.pool.QueryRow(ctx,
		`SELECT count(*) FROM questions WHERE section=$1 AND type=$2 AND active=true`, section, typ).Scan(&n)
	return n, err
}

// AllActiveBySectionType returns every active question of a section+type.
// The caller shuffles/assigns without replacement so parts never overlap.
func (r *Repo) AllActiveBySectionType(ctx context.Context, section, typ string) ([]*Question, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT `+questionCols+` FROM questions
		WHERE section=$1 AND type=$2 AND active=true ORDER BY id`, section, typ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Question
	for rows.Next() {
		q, err := scanQuestion(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, q)
	}
	return out, rows.Err()
}

// SeedUpsert inserts a question or updates an existing one, keyed by
// (section + type + question_text, enforced by idx_questions_seed_key).
// Re-seeding therefore self-heals content (e.g. highlight fixes) instead of
// duplicating rows.
func (r *Repo) SeedUpsert(ctx context.Context, q *Question) error {
	opts, _ := json.Marshal(q.Options)
	regions, _ := json.Marshal(q.HighlightRegions)
	var id int64
	err := r.pool.QueryRow(ctx, `
		INSERT INTO questions (section, type, question_text, passage, options, correct_index, explanation, highlight_regions, difficulty)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (section, type, question_text)
		DO UPDATE SET options = EXCLUDED.options,
		              correct_index = EXCLUDED.correct_index,
		              explanation = EXCLUDED.explanation,
		              highlight_regions = EXCLUDED.highlight_regions,
		              passage = EXCLUDED.passage,
		              difficulty = EXCLUDED.difficulty,
		              updated_at = now()
		RETURNING id`,
		q.Section, q.Type, q.QuestionText, q.Passage, opts, q.CorrectIndex, q.Explanation, regions, q.Difficulty,
	).Scan(&id)
	return err
}