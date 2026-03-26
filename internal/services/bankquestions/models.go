package bankquestions

type BankAnswerDTO struct {
	AnswerText string
	IsCorrect  bool
}

type (
	CreateRepoParams struct {
		QuestionText  string
		QuestionType  string
		DefaultPoints int
		BankId        int64
		Answers       []BankAnswerDTO
	}
	UpdateRepoParams struct {
		ID            int64
		QuestionText  string
		QuestionType  string
		DefaultPoints int
		BankId        int64
		Answers       []BankAnswerDTO
	}
	GetListRepoParams struct {
		BankId *int64
	}
)

type (
	CreateServiceParams struct {
		QuestionText  string
		QuestionType  string
		DefaultPoints int
		BankId        int64
		Answers       []BankAnswerDTO
	}
	UpdateServiceParams struct {
		ID            int64
		QuestionText  string
		QuestionType  string
		DefaultPoints int
		BankId        int64
		Answers       []BankAnswerDTO
	}
	GetListServiceParams struct {
		BankId *int64
	}
)
