package attempts

import (
	"time"

	"toefl-prep/backend/internal/grading"
	"toefl-prep/backend/internal/questions"
)

type Attempt struct {
	ID             int64      `json:"id"`
	UserID         int64      `json:"user_id"`
	ExamTemplateID int64      `json:"exam_template_id"`
	ExamTitle      string     `json:"exam_title,omitempty"`
	Mode           string     `json:"mode"`
	StartedAt      time.Time  `json:"started_at"`
	FinishedAt     *time.Time `json:"finished_at"`
	Status         string     `json:"status"`
	ScorePct       *int       `json:"score_pct"`
	Summary        any        `json:"summary,omitempty"`
	Deadline       *time.Time `json:"deadline,omitempty"`
}

type QuestionSnapshot struct {
	QuestionText     string                      `json:"question_text"`
	Passage          string                      `json:"passage,omitempty"`
	Options          []string                    `json:"options"`
	Explanation      string                      `json:"explanation"`
	HighlightRegions []questions.HighlightRegion `json:"highlight_regions"`
	Section          string                      `json:"section"`
	Type             string                      `json:"type"`
	Part             string                      `json:"part,omitempty"`
}

type AttemptItem struct {
	ID           int64            `json:"id"`
	AttemptID    int64            `json:"attempt_id"`
	QuestionID   int64            `json:"question_id"`
	Snapshot     QuestionSnapshot `json:"question_snapshot"`
	CorrectIndex int              `json:"-"`
	ChosenIndex  *int             `json:"chosen_index"`
	Flagged      bool             `json:"flagged"`
	TimeTakenMs  *int             `json:"time_taken_ms"`
	AnsweredAt   *time.Time       `json:"answered_at"`
}

type ReviewItem struct {
	ID               int64                      `json:"id"`
	Section          string                     `json:"section"`
	Type             string                     `json:"type"`
	Part             string                     `json:"part,omitempty"`
	QuestionText     string                     `json:"question_text"`
	Passage          string                     `json:"passage,omitempty"`
	Options          []string                   `json:"options"`
	CorrectIndex     int                        `json:"correct_index"`
	ChosenIndex      *int                       `json:"chosen_index"`
	IsCorrect        bool                       `json:"is_correct"`
	IsUnanswered     bool                       `json:"is_unanswered"`
	Explanation      string                     `json:"explanation"`
	HighlightRegions []questions.HighlightRegion `json:"highlight_regions"`
	TimeTakenMs      *int                       `json:"time_taken_ms"`
	Flagged          bool                       `json:"flagged"`
}

// Review is the full payload for GET /attempts/:id/review (FR-6.1).
type Review struct {
	Attempt *Attempt        `json:"attempt"`
	Report  grading.Report  `json:"report"`
	Items   []ReviewItem    `json:"items"`
}