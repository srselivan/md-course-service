package bankquestions

import (
	"time"

	"course-service/internal/domain"
)

type bankQuestionModel struct {
	ID            int64     `db:"id"`
	BankId        int64     `db:"bank_id"`
	QuestionText  string    `db:"question_text"`
	QuestionType  int16     `db:"question_type"`
	DefaultPoints int       `db:"default_points"`
	CreatedAt     time.Time `db:"created_at"`
}

type bankAnswerModel struct {
	ID         int64  `db:"id"`
	QuestionId int64  `db:"question_id"`
	AnswerText string `db:"answer_text"`
	IsCorrect  bool   `db:"is_correct"`
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
