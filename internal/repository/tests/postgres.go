package tests

import (
	"context"
	"course-service/internal/domain"
	"course-service/internal/repository"
	"course-service/internal/services/tests"
	"errors"
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

func (r *PostgresRepo) Create(ctx context.Context, params tests.CreateRepoParams) (domain.Test, error) {
	test := testModel{
		ItemId:             params.ItemId,
		Description:        params.Description,
		EffectiveFrom:      params.EffectiveFrom,
		EffectiveTill:      params.EffectiveTill,
		DurationSec:        params.DurationSec,
		MaxAttempts:        params.MaxAttempts,
		QuestionsTotal:     params.QuestionsTotal,
		GenerationSettings: params.GenerationSettings,
	}

	result := r.db.WithContext(ctx).Create(&test)

	if err := result.Error; err != nil {
		return domain.Test{}, fmt.Errorf("create: %w", err)
	}
	return test.toDomain(), nil
}

func (r *PostgresRepo) Update(ctx context.Context, params tests.UpdateRepoParams) (domain.Test, error) {
	test := testModel{
		ID: params.ID,
	}

	result := r.db.WithContext(ctx).Model(&test).Updates(
		map[string]interface{}{
			"item_id":             params.ItemId,
			"description":         params.Description,
			"effective_from":      params.EffectiveFrom,
			"effective_till":      params.EffectiveTill,
			"duration_sec":        params.DurationSec,
			"max_attempts":        params.MaxAttempts,
			"questions_total":     params.QuestionsTotal,
			"generation_settings": params.GenerationSettings,
		},
	)

	if err := result.Error; err != nil {
		return domain.Test{}, fmt.Errorf("update: %w", err)
	}
	return test.toDomain(), nil
}

func (r *PostgresRepo) Delete(ctx context.Context, id int64) error {
	test := testModel{
		ID: id,
	}

	result := r.db.WithContext(ctx).Delete(&test)

	if err := result.Error; err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (r *PostgresRepo) Get(ctx context.Context, id int64) (domain.Test, error) {
	test := testModel{
		ID: id,
	}

	result := r.db.WithContext(ctx).First(&test)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Test{}, repository.ErrNotFound
		}
		return domain.Test{}, fmt.Errorf("get: %w", err)
	}
	return test.toDomain(), nil
}

func (r *PostgresRepo) GetList(ctx context.Context, params tests.GetListRepoParams) ([]domain.Test, error) {
	var models []testModel

	result := r.db.WithContext(ctx).Find(&models)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("get list: %w", err)
	}

	testsList := make([]domain.Test, 0, len(models))
	for _, m := range models {
		testsList = append(testsList, m.toDomain())
	}

	return testsList, nil
}
