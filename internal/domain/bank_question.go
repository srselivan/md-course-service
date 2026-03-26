package domain

import "time"

type BankQuestion struct {
	ID            int64        `json:"id"`
	BankId        int64        `json:"bank_id"`
	QuestionText  string       `json:"question_text"`
	QuestionType  QuestionType `json:"question_type"`
	DefaultPoints int          `json:"default_points"`
	CreatedAt     time.Time    `json:"created_at"`
	Answers       []BankAnswer `json:"answers"`
}

type QuestionType int16

const (
	QuestionTypeSingleChoice QuestionType = iota
	QuestionTypeMultipleChoice
	QuestionTypeText
)
