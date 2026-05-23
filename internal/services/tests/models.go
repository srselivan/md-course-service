package tests

import (
	"course-service/internal/domain"
	"time"
)

type (
	CreateRepoParams struct {
		CourseID            int64
		CourseSectionItemID int64
		Title               string
		Description         string
		AvailableFrom       time.Time
		AvailableTo         time.Time
		DurationSeconds     int64
		MaxAttempts         int
		MaxScore            int
		QuestionsCount      int
		GenerationSettings  domain.JSONB
		BankIDs             []int64
	}
	UpdateRepoParams struct {
		ID                  int64
		CourseID            int64
		CourseSectionItemID int64
		Title               string
		Description         string
		AvailableFrom       time.Time
		AvailableTo         time.Time
		DurationSeconds     int64
		MaxAttempts         int
		MaxScore            int
		QuestionsCount      int
		GenerationSettings  domain.JSONB
		BankIDs             []int64
	}
	SelectBankQuestionIDsRepoParams struct {
		TestID           int64
		DifficultyLevels []int
		Limit            int
	}
	CreateAttemptRepoParams struct {
		TestID        int64
		UserID        int64
		AttemptNumber int
		StartedAt     time.Time
		ExpiresAt     time.Time
		Questions     []domain.Question
	}
	SaveTextAnswerRepoParams struct {
		AttemptQuestionID int64
		TextResponse      string
	}
	SaveChoiceAnswersRepoParams struct {
		AttemptQuestionID int64
		SelectedAnswerIDs []int64
	}
	AttemptQuestionScore struct {
		AttemptQuestionID int64
		ScoreAwarded      int
		InstructorComment *string
	}
	CompleteAttemptRepoParams struct {
		AttemptID   int64
		Status      domain.TestAttemptStatus
		GradeStatus domain.TestGradeStatus
		CompletedAt time.Time
		TotalScore  int
		Scores      []AttemptQuestionScore
	}
	GradeAttemptRepoParams struct {
		AttemptID   int64
		GradeStatus domain.TestGradeStatus
		CompletedAt time.Time
		TotalScore  int
		Scores      []AttemptQuestionScore
	}
	GetAttemptStateRepoParams struct {
		AttemptID      int64
		Limit          int
		Offset         int
		IncludeCorrect bool
	}
	GetAttemptsRepoParams struct {
		TestID *int64
		UserID *int64
		Status *domain.TestAttemptStatus
	}
)

type (
	CreateServiceParams struct {
		CourseID            int64
		CourseSectionItemID int64
		Title               string
		Description         string
		AvailableFrom       time.Time
		AvailableTo         time.Time
		DurationSeconds     int64
		MaxAttempts         int
		MaxScore            int
		QuestionsCount      int
		GenerationSettings  domain.JSONB
		BankIDs             []int64
	}
	UpdateServiceParams struct {
		ID                  int64
		CourseID            int64
		CourseSectionItemID int64
		Title               string
		Description         string
		AvailableFrom       time.Time
		AvailableTo         time.Time
		DurationSeconds     int64
		MaxAttempts         int
		MaxScore            int
		QuestionsCount      int
		GenerationSettings  domain.JSONB
		BankIDs             []int64
	}
	StartAttemptServiceParams struct {
		TestID int64
		UserID int64
		Limit  int
		Offset int
	}
	GetAttemptStateServiceParams struct {
		AttemptID int64
		Limit     int
		Offset    int
	}
	SaveAnswerServiceParams struct {
		AttemptID         int64
		AttemptQuestionID int64
		SelectedAnswerIDs []int64
		TextResponse      *string
	}
	SubmitWithAnswersServiceParams struct {
		AttemptID int64
		Answers   []SaveAnswerServiceParams
	}
	GradeAttemptServiceParams struct {
		AttemptID int64
		Scores    []AttemptQuestionScore
	}
	GetAttemptsServiceParams struct {
		TestID *int64
		UserID *int64
		Status *domain.TestAttemptStatus
	}
)
