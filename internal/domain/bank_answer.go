package domain

type BankAnswer struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"questionId"`
	Text       string `json:"text"`
	IsCorrect  bool   `json:"isCorrect"`
}
