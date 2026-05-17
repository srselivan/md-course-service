package request

import (
	"course-service/internal/domain"
	"course-service/internal/services/tests"
	"time"
)

type CreateTest struct {
	CourseID           int64        `json:"course_id"`
	CourseSectionID    int64        `json:"course_section_id"`
	Title              string       `json:"title"`
	Description        string       `json:"description"`
	AvailableFrom      time.Time    `json:"available_from"`
	AvailableTo        time.Time    `json:"available_to"`
	DurationSeconds    int64        `json:"duration_seconds"`
	MaxAttempts        int          `json:"max_attempts"`
	MaxScore           int          `json:"max_score"`
	QuestionsCount     int          `json:"questions_count"`
	GenerationSettings domain.JSONB `json:"generation_settings"`
	BankIDs            []int64      `json:"bank_ids"`
}

func (r CreateTest) ToService() tests.CreateServiceParams {
	return tests.CreateServiceParams{
		CourseID:           r.CourseID,
		CourseSectionID:    r.CourseSectionID,
		Title:              r.Title,
		Description:        r.Description,
		AvailableFrom:      r.AvailableFrom,
		AvailableTo:        r.AvailableTo,
		DurationSeconds:    r.DurationSeconds,
		MaxAttempts:        r.MaxAttempts,
		MaxScore:           r.MaxScore,
		QuestionsCount:     r.QuestionsCount,
		GenerationSettings: r.GenerationSettings,
		BankIDs:            r.BankIDs,
	}
}

type UpdateTest struct {
	ID                 int64        `json:"id"`
	CourseID           int64        `json:"course_id"`
	CourseSectionID    int64        `json:"course_section_id"`
	Title              string       `json:"title"`
	Description        string       `json:"description"`
	AvailableFrom      time.Time    `json:"available_from"`
	AvailableTo        time.Time    `json:"available_to"`
	DurationSeconds    int64        `json:"duration_seconds"`
	MaxAttempts        int          `json:"max_attempts"`
	MaxScore           int          `json:"max_score"`
	QuestionsCount     int          `json:"questions_count"`
	GenerationSettings domain.JSONB `json:"generation_settings"`
	BankIDs            []int64      `json:"bank_ids"`
}

func (r UpdateTest) ToService(id int64) tests.UpdateServiceParams {
	return tests.UpdateServiceParams{
		ID:                 id,
		CourseID:           r.CourseID,
		CourseSectionID:    r.CourseSectionID,
		Title:              r.Title,
		Description:        r.Description,
		AvailableFrom:      r.AvailableFrom,
		AvailableTo:        r.AvailableTo,
		DurationSeconds:    r.DurationSeconds,
		MaxAttempts:        r.MaxAttempts,
		MaxScore:           r.MaxScore,
		QuestionsCount:     r.QuestionsCount,
		GenerationSettings: r.GenerationSettings,
		BankIDs:            r.BankIDs,
	}
}

type StartAttempt struct {
	TestID int64 `json:"test_id"`
	UserID int64 `json:"user_id"`
}

type SaveAttemptAnswer struct {
	AttemptQuestionID int64   `json:"attempt_question_id"`
	SelectedAnswerIDs []int64 `json:"selected_answer_ids"`
	SelectedAnswerID  *int64  `json:"selected_answer_id"`
	TextResponse      *string `json:"text_response"`
}

func (r SaveAttemptAnswer) ToService(attemptID int64) tests.SaveAnswerServiceParams {
	selectedAnswerIDs := r.SelectedAnswerIDs
	if r.SelectedAnswerID != nil && len(selectedAnswerIDs) == 0 {
		selectedAnswerIDs = []int64{*r.SelectedAnswerID}
	}

	return tests.SaveAnswerServiceParams{
		AttemptID:         attemptID,
		AttemptQuestionID: r.AttemptQuestionID,
		SelectedAnswerIDs: selectedAnswerIDs,
		TextResponse:      r.TextResponse,
	}
}

type GradeAttempt struct {
	Scores []GradeAttemptQuestion `json:"scores"`
}

type GradeAttemptQuestion struct {
	AttemptQuestionID int64   `json:"attempt_question_id"`
	ScoreAwarded      int     `json:"score_awarded"`
	InstructorComment *string `json:"instructor_comment"`
}

func (r GradeAttempt) ToService(attemptID int64) tests.GradeAttemptServiceParams {
	scores := make([]tests.AttemptQuestionScore, 0, len(r.Scores))
	for _, score := range r.Scores {
		scores = append(scores, tests.AttemptQuestionScore{
			AttemptQuestionID: score.AttemptQuestionID,
			ScoreAwarded:      score.ScoreAwarded,
			InstructorComment: score.InstructorComment,
		})
	}

	return tests.GradeAttemptServiceParams{
		AttemptID: attemptID,
		Scores:    scores,
	}
}
