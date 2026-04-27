package banks

import (
	"context"
	"fmt"

	"course-service/internal/domain"
	"course-service/internal/services/banks"

	"gorm.io/gorm"
)

const banksTable = "bank"

type PostgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) Create(ctx context.Context, params banks.CreateRepoParams) (domain.Bank, error) {
	bank := bankModel{
		Title:       params.Title,
		Description: params.Description,
		CourseId:    params.CourseId,
	}

	result := r.db.WithContext(ctx).Table(banksTable).Create(&bank)

	if err := result.Error; err != nil {
		return domain.Bank{}, fmt.Errorf("create: %w", err)
	}
	return bank.toDomain(), nil
}

func (r *PostgresRepo) Update(ctx context.Context, params banks.UpdateRepoParams) (domain.Bank, error) {
	bank := bankModel{
		ID: params.ID,
	}

	result := r.db.WithContext(ctx).Table(banksTable).Model(&bank).Updates(
		map[string]interface{}{
			"title":       params.Title,
			"description": params.Description,
			"course_id":   params.CourseId,
		},
	)

	if err := result.Error; err != nil {
		return domain.Bank{}, fmt.Errorf("update: %w", err)
	}
	return bank.toDomain(), nil
}

func (r *PostgresRepo) Delete(ctx context.Context, id int64) error {
	bank := bankModel{
		ID: id,
	}

	result := r.db.WithContext(ctx).Table(banksTable).Delete(&bank)

	if err := result.Error; err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (r *PostgresRepo) GetList(ctx context.Context, params banks.GetListRepoParams) ([]domain.Bank, error) {
	var models []bankModel

	query := r.db.WithContext(ctx).Table(banksTable)
	if params.CourseId != nil {
		query = query.Where("course_id = ?", *params.CourseId)
	}

	result := query.Find(&models)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("get list: %w", err)
	}

	banksList := make([]domain.Bank, 0, len(models))
	for _, m := range models {
		banksList = append(banksList, m.toDomain())
	}

	return banksList, nil
}
