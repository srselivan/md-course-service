package assignments

import (
	"context"
	"course-service/internal/domain"
	"course-service/internal/repository"
	"course-service/internal/services/assignments"
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

func (r *PostgresRepo) Create(ctx context.Context, params assignments.CreateRepoParams) (domain.Assignment, error) {
	assignment := assignmentModel{
		ItemId:       params.ItemId,
		Description:  params.Description,
		MaxScore:     params.MaxScore,
		DeadlineDays: params.DeadlineDays,
	}

	result := r.db.WithContext(ctx).Create(&assignment)

	if err := result.Error; err != nil {
		return domain.Assignment{}, fmt.Errorf("create: %w", err)
	}
	return assignment.toDomain(), nil
}

func (r *PostgresRepo) Update(ctx context.Context, params assignments.UpdateRepoParams) (domain.Assignment, error) {
	assignment := assignmentModel{
		ItemId: params.ItemId,
	}

	result := r.db.WithContext(ctx).Model(&assignment).Updates(
		map[string]interface{}{
			"description":   params.Description,
			"max_score":     params.MaxScore,
			"deadline_days": params.DeadlineDays,
		},
	)

	if err := result.Error; err != nil {
		return domain.Assignment{}, fmt.Errorf("update: %w", err)
	}
	return assignment.toDomain(), nil
}

func (r *PostgresRepo) Delete(ctx context.Context, itemId int64) error {
	assignment := assignmentModel{
		ItemId: itemId,
	}

	result := r.db.WithContext(ctx).Delete(&assignment)

	if err := result.Error; err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (r *PostgresRepo) Get(ctx context.Context, itemId int64) (domain.Assignment, error) {
	assignment := assignmentModel{
		ItemId: itemId,
	}

	result := r.db.WithContext(ctx).First(&assignment)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Assignment{}, repository.ErrNotFound
		}
		return domain.Assignment{}, fmt.Errorf("get: %w", err)
	}
	return assignment.toDomain(), nil
}

func (r *PostgresRepo) GetList(ctx context.Context, params assignments.GetListRepoParams) ([]domain.Assignment, error) {
	var models []assignmentModel

	result := r.db.WithContext(ctx).Find(&models)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("get list: %w", err)
	}

	assignmentsList := make([]domain.Assignment, 0, len(models))
	for _, m := range models {
		assignmentsList = append(assignmentsList, m.toDomain())
	}

	return assignmentsList, nil
}
