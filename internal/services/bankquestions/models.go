package bankquestions

import "course-service/internal/domain"

type BankAnswerDTO struct {
	ID         int64
	AnswerText string
	IsCorrect  bool
}

type (
	CreateRepoParams struct {
		QuestionText  string
		QuestionType  domain.QuestionType
		DefaultPoints int
		BankId        int64
		Answers       []BankAnswerDTO
	}
	BulkCreateRepoParams struct {
		Questions []CreateRepoParams
	}
	UpdateRepoParams struct {
		ID            int64
		QuestionText  string
		QuestionType  domain.QuestionType
		DefaultPoints int
		BankId        int64
		Answers       []BankAnswerDTO
	}
	BulkUpdateRepoParams struct {
		Id        int64
		Questions []UpdateRepoParams
	}
	GetListRepoParams struct {
		BankId *int64
	}
)

type (
	CreateServiceParams struct {
		QuestionText  string
		QuestionType  domain.QuestionType
		DefaultPoints int
		BankId        int64
		Answers       []BankAnswerDTO
	}
	BulkCreateServiceParams struct {
		Questions []CreateServiceParams
	}
	UpdateServiceParams struct {
		ID            int64
		QuestionText  string
		QuestionType  domain.QuestionType
		DefaultPoints int
		BankId        int64
		Answers       []BankAnswerDTO
	}
	BulkUpdateServiceParams struct {
		Id        int64
		Questions []UpdateServiceParams
	}
	GetListServiceParams struct {
		BankId *int64
	}
)
