package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type JSONB []byte

func (JSONB) SwaggerType() string {
	return "object"
}

func (j *JSONB) Scan(value any) error {
	if value == nil {
		*j = JSONB("{}")
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*j = append((*j)[0:0], v...)
	case string:
		*j = append((*j)[0:0], v...)
	default:
		return fmt.Errorf("unsupported JSONB scan type %T", value)
	}

	return nil
}

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return []byte("{}"), nil
	}
	if !json.Valid(j) {
		return nil, errors.New("invalid json")
	}
	return []byte(j), nil
}

func (j JSONB) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("{}"), nil
	}
	return j, nil
}

func (j *JSONB) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*j = JSONB("{}")
		return nil
	}
	if !json.Valid(data) {
		return errors.New("invalid json")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

type Test struct {
	ID                  int64     `json:"id" example:"10"`
	CourseID            int64     `json:"courseId" example:"1"`
	CourseSectionItemID int64     `json:"courseSectionItemId" example:"15"`
	Title               string    `json:"title" example:"Midterm Exam"`
	Description         string    `json:"description" example:"Final assessment for module 1"`
	AvailableFrom       time.Time `json:"availableFrom" example:"2026-06-01T09:00:00Z"`
	AvailableTo         time.Time `json:"availableTo" example:"2026-06-30T23:59:59Z"`
	DurationSeconds     int64     `json:"durationSeconds" example:"3600"`
	MaxAttempts         int       `json:"maxAttempts" example:"3"`
	MaxScore            int       `json:"maxScore" example:"100"`
	QuestionsCount      int       `json:"questionsCount" example:"20"`
	GenerationSettings  JSONB     `json:"generationSettings" swaggertype:"object"`
	CreatedAt           time.Time `json:"createdAt" example:"2026-05-01T12:00:00Z"`
	BankIDs             []int64   `json:"bankIds,omitempty" example:"1,2"`
}

type TestGenerationSettings struct {
	ShuffleQuestions bool  `json:"shuffleQuestions" example:"true"`
	DifficultyLevels []int `json:"difficultyLevels" example:"1,2,3"`
}

type QuestionBank struct {
	ID       int64  `json:"id"`
	CourseID int64  `json:"courseId"`
	Title    string `json:"title"`
}

type TestQuestionBank struct {
	TestID int64 `json:"testId"`
	BankID int64 `json:"bankId"`
}

type TestingQuestionType string

const (
	TestingQuestionTypeSingle   TestingQuestionType = "SINGLE"
	TestingQuestionTypeMultiple TestingQuestionType = "MULTIPLE"
	TestingQuestionTypeText     TestingQuestionType = "TEXT"
)

type Question struct {
	ID              int64               `json:"id" example:"100"`
	BankID          int64               `json:"bankId" example:"1"`
	Type            TestingQuestionType `json:"type" example:"SINGLE" enums:"SINGLE,MULTIPLE,TEXT"`
	Text            string              `json:"text" example:"What is 2+2?"`
	Points          int                 `json:"points" example:"5"`
	DifficultyLevel int                 `json:"difficultyLevel" example:"2"`
	Answers         []QuestionAnswer    `json:"answers,omitempty"`
}

type QuestionAnswer struct {
	ID         int64  `json:"id" example:"1001"`
	QuestionID int64  `json:"questionId" example:"100"`
	Text       string `json:"text" example:"4"`
	IsCorrect  bool   `json:"isCorrect" example:"true"`
}

type PublicQuestionAnswer struct {
	ID         int64  `json:"id" example:"1001"`
	QuestionID int64  `json:"questionId" example:"100"`
	Text       string `json:"text" example:"4"`
}

type TestAttemptStatus string

const (
	TestAttemptStatusInProgress   TestAttemptStatus = "IN_PROGRESS"
	TestAttemptStatusCompleted    TestAttemptStatus = "COMPLETED"
	TestAttemptStatusNeedsGrading TestAttemptStatus = "NEEDS_GRADING"
)

type TestGradeStatus string

const (
	TestGradeStatusNone      TestGradeStatus = "NONE"
	TestGradeStatusReady     TestGradeStatus = "READY"
	TestGradeStatusPreGraded TestGradeStatus = "PRE_GRADED"
)

type TestAttempt struct {
	ID            int64             `json:"id" example:"200"`
	TestID        int64             `json:"testId" example:"10"`
	UserID        int64             `json:"userId" example:"42"`
	AttemptNumber int               `json:"attemptNumber" example:"1"`
	Status        TestAttemptStatus `json:"status" example:"IN_PROGRESS" enums:"IN_PROGRESS,COMPLETED,NEEDS_GRADING"`
	GradeStatus   TestGradeStatus   `json:"gradeStatus" example:"NONE" enums:"NONE,READY,PRE_GRADED"`
	StartedAt     time.Time         `json:"startedAt" example:"2026-06-15T10:00:00Z"`
	ExpiresAt     time.Time         `json:"expiresAt" example:"2026-06-15T11:00:00Z"`
	CompletedAt   *time.Time        `json:"completedAt,omitempty" example:"2026-06-15T10:45:00Z"`
	TotalScore    int               `json:"totalScore" example:"85"`
}

type AttemptQuestion struct {
	ID                int64               `json:"id" example:"501"`
	AttemptID         int64               `json:"attemptId" example:"200"`
	QuestionID        int64               `json:"questionId" example:"100"`
	OrderIndex        int                 `json:"orderIndex" example:"0"`
	Type              TestingQuestionType `json:"type" example:"SINGLE" enums:"SINGLE,MULTIPLE,TEXT"`
	Text              string              `json:"text" example:"What is 2+2?"`
	Points            int                 `json:"points" example:"5"`
	ScoreAwarded      int                 `json:"scoreAwarded" example:"5"`
	InstructorComment *string             `json:"instructorComment,omitempty" example:"Correct"`
}

type AttemptAnswer struct {
	ID                int64   `json:"id" example:"601"`
	AttemptQuestionID int64   `json:"attemptQuestionId" example:"501"`
	SelectedAnswerID  *int64  `json:"selectedAnswerId,omitempty" example:"1001"`
	TextResponse      *string `json:"textResponse,omitempty" example:"The answer is 42"`
}

type AttemptPagination struct {
	Limit  int `json:"limit" example:"10"`
	Offset int `json:"offset" example:"0"`
	Total  int `json:"total" example:"20"`
}

type AttemptQuestionState struct {
	AttemptQuestion
	Question Question               `json:"question"`
	Answers  []PublicQuestionAnswer `json:"answers,omitempty"`
	Response []AttemptAnswer        `json:"response,omitempty"`
}

type AttemptState struct {
	Attempt          TestAttempt            `json:"attempt"`
	RemainingSeconds int64                  `json:"remainingSeconds" example:"2700"`
	Questions        []AttemptQuestionState `json:"questions"`
	Pagination       AttemptPagination      `json:"pagination"`
}

type TestWithAttempts struct {
	Test     Test          `json:"test"`
	Attempts []TestAttempt `json:"attempts"`
}
