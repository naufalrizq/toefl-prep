-- +goose Up
CREATE TABLE exam_templates (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    section_filters JSONB NOT NULL,
    shuffle BOOLEAN NOT NULL DEFAULT TRUE,
    mode TEXT NOT NULL DEFAULT 'both' CHECK (mode IN ('per_question', 'overall', 'both')),
    seconds_per_question INT,
    total_minutes INT,
    published BOOLEAN NOT NULL DEFAULT FALSE,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE exam_templates;