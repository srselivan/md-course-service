package tests

import (
	"time"

	"course-service/internal/domain"
)

type testModel struct {
	ID                 int64        `db:"id"`
	CourseID           int64        `db:"course_id"`
	CourseSectionID    int64        `db:"course_section_id"`
	Title              string       `db:"title"`
	Description        string       `db:"description"`
	AvailableFrom      time.Time    `db:"available_from"`
	AvailableTo        time.Time    `db:"available_to"`
	DurationSeconds    int64        `db:"duration_seconds"`
	MaxAttempts        int          `db:"max_attempts"`
	MaxScore           int          `db:"max_score"`
	QuestionsCount     int          `db:"questions_count"`
	GenerationSettings domain.JSONB `db:"generation_settings"`
	CreatedAt          time.Time    `db:"created_at"`
	BankIDs            []int64      `db:"bank_ids"`
}

func (m testModel) toDomain() domain.Test {
	return domain.Test{
		ID:                 m.ID,
		CourseID:           m.CourseID,
		CourseSectionID:    m.CourseSectionID,
		Title:              m.Title,
		Description:        m.Description,
		AvailableFrom:      m.AvailableFrom,
		AvailableTo:        m.AvailableTo,
		DurationSeconds:    m.DurationSeconds,
		MaxAttempts:        m.MaxAttempts,
		MaxScore:           m.MaxScore,
		QuestionsCount:     m.QuestionsCount,
		GenerationSettings: m.GenerationSettings,
		CreatedAt:          m.CreatedAt,
		BankIDs:            m.BankIDs,
	}
}

type testAttemptModel struct {
	ID            int64                    `db:"id"`
	TestID        int64                    `db:"test_id"`
	UserID        int64                    `db:"user_id"`
	AttemptNumber int                      `db:"attempt_number"`
	Status        domain.TestAttemptStatus `db:"status"`
	GradeStatus   domain.TestGradeStatus   `db:"grade_status"`
	StartedAt     time.Time                `db:"started_at"`
	ExpiresAt     time.Time                `db:"expires_at"`
	CompletedAt   *time.Time               `db:"completed_at"`
	TotalScore    int                      `db:"total_score"`
}

func (m testAttemptModel) toDomain() domain.TestAttempt {
	return domain.TestAttempt{
		ID:            m.ID,
		TestID:        m.TestID,
		UserID:        m.UserID,
		AttemptNumber: m.AttemptNumber,
		Status:        m.Status,
		GradeStatus:   m.GradeStatus,
		StartedAt:     m.StartedAt,
		ExpiresAt:     m.ExpiresAt,
		CompletedAt:   m.CompletedAt,
		TotalScore:    m.TotalScore,
	}
}

type testAttemptModels []testAttemptModel

func (models testAttemptModels) toDomain() []domain.TestAttempt {
	result := make([]domain.TestAttempt, 0, len(models))
	for _, model := range models {
		result = append(result, model.toDomain())
	}
	return result
}

type attemptQuestionModel struct {
	ID                int64                      `db:"id"`
	AttemptID         int64                      `db:"attempt_id"`
	QuestionID        int64                      `db:"question_id"`
	OrderIndex        int                        `db:"order_index"`
	QuestionType      domain.TestingQuestionType `db:"question_type"`
	QuestionText      string                     `db:"question_text"`
	Points            int                        `db:"points"`
	ScoreAwarded      int                        `db:"score_awarded"`
	InstructorComment *string                    `db:"instructor_comment"`
}

func (m attemptQuestionModel) toDomain() domain.AttemptQuestion {
	return domain.AttemptQuestion{
		ID:                m.ID,
		AttemptID:         m.AttemptID,
		QuestionID:        m.QuestionID,
		OrderIndex:        m.OrderIndex,
		Type:              m.QuestionType,
		Text:              m.QuestionText,
		Points:            m.Points,
		ScoreAwarded:      m.ScoreAwarded,
		InstructorComment: m.InstructorComment,
	}
}

type attemptScoreModel struct {
	ID                int64   `db:"id"`
	ScoreAwarded      int     `db:"score_awarded"`
	InstructorComment *string `db:"instructor_comment"`
}

type questionAnswerModel struct {
	ID         int64  `db:"id"`
	QuestionID int64  `db:"question_id"`
	Text       string `db:"text"`
	IsCorrect  bool   `db:"is_correct"`
}

type attemptAnswerModel struct {
	ID                int64   `db:"id"`
	AttemptQuestionID int64   `db:"attempt_question_id"`
	SelectedAnswerID  *int64  `db:"selected_answer_id"`
	TextResponse      *string `db:"text_response"`
}

type bankQuestionRow struct {
	ID            int64  `db:"id"`
	BankID        int64  `db:"bank_id"`
	QuestionType  int16  `db:"question_type"`
	QuestionText  string `db:"question_text"`
	DefaultPoints int    `db:"default_points"`
}

func (q bankQuestionRow) toDomainQuestion() domain.Question {
	return domain.Question{
		ID:     q.ID,
		BankID: q.BankID,
		Type:   toTestingQuestionType(domain.QuestionType(q.QuestionType)),
		Text:   q.QuestionText,
		Points: q.DefaultPoints,
	}
}

func toTestingQuestionType(questionType domain.QuestionType) domain.TestingQuestionType {
	switch questionType {
	case domain.QuestionTypeMultipleChoice:
		return domain.TestingQuestionTypeMultiple
	case domain.QuestionTypeText:
		return domain.TestingQuestionTypeText
	default:
		return domain.TestingQuestionTypeSingle
	}
}
