-- +goose Up
CREATE TABLE questions (
    id BIGSERIAL PRIMARY KEY,
    section TEXT NOT NULL CHECK (section IN ('structure', 'vocabulary')),
    type TEXT NOT NULL CHECK (type IN ('sentence-completion', 'vocab-multiple-choice')),
    question_text TEXT NOT NULL,
    options JSONB NOT NULL,
    correct_index INT NOT NULL CHECK (correct_index BETWEEN 0 AND 3),
    explanation TEXT NOT NULL,
    highlight_regions JSONB NOT NULL DEFAULT '[]'::jsonb,
    difficulty TEXT NOT NULL DEFAULT 'medium' CHECK (difficulty IN ('easy', 'medium', 'hard')),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_questions_section_active ON questions(section, active);
CREATE INDEX idx_questions_type ON questions(type);

-- +goose Down
DROP TABLE questions;