package coursesections

import (
	"context"
	"course-service/internal/domain"
	"course-service/internal/repository"
	"course-service/internal/services/coursesections"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const courseSectionTable = "course_section"

type PostgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) Create(ctx context.Context, params coursesections.CreateRepoParams) (domain.CourseSection, error) {
	course := courseSectionModel{
		CourseId:  params.CourseId,
		ParentId:  params.ParentId,
		Title:     params.Title,
		SortOrder: params.SortOrder,
	}

	result := r.db.WithContext(ctx).Table(courseSectionTable).Create(&course)

	if err := result.Error; err != nil {
		return domain.CourseSection{}, fmt.Errorf("create: %w", err)
	}
	return course.toDomain(), nil
}

func (r *PostgresRepo) Update(ctx context.Context, params coursesections.UpdateRepoParams) (domain.CourseSection, error) {
	course := courseSectionModel{
		ID: params.ID,
	}

	result := r.db.WithContext(ctx).Table(courseSectionTable).Model(&course).Updates(
		map[string]interface{}{
			"course_id":  params.CourseId,
			"parent_id":  params.ParentId,
			"title":      params.Title,
			"sort_order": params.SortOrder,
		},
	)

	if err := result.Error; err != nil {
		return domain.CourseSection{}, fmt.Errorf("update: %w", err)
	}
	return course.toDomain(), nil
}

func (r *PostgresRepo) Delete(ctx context.Context, id int64) error {
	course := courseSectionModel{
		ID: id,
	}

	result := r.db.WithContext(ctx).Table(courseSectionTable).Delete(&course)

	if err := result.Error; err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (r *PostgresRepo) Get(ctx context.Context, id int64) (domain.CourseSection, error) {
	course := courseSectionModel{
		ID: id,
	}

	result := r.db.WithContext(ctx).Table(courseSectionTable).First(&course)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.CourseSection{}, repository.ErrNotFound
		}
		return domain.CourseSection{}, fmt.Errorf("get: %w", err)
	}
	return course.toDomain(), nil
}

func (r *PostgresRepo) GetList(ctx context.Context, params coursesections.GetListRepoParams) ([]domain.CourseSection, error) {
	var models []courseSectionModel

	result := r.db.WithContext(ctx).Table(courseSectionTable).Find(&models)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("get list: %w", err)
	}

	sections := make([]domain.CourseSection, 0, len(models))
	for _, m := range models {
		sections = append(sections, m.toDomain())
	}

	return sections, nil
}
