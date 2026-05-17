package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type JSONB []byte

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
	ID                 int64     `json:"id"`
	CourseID           int64     `json:"course_id"`
	CourseSectionID    int64     `json:"course_section_id"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	AvailableFrom      time.Time `json:"available_from"`
	AvailableTo        time.Time `json:"available_to"`
	DurationSeconds    int64     `json:"duration_seconds"`
	MaxAttempts        int       `json:"max_attempts"`
	MaxScore           int       `json:"max_score"`
	QuestionsCount     int       `json:"questions_count"`
	GenerationSettings JSONB     `json:"generation_settings"`
	CreatedAt          time.Time `json:"created_at"`
	BankIDs            []int64   `json:"bank_ids,omitempty"`
}

type TestGenerationSettings struct {
	ShuffleQuestions bool  `json:"shuffle_questions"`
	DifficultyLevels []int `json:"difficulty_levels"`
	QuestionsCount   int   `json:"questions_count"`
}

type QuestionBank struct {
	ID       int64  `json:"id"`
	CourseID int64  `json:"course_id"`
	Title    string `json:"title"`
}

type TestQuestionBank struct {
	TestID int64 `json:"test_id"`
	BankID int64 `json:"bank_id"`
}

type TestingQuestionType string

const (
	TestingQuestionTypeSingle   TestingQuestionType = "SINGLE"
	TestingQuestionTypeMultiple TestingQuestionType = "MULTIPLE"
	TestingQuestionTypeText     TestingQuestionType = "TEXT"
)

type Question struct {
	ID              int64               `json:"id"`
	BankID          int64               `json:"bank_id"`
	Type            TestingQuestionType `json:"type"`
	Text            string              `json:"text"`
	Points          int                 `json:"points"`
	DifficultyLevel int                 `json:"difficulty_level"`
	Answers         []QuestionAnswer    `json:"answers,omitempty"`
}

type QuestionAnswer struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"question_id"`
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
}

type PublicQuestionAnswer struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"question_id"`
	Text       string `json:"text"`
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
	ID            int64             `json:"id"`
	TestID        int64             `json:"test_id"`
	UserID        int64             `json:"user_id"`
	AttemptNumber int               `json:"attempt_number"`
	Status        TestAttemptStatus `json:"status"`
	GradeStatus   TestGradeStatus   `json:"grade_status"`
	StartedAt     time.Time         `json:"started_at"`
	ExpiresAt     time.Time         `json:"expires_at"`
	CompletedAt   *time.Time        `json:"completed_at,omitempty"`
	TotalScore    int               `json:"total_score"`
}

type AttemptQuestion struct {
	ID                int64               `json:"id"`
	AttemptID         int64               `json:"attempt_id"`
	QuestionID        int64               `json:"question_id"`
	OrderIndex        int                 `json:"order_index"`
	Type              TestingQuestionType `json:"type"`
	Text              string              `json:"text"`
	Points            int                 `json:"points"`
	ScoreAwarded      int                 `json:"score_awarded"`
	InstructorComment *string             `json:"instructor_comment,omitempty"`
}

type AttemptAnswer struct {
	ID                int64   `json:"id"`
	AttemptQuestionID int64   `json:"attempt_question_id"`
	SelectedAnswerID  *int64  `json:"selected_answer_id,omitempty"`
	TextResponse      *string `json:"text_response,omitempty"`
}

type AttemptPagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type AttemptQuestionState struct {
	AttemptQuestion
	Question Question               `json:"question"`
	Answers  []PublicQuestionAnswer `json:"answers,omitempty"`
	Response []AttemptAnswer        `json:"response,omitempty"`
}

type AttemptState struct {
	Attempt          TestAttempt            `json:"attempt"`
	RemainingSeconds int64                  `json:"remaining_seconds"`
	Questions        []AttemptQuestionState `json:"questions"`
	Pagination       AttemptPagination      `json:"pagination"`
}

type TestWithAttempts struct {
	Test     Test          `json:"test"`
	Attempts []TestAttempt `json:"attempts"`
}
