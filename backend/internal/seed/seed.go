// Package seed loads the initial question bank from an embedded JSON file.
// The file uses the same Draft contract as prompts/question-generator.md so
// the same content can be pasted through the AI import pipeline.
package seed

import (
	_ "embed"
	"encoding/json"

	"toefl-prep/backend/internal/questions"
)

//go:embed questions.json
var questionsJSON []byte

func Load() ([]*questions.Question, error) {
	var drafts []questions.Draft
	if err := json.Unmarshal(questionsJSON, &drafts); err != nil {
		return nil, err
	}
	out := make([]*questions.Question, 0, len(drafts))
	for _, d := range drafts {
		out = append(out, questions.FromDraft(d))
	}
	return out, nil
}