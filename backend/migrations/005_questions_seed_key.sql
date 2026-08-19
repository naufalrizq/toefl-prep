-- +goose Up
-- Seed content uses (section, type, question_text) as the stable key so that
-- re-seeding updates existing rows (self-heals content fixes) instead of
-- inserting duplicates.
CREATE UNIQUE INDEX idx_questions_seed_key ON questions (section, type, question_text);

-- +goose Down
DROP INDEX idx_questions_seed_key;