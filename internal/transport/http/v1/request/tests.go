package request

import (
	"encoding/json"
	"time"

	"course-service/internal/domain"
	"course-service/internal/services/tests"
)

type CreateTest struct {
	CourseID            int64                         `json:"courseId" example:"1"`
	CourseSectionItemID int64                         `json:"courseSectionItemId" example:"15"`
	Title               string                        `json:"title" example:"Midterm Exam"`
	Description         string                        `json:"description" example:"Final assessment for module 1"`
	AvailableFrom       time.Time                     `json:"availableFrom" example:"2026-06-01T09:00:00Z"`
	AvailableTo         time.Time                     `json:"availableTo" example:"2026-06-30T23:59:59Z"`
	DurationSeconds     int64                         `json:"durationSeconds" example:"3600"`
	MaxAttempts         int                           `json:"maxAttempts" example:"3"`
	MaxScore            int                           `json:"maxScore" example:"100"`
	QuestionsCount      int                           `json:"questionsCount" example:"20"`
	GenerationSettings  domain.TestGenerationSettings `json:"generationSettings"`
	BankIDs             []int64                       `json:"bankIds" example:"1,2"`
}

func (r CreateTest) ToService() tests.CreateServiceParams {
	settings, _ := json.Marshal(r.GenerationSettings)
	return tests.CreateServiceParams{
		CourseID:            r.CourseID,
		CourseSectionItemID: r.CourseSectionItemID,
		Title:               r.Title,
		Description:         r.Description,
		AvailableFrom:       r.AvailableFrom,
		AvailableTo:         r.AvailableTo,
		DurationSeconds:     r.DurationSeconds,
		MaxAttempts:         r.MaxAttempts,
		MaxScore:            r.MaxScore,
		QuestionsCount:      r.QuestionsCount,
		GenerationSettings:  settings,
		BankIDs:             r.BankIDs,
	}
}

type UpdateTest struct {
	ID                  int64                         `json:"id" example:"10"`
	CourseID            int64                         `json:"courseId" example:"1"`
	CourseSectionItemID int64                         `json:"courseSectionItemId" example:"15"`
	Title               string                        `json:"title" example:"Midterm Exam (updated)"`
	Description         string                        `json:"description" example:"Updated description"`
	AvailableFrom       time.Time                     `json:"availableFrom" example:"2026-06-01T09:00:00Z"`
	AvailableTo         time.Time                     `json:"availableTo" example:"2026-06-30T23:59:59Z"`
	DurationSeconds     int64                         `json:"durationSeconds" example:"3600"`
	MaxAttempts         int                           `json:"maxAttempts" example:"3"`
	MaxScore            int                           `json:"maxScore" example:"100"`
	QuestionsCount      int                           `json:"questionsCount" example:"20"`
	GenerationSettings  domain.TestGenerationSettings `json:"generationSettings"`
	BankIDs             []int64                       `json:"bankIds" example:"1,2"`
}

func (r UpdateTest) ToService(id int64) tests.UpdateServiceParams {
	settings, _ := json.Marshal(r.GenerationSettings)
	return tests.UpdateServiceParams{
		ID:                  id,
		CourseID:            r.CourseID,
		CourseSectionItemID: r.CourseSectionItemID,
		Title:               r.Title,
		Description:         r.Description,
		AvailableFrom:       r.AvailableFrom,
		AvailableTo:         r.AvailableTo,
		DurationSeconds:     r.DurationSeconds,
		MaxAttempts:         r.MaxAttempts,
		MaxScore:            r.MaxScore,
		QuestionsCount:      r.QuestionsCount,
		GenerationSettings:  domain.JSONB(settings),
		BankIDs:             r.BankIDs,
	}
}

type StartAttempt struct {
	TestID int64 `json:"testId" example:"10"`
	UserID int64 `json:"userId" example:"42"`
}

type SaveAttemptAnswer struct {
	AttemptQuestionID int64   `json:"attemptQuestionId" example:"501"`
	SelectedAnswerIDs []int64 `json:"selectedAnswerIds" example:"1001,1003"`
	SelectedAnswerID  *int64  `json:"selectedAnswerId" example:"1001"`
	TextResponse      *string `json:"textResponse" example:"The answer is 42"`
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

type SubmitAttemptAnswers struct {
	Answers []SaveAttemptAnswer `json:"answers"`
}

func (r SubmitAttemptAnswers) ToService(attemptID int64) []tests.SaveAnswerServiceParams {
	result := make([]tests.SaveAnswerServiceParams, 0, len(r.Answers))
	for _, a := range r.Answers {
		result = append(result, a.ToService(attemptID))
	}
	return result
}

type GradeAttempt struct {
	Scores []GradeAttemptQuestion `json:"scores"`
}

type GradeAttemptQuestion struct {
	AttemptQuestionID int64   `json:"attemptQuestionId" example:"501"`
	ScoreAwarded      int     `json:"scoreAwarded" example:"8"`
	InstructorComment *string `json:"instructorComment" example:"Partially correct reasoning"`
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
