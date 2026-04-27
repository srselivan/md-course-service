package bankquestions

import (
	"context"
	"fmt"

	"course-service/internal/domain"
	"course-service/internal/services/bankquestions"

	"gorm.io/gorm"
)

const (
	bankQuestionsTable = "bank_question"
	bankAnswersTable   = "bank_answers"
)

type PostgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) Create(ctx context.Context, params bankquestions.CreateRepoParams) (domain.BankQuestion, error) {
	panic("implement me")
}

func (r *PostgresRepo) BulkCreate(ctx context.Context, params bankquestions.BulkCreateRepoParams) error {
	var questionsDB []bankQuestionModel
	for _, question := range params.Questions {
		questionsDB = append(questionsDB, bankQuestionModel{
			BankId:        question.BankId,
			QuestionText:  question.QuestionText,
			QuestionType:  int16(question.QuestionType),
			DefaultPoints: question.DefaultPoints,
		})
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(bankQuestionsTable).Create(&questionsDB).Error; err != nil {
			return fmt.Errorf("bulk insert %s: %w", bankQuestionsTable, err)
		}

		var answersDB []bankAnswerModel
		for i, question := range params.Questions {
			for _, answer := range question.Answers {
				answersDB = append(answersDB, bankAnswerModel{
					QuestionId: questionsDB[i].ID,
					AnswerText: answer.AnswerText,
					IsCorrect:  answer.IsCorrect,
				})
			}
		}

		if err := tx.Table(bankAnswersTable).Create(&answersDB).Error; err != nil {
			return fmt.Errorf("bulk insert %s: %w", bankAnswersTable, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (r *PostgresRepo) Update(ctx context.Context, params bankquestions.UpdateRepoParams) (domain.BankQuestion, error) {
	panic("implement me")
}

func (r *PostgresRepo) BulkUpdate(ctx context.Context, params bankquestions.BulkUpdateRepoParams) error {
	var questionsDB []bankQuestionModel
	for _, question := range params.Questions {
		questionsDB = append(questionsDB, bankQuestionModel{
			BankId:        question.BankId,
			QuestionText:  question.QuestionText,
			QuestionType:  int16(question.QuestionType),
			DefaultPoints: question.DefaultPoints,
		})
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(bankQuestionsTable).Where("bank_id = ?", params.Id).Delete(&bankQuestionModel{}).Error; err != nil {
			return fmt.Errorf("bulk delete %s: %w", bankQuestionsTable, err)
		}

		if err := tx.Table(bankQuestionsTable).Create(&questionsDB).Error; err != nil {
			return fmt.Errorf("bulk insert %s: %w", bankQuestionsTable, err)
		}

		var answersDB []bankAnswerModel
		for i, question := range params.Questions {
			for _, answer := range question.Answers {
				answersDB = append(answersDB, bankAnswerModel{
					QuestionId: questionsDB[i].ID,
					AnswerText: answer.AnswerText,
					IsCorrect:  answer.IsCorrect,
				})
			}
		}

		if err := tx.Table(bankAnswersTable).Create(&answersDB).Error; err != nil {
			return fmt.Errorf("bulk insert %s: %w", bankAnswersTable, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (r *PostgresRepo) GetList(ctx context.Context, params bankquestions.GetListRepoParams) ([]domain.BankQuestion, error) {
	var questions []bankQuestionModel

	query := r.db.WithContext(ctx).Model(&bankQuestionModel{})
	if params.BankId != nil {
		query = query.Where("bank_id = ?", *params.BankId)
	}

	if err := query.Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("get questions: %w", err)
	}

	if len(questions) == 0 {
		return []domain.BankQuestion{}, nil
	}

	ids := make([]int64, 0, len(questions))
	for _, q := range questions {
		ids = append(ids, q.ID)
	}

	var answers []bankAnswerModel
	if err := r.db.WithContext(ctx).
		Where("question_id IN ?", ids).
		Find(&answers).Error; err != nil {
		return nil, fmt.Errorf("get answers: %w", err)
	}

	answersByQuestion := make(map[int64][]bankAnswerModel, len(questions))
	for _, a := range answers {
		answersByQuestion[a.QuestionId] = append(answersByQuestion[a.QuestionId], a)
	}

	result := make([]domain.BankQuestion, 0, len(questions))
	for _, q := range questions {
		result = append(result, q.toDomain(answersByQuestion[q.ID]))
	}

	return result, nil
}
