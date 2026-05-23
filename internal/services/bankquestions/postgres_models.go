package bankquestions

import (
	"time"

	"course-service/internal/domain"
)

type bankQuestionModel struct {
	ID        int64     `db:"id"`
	BankId    int64     `db:"bank_id"`
	Text      string    `db:"text"`
	Type      string    `db:"type"`
	Points    int       `db:"points"`
	CreatedAt time.Time `db:"created_at"`
}

type bankAnswerModel struct {
	ID         int64  `db:"id"`
	QuestionId int64  `db:"question_id"`
	Text       string `db:"text"`
	IsCorrect  bool   `db:"is_correct"`
}

func (q bankQuestionModel) toDomain(answers []bankAnswerModel) domain.BankQuestion {
	dAnswers := make([]domain.BankAnswer, 0, len(answers))
	for _, a := range answers {
		dAnswers = append(dAnswers, domain.BankAnswer{
			ID:         a.ID,
			QuestionID: a.QuestionId,
			Text:       a.Text,
			IsCorrect:  a.IsCorrect,
		})
	}

	return domain.BankQuestion{
		ID:        q.ID,
		BankId:    q.BankId,
		Text:      q.Text,
		Type:      domain.TestingQuestionType(q.Type),
		Points:    q.Points,
		CreatedAt: q.CreatedAt,
		Answers:   dAnswers,
	}
}
