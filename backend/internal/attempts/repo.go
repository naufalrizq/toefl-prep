package attempts

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	httppkg "toefl-prep/backend/internal/http"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo { return &Repo{pool: pool} }

// InsertAttempt inserts the attempt and its snapshot items atomically.
func (r *Repo) InsertAttempt(ctx context.Context, a *Attempt, items []*AttemptItem) (int64, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var id int64
	if err := tx.QueryRow(ctx, `
		INSERT INTO attempts (user_id, exam_template_id, mode)
		VALUES ($1,$2,$3) RETURNING id`,
		a.UserID, a.ExamTemplateID, a.Mode,
	).Scan(&id); err != nil {
		return 0, err
	}

	for _, it := range items {
		snap, _ := json.Marshal(it.Snapshot)
		if err := tx.QueryRow(ctx, `
			INSERT INTO attempt_items (attempt_id, question_id, question_snapshot, correct_index)
			VALUES ($1,$2,$3,$4) RETURNING id`,
			id, it.QuestionID, snap, it.CorrectIndex).Scan(&it.ID); err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return id, nil
}

const attemptCols = `id, user_id, exam_template_id, mode, started_at, finished_at, status, score_pct, summary`

func scanAttempt(row pgx.Row, title string) (*Attempt, error) {
	var (
		a         Attempt
		summary   []byte
		scorePct  *int
		finished  *time.Time
	)
	if err := row.Scan(&a.ID, &a.UserID, &a.ExamTemplateID, &a.Mode, &a.StartedAt,
		&finished, &a.Status, &scorePct, &summary); err != nil {
		return nil, err
	}
	a.ExamTitle = title
	a.FinishedAt = finished
	a.ScorePct = scorePct
	if len(summary) > 0 {
		var v any
		if err := json.Unmarshal(summary, &v); err == nil {
			a.Summary = v
		}
	}
	return &a, nil
}

func (r *Repo) GetByID(ctx context.Context, userID, id int64) (*Attempt, error) {
	a, err := scanAttempt(r.pool.QueryRow(ctx, `
		SELECT a.id, a.user_id, a.exam_template_id, a.mode, a.started_at, a.finished_at,
		       a.status, a.score_pct, a.summary
		FROM attempts a WHERE a.id=$1 AND a.user_id=$2`, id, userID), "")
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, httppkg.ErrNotFound
		}
		return nil, err
	}
	return a, nil
}

func (r *Repo) GetByIDAdmin(ctx context.Context, id int64) (*Attempt, error) {
	a, err := scanAttempt(r.pool.QueryRow(ctx, `
		SELECT a.id, a.user_id, a.exam_template_id, a.mode, a.started_at, a.finished_at,
		       a.status, a.score_pct, a.summary
		FROM attempts a WHERE a.id=$1`, id), "")
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, httppkg.ErrNotFound
		}
		return nil, err
	}
	return a, nil
}

func (r *Repo) ActiveByTemplate(ctx context.Context, userID, templateID int64) (*Attempt, error) {
	a, err := scanAttempt(r.pool.QueryRow(ctx, `
		SELECT a.id, a.user_id, a.exam_template_id, a.mode, a.started_at, a.finished_at,
		       a.status, a.score_pct, a.summary
		FROM attempts a WHERE a.user_id=$1 AND a.exam_template_id=$2 AND a.status='in_progress'
		ORDER BY a.started_at DESC LIMIT 1`, userID, templateID), "")
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}

func (r *Repo) Items(ctx context.Context, attemptID int64) ([]*AttemptItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, attempt_id, question_id, question_snapshot, correct_index, chosen_index, flagged, time_taken_ms, answered_at
		FROM attempt_items WHERE attempt_id=$1 ORDER BY id`, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*AttemptItem
	for rows.Next() {
		var (
			it       AttemptItem
			snap     []byte
			chosen   *int
			timeMs   *int
			answered *time.Time
		)
		if err := rows.Scan(&it.ID, &it.AttemptID, &it.QuestionID, &snap, &it.CorrectIndex,
			&chosen, &it.Flagged, &timeMs, &answered); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(snap, &it.Snapshot); err != nil {
			return nil, err
		}
		it.ChosenIndex = chosen
		it.TimeTakenMs = timeMs
		it.AnsweredAt = answered
		out = append(out, &it)
	}
	return out, rows.Err()
}

func (r *Repo) SetAnswer(ctx context.Context, attemptID, itemID int64, chosen *int, timeMs *int) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE attempt_items SET chosen_index=$1, time_taken_ms=$2, answered_at=now()
		WHERE id=$3 AND attempt_id=$4`, chosen, timeMs, itemID, attemptID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return httppkg.ErrNotFound
	}
	return nil
}

func (r *Repo) SetFlagged(ctx context.Context, attemptID, itemID int64, flagged bool) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE attempt_items SET flagged=$3 WHERE id=$2 AND attempt_id=$1`, attemptID, itemID, flagged)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return httppkg.ErrNotFound
	}
	return nil
}

func (r *Repo) Finalize(ctx context.Context, a *Attempt) error {
	summary, _ := json.Marshal(a.Summary)
	_, err := r.pool.Exec(ctx, `
		UPDATE attempts SET status='submitted', finished_at=$2, score_pct=$3, summary=$4
		WHERE id=$1`, a.ID, a.FinishedAt, a.ScorePct, summary)
	return err
}

func (r *Repo) List(ctx context.Context, userID int64, limit int) ([]*Attempt, error) {
	if limit < 1 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `
		SELECT a.id, a.user_id, a.exam_template_id, a.mode, a.started_at, a.finished_at,
		       a.status, a.score_pct, a.summary, e.title
		FROM attempts a JOIN exam_templates e ON e.id = a.exam_template_id
		WHERE a.user_id=$1 ORDER BY a.started_at DESC LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*Attempt
	for rows.Next() {
		var (
			title    string
			summary  []byte
			scorePct *int
			finished *time.Time
			a        Attempt
		)
		if err := rows.Scan(&a.ID, &a.UserID, &a.ExamTemplateID, &a.Mode, &a.StartedAt,
			&finished, &a.Status, &scorePct, &summary, &title); err != nil {
			return nil, err
		}
		a.ExamTitle = title
		a.FinishedAt = finished
		a.ScorePct = scorePct
		if len(summary) > 0 {
			var v any
			if json.Unmarshal(summary, &v) == nil {
				a.Summary = v
			}
		}
		out = append(out, &a)
	}
	return out, rows.Err()
}