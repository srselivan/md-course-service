package lectures

import (
	"context"
	"course-service/internal/domain"
	"course-service/internal/repository"
	"course-service/internal/services/lectures"
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

func (r *PostgresRepo) Create(ctx context.Context, params lectures.CreateRepoParams) (domain.Lecture, error) {
	lecture := lectureModel{
		ItemId:            params.ItemId,
		Content:           params.Content,
		VideoURL:          params.VideoURL,
		ReadingTimeMinute: params.ReadingTimeMinute,
	}

	result := r.db.WithContext(ctx).Create(&lecture)

	if err := result.Error; err != nil {
		return domain.Lecture{}, fmt.Errorf("create: %w", err)
	}
	return lecture.toDomain(), nil
}

func (r *PostgresRepo) Update(ctx context.Context, params lectures.UpdateRepoParams) (domain.Lecture, error) {
	lecture := lectureModel{
		ItemId: params.ItemId,
	}

	result := r.db.WithContext(ctx).Model(&lecture).Updates(
		map[string]interface{}{
			"content":              params.Content,
			"video_url":            params.VideoURL,
			"reading_time_minutes": params.ReadingTimeMinute,
		},
	)

	if err := result.Error; err != nil {
		return domain.Lecture{}, fmt.Errorf("update: %w", err)
	}
	return lecture.toDomain(), nil
}

func (r *PostgresRepo) Delete(ctx context.Context, itemId int64) error {
	lecture := lectureModel{
		ItemId: itemId,
	}

	result := r.db.WithContext(ctx).Delete(&lecture)

	if err := result.Error; err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (r *PostgresRepo) Get(ctx context.Context, itemId int64) (domain.Lecture, error) {
	lecture := lectureModel{
		ItemId: itemId,
	}

	result := r.db.WithContext(ctx).First(&lecture)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Lecture{}, repository.ErrNotFound
		}
		return domain.Lecture{}, fmt.Errorf("get: %w", err)
	}
	return lecture.toDomain(), nil
}

func (r *PostgresRepo) GetList(ctx context.Context, params lectures.GetListRepoParams) ([]domain.Lecture, error) {
	var models []lectureModel

	result := r.db.WithContext(ctx).Find(&models)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("get list: %w", err)
	}

	lecturesList := make([]domain.Lecture, 0, len(models))
	for _, m := range models {
		lecturesList = append(lecturesList, m.toDomain())
	}

	return lecturesList, nil
}
