package request

import (
	"github.com/samber/lo"

	"course-service/internal/domain"
	"course-service/internal/services/bankquestions"
)

type (
	CreateBankQuestion struct {
		QuestionText  string       `json:"question_text"`
		QuestionType  int16        `json:"question_type"`
		DefaultPoints int          `json:"default_points"`
		BankId        int64        `json:"bank_id"`
		Answers       []BankAnswer `json:"answers"`
	}
	BulkCreateBankQuestion struct {
		Questions []CreateBankQuestion `json:"questions"`
	}
	UpdateBankQuestion struct {
		QuestionText  string       `json:"question_text"`
		QuestionType  int16        `json:"question_type"`
		DefaultPoints int          `json:"default_points"`
		BankId        int64        `json:"bank_id"`
		Answers       []BankAnswer `json:"answers"`
	}
	BulkUpdateBankQuestion struct {
		Questions []UpdateBankQuestion `json:"questions"`
	}
	BankAnswer struct {
		AnswerText string `json:"answer_text"`
		IsCorrect  bool   `json:"is_correct"`
	}
)

func (cbq CreateBankQuestion) ToService() bankquestions.CreateServiceParams {
	return bankquestions.CreateServiceParams{
		QuestionText:  cbq.QuestionText,
		QuestionType:  domain.QuestionType(cbq.QuestionType),
		DefaultPoints: cbq.DefaultPoints,
		BankId:        cbq.BankId,
		Answers:       lo.Map(cbq.Answers, func(item BankAnswer, _ int) bankquestions.BankAnswerDTO { return item.ToService() }),
	}
}

func (cbq UpdateBankQuestion) ToService() bankquestions.UpdateServiceParams {
	return bankquestions.UpdateServiceParams{
		QuestionText:  cbq.QuestionText,
		QuestionType:  domain.QuestionType(cbq.QuestionType),
		DefaultPoints: cbq.DefaultPoints,
		BankId:        cbq.BankId,
		Answers:       lo.Map(cbq.Answers, func(item BankAnswer, _ int) bankquestions.BankAnswerDTO { return item.ToService() }),
	}
}

func (ba BankAnswer) ToService() bankquestions.BankAnswerDTO {
	return bankquestions.BankAnswerDTO{
		AnswerText: ba.AnswerText,
		IsCorrect:  ba.IsCorrect,
	}
}
