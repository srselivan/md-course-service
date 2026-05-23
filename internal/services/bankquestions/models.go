package bankquestions

import "course-service/internal/domain"

type BankAnswerDTO struct {
	ID        int64
	Text      string
	IsCorrect bool
}

type (
	CreateRepoParams struct {
		Text    string
		Type    domain.TestingQuestionType
		Points  int
		BankId  int64
		Answers []BankAnswerDTO
	}
	BulkCreateRepoParams struct {
		Questions []CreateRepoParams
	}
	UpdateRepoParams struct {
		ID      int64
		Text    string
		Type    domain.TestingQuestionType
		Points  int
		BankId  int64
		Answers []BankAnswerDTO
	}
	BulkUpdateRepoParams struct {
		Id        int64
		Questions []UpdateRepoParams
	}
	GetListRepoParams struct {
		BankId int64
		Limit  int64
		Offset int64
		Type   *string
		Filter *string
	}
	GetListRepoResult struct {
		Items []domain.BankQuestion
		Total int64
	}
)

type (
	CreateServiceParams struct {
		Text    string
		Type    domain.TestingQuestionType
		Points  int
		BankId  int64
		Answers []BankAnswerDTO
	}
	BulkCreateServiceParams struct {
		Questions []CreateServiceParams
	}
	UpdateServiceParams struct {
		ID      int64
		Text    string
		Type    domain.TestingQuestionType
		Points  int
		BankId  int64
		Answers []BankAnswerDTO
	}
	BulkUpdateServiceParams struct {
		Id        int64
		Questions []UpdateServiceParams
	}
	GetListServiceParams struct {
		BankId int64
		Limit  int64
		Offset int64
		Type   *string
		Filter *string
	}
)
