package request

// CreateTestExample is a full request body example for POST /tests (Swagger documentation).
type CreateTestExample struct {
	CourseID            int64                 `json:"courseId" example:"1"`
	CourseSectionItemID int64                 `json:"courseSectionItemId" example:"15"`
	Title               string                `json:"title" example:"Midterm Exam"`
	Description         string                `json:"description" example:"Final assessment for module 1"`
	AvailableFrom       string                `json:"availableFrom" example:"2026-06-01T09:00:00Z"`
	AvailableTo         string                `json:"availableTo" example:"2026-06-30T23:59:59Z"`
	DurationSeconds     int64                 `json:"durationSeconds" example:"3600"`
	MaxAttempts         int                   `json:"maxAttempts" example:"3"`
	MaxScore            int                   `json:"maxScore" example:"100"`
	QuestionsCount      int                   `json:"questionsCount" example:"20"`
	GenerationSettings  CreateTestGenSettings `json:"generationSettings"`
	BankIDs             []int64               `json:"bankIds" example:"1,2"`
}

// CreateTestGenSettings documents generationSettings object shape.
type CreateTestGenSettings struct {
	ShuffleQuestions bool  `json:"shuffleQuestions" example:"true"`
	DifficultyLevels []int `json:"difficultyLevels" example:"1,2,3"`
}

// SaveAttemptAnswerSingleExample documents answer for SINGLE choice question.
type SaveAttemptAnswerSingleExample struct {
	AttemptQuestionID int64 `json:"attemptQuestionId" example:"501"`
	SelectedAnswerID  int64 `json:"selectedAnswerId" example:"1001"`
}

// SaveAttemptAnswerMultipleExample documents answer for MULTIPLE choice question.
type SaveAttemptAnswerMultipleExample struct {
	AttemptQuestionID int64   `json:"attemptQuestionId" example:"502"`
	SelectedAnswerIDs []int64 `json:"selectedAnswerIds" example:"1001,1003"`
}

// SaveAttemptAnswerTextExample documents answer for TEXT question.
type SaveAttemptAnswerTextExample struct {
	AttemptQuestionID int64  `json:"attemptQuestionId" example:"503"`
	TextResponse      string `json:"textResponse" example:"Photosynthesis converts light energy into chemical energy."`
}

// GradeAttemptExample is a full request body example for POST .../grades.
type GradeAttemptExample struct {
	Scores []GradeAttemptQuestion `json:"scores"`
}
