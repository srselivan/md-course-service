package request

import (
	"github.com/samber/lo"

	"course-service/internal/domain"
	"course-service/internal/services/bankquestions"
)

type (
	CreateBankQuestion struct {
		Text    string                     `json:"text" example:"What is 2+2?"`
		Type    domain.TestingQuestionType `json:"type" example:"SINGLE" enums:"SINGLE,MULTIPLE,TEXT"`
		Points  int                        `json:"points" example:"5"`
		BankId  int64                      `json:"bankId" example:"1"`
		Answers []BankAnswerRequest        `json:"answers"`
	}
	BulkCreateBankQuestion struct {
		Questions []CreateBankQuestion `json:"questions"`
	}
	UpdateBankQuestion struct {
		Text    string                     `json:"text" example:"What is 2+2?"`
		Type    domain.TestingQuestionType `json:"type" example:"SINGLE" enums:"SINGLE,MULTIPLE,TEXT"`
		Points  int                        `json:"points" example:"5"`
		BankId  int64                      `json:"bankId" example:"1"`
		Answers []BankAnswerRequest        `json:"answers"`
	}
	BulkUpdateBankQuestion struct {
		Questions []UpdateBankQuestion `json:"questions"`
	}
	BankAnswerRequest struct {
		Text      string `json:"text" example:"4"`
		IsCorrect bool   `json:"isCorrect" example:"true"`
	}
)

func (cbq CreateBankQuestion) ToService() bankquestions.CreateServiceParams {
	return bankquestions.CreateServiceParams{
		Text:    cbq.Text,
		Type:    cbq.Type,
		Points:  cbq.Points,
		BankId:  cbq.BankId,
		Answers: lo.Map(cbq.Answers, func(item BankAnswerRequest, _ int) bankquestions.BankAnswerDTO { return item.ToService() }),
	}
}

func (cbq UpdateBankQuestion) ToService(id int64) bankquestions.UpdateServiceParams {
	return bankquestions.UpdateServiceParams{
		ID:      id,
		Text:    cbq.Text,
		Type:    cbq.Type,
		Points:  cbq.Points,
		BankId:  cbq.BankId,
		Answers: lo.Map(cbq.Answers, func(item BankAnswerRequest, _ int) bankquestions.BankAnswerDTO { return item.ToService() }),
	}
}

func (ba BankAnswerRequest) ToService() bankquestions.BankAnswerDTO {
	return bankquestions.BankAnswerDTO{
		Text:      ba.Text,
		IsCorrect: ba.IsCorrect,
	}
}
