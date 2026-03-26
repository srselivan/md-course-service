package bankquestions

import (
	"context"
	"course-service/internal/domain"
	"course-service/internal/services/bankquestions"
	"fmt"

	"gorm.io/gorm"
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
	var questionModel bankQuestionModel
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		questionModel = bankQuestionModel{
			BankId:        params.BankId,
			QuestionText:  params.QuestionText,
			QuestionType:  params.QuestionType,
			DefaultPoints: params.DefaultPoints,
		}

		if err := tx.Create(&questionModel).Error; err != nil {
			return fmt.Errorf("create question: %w", err)
		}

		if len(params.Answers) == 0 {
			return nil
		}

		answerModels := make([]bankAnswerModel, 0, len(params.Answers))
		for _, a := range params.Answers {
			answerModels = append(answerModels, bankAnswerModel{
				QuestionId: questionModel.ID,
				AnswerText: a.AnswerText,
				IsCorrect:  a.IsCorrect,
			})
		}

		if err := tx.Create(&answerModels).Error; err != nil {
			return fmt.Errorf("create answers: %w", err)
		}

		return nil
	})
	if err != nil {
		return domain.BankQuestion{}, err
	}

	var answerModels []bankAnswerModel
	if err := r.db.WithContext(ctx).Where("question_id = ?", questionModel.ID).Find(&answerModels).Error; err != nil {
		return domain.BankQuestion{}, fmt.Errorf("load answers after create: %w", err)
	}

	return questionModel.toDomain(answerModels), nil
}

func (r *PostgresRepo) Update(ctx context.Context, params bankquestions.UpdateRepoParams) (domain.BankQuestion, error) {
	questionModel := bankQuestionModel{
		ID: params.ID,
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&questionModel).Updates(map[string]interface{}{
			"bank_id":        params.BankId,
			"question_text":  params.QuestionText,
			"question_type":  params.QuestionType,
			"default_points": params.DefaultPoints,
		}).Error; err != nil {
			return fmt.Errorf("update question: %w", err)
		}

		if err := tx.Where("question_id = ?", params.ID).Delete(&bankAnswerModel{}).Error; err != nil {
			return fmt.Errorf("delete old answers: %w", err)
		}

		if len(params.Answers) == 0 {
			return nil
		}

		answerModels := make([]bankAnswerModel, 0, len(params.Answers))
		for _, a := range params.Answers {
			answerModels = append(answerModels, bankAnswerModel{
				QuestionId: params.ID,
				AnswerText: a.AnswerText,
				IsCorrect:  a.IsCorrect,
			})
		}

		if err := tx.Create(&answerModels).Error; err != nil {
			return fmt.Errorf("create answers: %w", err)
		}

		return nil
	})
	if err != nil {
		return domain.BankQuestion{}, err
	}

	var answerModels []bankAnswerModel
	if err := r.db.WithContext(ctx).Where("question_id = ?", params.ID).Find(&answerModels).Error; err != nil {
		return domain.BankQuestion{}, fmt.Errorf("load answers after update: %w", err)
	}

	return questionModel.toDomain(answerModels), nil
}

func (r *PostgresRepo) Delete(ctx context.Context, id int64) error {
	question := bankQuestionModel{
		ID: id,
	}

	if err := r.db.WithContext(ctx).Delete(&question).Error; err != nil {
		return fmt.Errorf("delete question: %w", err)
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
