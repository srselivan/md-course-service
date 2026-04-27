package bankquestions

import (
	"time"

	"course-service/internal/domain"
)

type bankQuestionModel struct {
	ID            int64     `gorm:"primaryKey"`
	BankId        int64     `gorm:"column:bank_id;not null"`
	QuestionText  string    `gorm:"column:question_text;type:text;not null"`
	QuestionType  int16     `gorm:"column:question_type;not null"`
	DefaultPoints int       `gorm:"column:default_points"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
}

type bankAnswerModel struct {
	ID         int64  `gorm:"primaryKey"`
	QuestionId int64  `gorm:"column:question_id;not null"`
	AnswerText string `gorm:"column:answer_text;type:text;not null"`
	IsCorrect  bool   `gorm:"column:is_correct;not null"`
}

func (q bankQuestionModel) toDomain(answers []bankAnswerModel) domain.BankQuestion {
	dAnswers := make([]domain.BankAnswer, 0, len(answers))
	for _, a := range answers {
		dAnswers = append(dAnswers, domain.BankAnswer{
			ID:         a.ID,
			QuestionId: a.QuestionId,
			AnswerText: a.AnswerText,
			IsCorrect:  a.IsCorrect,
		})
	}

	return domain.BankQuestion{
		ID:            q.ID,
		BankId:        q.BankId,
		QuestionText:  q.QuestionText,
		QuestionType:  domain.QuestionType(q.QuestionType),
		DefaultPoints: q.DefaultPoints,
		CreatedAt:     q.CreatedAt,
		Answers:       dAnswers,
	}
}
