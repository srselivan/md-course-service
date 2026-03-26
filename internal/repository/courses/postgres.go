package courses

import (
	"context"
	"course-service/internal/domain"
	"course-service/internal/repository"
	"course-service/internal/services/courses"
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

func (r *PostgresRepo) Create(ctx context.Context, params courses.CreateRepoParams) (domain.Course, error) {
	course := courseModel{
		Title:       params.Title,
		Description: params.Description,
		OwnerUserId: params.OwnerUserId,
	}

	result := r.db.WithContext(ctx).Table("course").Create(&course)

	if err := result.Error; err != nil {
		return domain.Course{}, fmt.Errorf("create: %w", err)
	}
	return course.toDomain(), nil
}

func (r *PostgresRepo) Update(ctx context.Context, params courses.UpdateRepoParams) (domain.Course, error) {
	course := courseModel{
		ID: params.ID,
	}

	result := r.db.WithContext(ctx).Table("course").Model(&course).Updates(
		map[string]interface{}{
			"title":         params.Title,
			"description":   params.Description,
			"owner_user_id": params.OwnerUserId,
		},
	)

	if err := result.Error; err != nil {
		return domain.Course{}, fmt.Errorf("update: %w", err)
	}
	return course.toDomain(), nil
}

func (r *PostgresRepo) Delete(ctx context.Context, id int64) error {
	course := courseModel{
		ID: id,
	}

	result := r.db.WithContext(ctx).Table("course").Delete(&course)

	if err := result.Error; err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (r *PostgresRepo) Get(ctx context.Context, id int64) (domain.Course, error) {
	course := courseModel{
		ID: id,
	}

	result := r.db.WithContext(ctx).Table("course").First(&course)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Course{}, repository.ErrNotFound
		}
		return domain.Course{}, fmt.Errorf("get: %w", err)
	}
	return course.toDomain(), nil
}

func (r *PostgresRepo) GetList(ctx context.Context, params courses.GetListRepoParams) ([]domain.Course, error) {
	var models []courseModel

	result := r.db.WithContext(ctx).Table("course").Find(&models)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("get list: %w", err)
	}

	coursesList := make([]domain.Course, 0, len(models))
	for _, m := range models {
		coursesList = append(coursesList, m.toDomain())
	}

	return coursesList, nil
}

func (r *PostgresRepo) SetListener(ctx context.Context, params courses.SetListenerRepoParams) error {
	listener := courseListenerModel{
		CourseId: params.CourseId,
		GroupId:  params.GroupId,
	}

	result := r.db.WithContext(ctx).Table("course").Create(&listener)
	if err := result.Error; err != nil {
		return fmt.Errorf("set listener: %w", err)
	}

	return nil
}

func (r *PostgresRepo) GetListenersList(ctx context.Context, params courses.GetListenersListRepoParams) ([]int64, error) {
	var listeners []courseListenerModel

	result := r.db.WithContext(ctx).Table("course").Where("course_id = ?", params.CourseId).Find(&listeners)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("get listeners: %w", err)
	}

	groupIds := make([]int64, 0, len(listeners))
	for _, l := range listeners {
		groupIds = append(groupIds, l.GroupId)
	}

	return groupIds, nil
}
