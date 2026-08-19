-- +goose Up
CREATE TABLE attempts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    exam_template_id BIGINT NOT NULL REFERENCES exam_templates(id),
    mode TEXT NOT NULL CHECK (mode IN ('per_question', 'overall')),
    started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at TIMESTAMPTZ,
    status TEXT NOT NULL DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'submitted')),
    score_pct INT,
    summary JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_attempts_user ON attempts(user_id, started_at);

CREATE TABLE attempt_items (
    id BIGSERIAL PRIMARY KEY,
    attempt_id BIGINT NOT NULL REFERENCES attempts(id) ON DELETE CASCADE,
    question_id BIGINT NOT NULL REFERENCES questions(id),
    question_snapshot JSONB NOT NULL,
    correct_index INT NOT NULL,
    chosen_index INT,
    flagged BOOLEAN NOT NULL DEFAULT FALSE,
    time_taken_ms INT,
    answered_at TIMESTAMPTZ
);
CREATE INDEX idx_attempt_items_attempt ON attempt_items(attempt_id);

-- +goose Down
DROP TABLE attempt_items;
DROP TABLE attempts;