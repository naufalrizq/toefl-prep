-- +goose Up
-- Extend the question bank to 4 sections with dedicated item types and add a
-- passage field for reading-comprehension items.
ALTER TABLE questions DROP CONSTRAINT questions_section_check;
ALTER TABLE questions DROP CONSTRAINT questions_type_check;
ALTER TABLE questions ADD CONSTRAINT questions_section_check
    CHECK (section IN ('structure', 'vocabulary', 'reading', 'grammar_adv'));
ALTER TABLE questions ADD CONSTRAINT questions_type_check
    CHECK (type IN ('sentence-completion', 'vocab-multiple-choice',
                    'reading-comprehension', 'error-identification'));
ALTER TABLE questions ADD COLUMN passage TEXT;

-- Migrate exam_templates.section_filters from the flat
--   {"structure": 8, "vocabulary": 4}
-- shape to the part-based
--   {"structure": {"parts": [{"title": "Part 1", "type": "sentence-completion", "count": 8}]}, ...}
-- shape. Type is inferred from the section.
UPDATE exam_templates
SET section_filters = (
    SELECT jsonb_object_agg(section, jsonb_build_object('parts', jsonb_build_array(
        jsonb_build_object(
            'title', 'Part 1',
            'type', CASE section
                WHEN 'structure' THEN 'sentence-completion'
                WHEN 'vocabulary' THEN 'vocab-multiple-choice'
                WHEN 'reading' THEN 'reading-comprehension'
                ELSE 'error-identification'
            END,
            'count', (f.value)::int
        )
    )))
    FROM jsonb_each(section_filters) f(section, value)
);

-- +goose Down
UPDATE exam_templates
SET section_filters = (
    SELECT jsonb_object_agg(section, (
        SELECT sum((p.value->>'count')::int) FROM jsonb_array_elements(f.value->'parts') p
    ))
    FROM jsonb_each(section_filters) f(section, value)
);
ALTER TABLE questions DROP COLUMN passage;
ALTER TABLE questions DROP CONSTRAINT questions_section_check;
ALTER TABLE questions DROP CONSTRAINT questions_type_check;
ALTER TABLE questions ADD CONSTRAINT questions_section_check
    CHECK (section IN ('structure', 'vocabulary'));
ALTER TABLE questions ADD CONSTRAINT questions_type_check
    CHECK (type IN ('sentence-completion', 'vocab-multiple-choice'));