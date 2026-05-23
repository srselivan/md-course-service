package domain

import "time"

type BankQuestion struct {
	ID        int64               `json:"id"`
	BankId    int64               `json:"bankId"`
	Text      string              `json:"text"`
	Type      TestingQuestionType `json:"type"`
	Points    int                 `json:"points"`
	CreatedAt time.Time           `json:"createdAt"`
	Answers   []BankAnswer        `json:"answers"`
}

type BankQuestionsListMeta struct {
	Total  int64 `json:"total"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

type BankQuestionsListResponse struct {
	Meta BankQuestionsListMeta `json:"meta"`
	Data []BankQuestion        `json:"data"`
}
