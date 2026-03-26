package domain

type BankAnswer struct {
	ID         int64  `json:"id"`
	QuestionId int64  `json:"question_id"`
	AnswerText string `json:"answer_text"`
	IsCorrect  bool   `json:"is_correct"`
}
