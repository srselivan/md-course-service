package request

type BankAnswer struct {
	AnswerText string `json:"answer_text"`
	IsCorrect  bool   `json:"is_correct"`
}

type (
	CreateBankQuestion struct {
		QuestionText  string       `json:"question_text"`
		QuestionType  string       `json:"question_type"`
		DefaultPoints int          `json:"default_points"`
		BankId        int64        `json:"bank_id"`
		Answers       []BankAnswer `json:"answers"`
	}
	UpdateBankQuestion struct {
		QuestionText  string       `json:"question_text"`
		QuestionType  string       `json:"question_type"`
		DefaultPoints int          `json:"default_points"`
		BankId        int64        `json:"bank_id"`
		Answers       []BankAnswer `json:"answers"`
	}
)
