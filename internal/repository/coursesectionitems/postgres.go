package coursesectionitems

import (
	"context"
	"course-service/internal/domain"
	"course-service/internal/repository"
	"course-service/internal/services/coursesectionitems"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const courseSectionItemTable = "course_section_item"

type PostgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) Create(ctx context.Context, params coursesectionitems.CreateRepoParams) (domain.CourseSectionItem, error) {
	item := courseSectionItemModel{
		SectionId:   params.SectionId,
		ItemType:    params.ItemType,
		Title:       params.Title,
		SortOrder:   params.SortOrder,
		IsPublished: params.IsPublished,
	}

	result := r.db.WithContext(ctx).Table(courseSectionItemTable).Create(&item)

	if err := result.Error; err != nil {
		return domain.CourseSectionItem{}, fmt.Errorf("create: %w", err)
	}
	return item.toDomain(), nil
}

func (r *PostgresRepo) Update(ctx context.Context, params coursesectionitems.UpdateRepoParams) (domain.CourseSectionItem, error) {
	item := courseSectionItemModel{
		ID: params.ID,
	}

	result := r.db.WithContext(ctx).Table(courseSectionItemTable).Model(&item).Updates(
		map[string]interface{}{
			"section_id":   params.SectionId,
			"item_type":    params.ItemType,
			"title":        params.Title,
			"sort_order":   params.SortOrder,
			"is_published": params.IsPublished,
		},
	)

	if err := result.Error; err != nil {
		return domain.CourseSectionItem{}, fmt.Errorf("update: %w", err)
	}
	return item.toDomain(), nil
}

func (r *PostgresRepo) Delete(ctx context.Context, id int64) error {
	item := courseSectionItemModel{
		ID: id,
	}

	result := r.db.WithContext(ctx).Table(courseSectionItemTable).Delete(&item)

	if err := result.Error; err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (r *PostgresRepo) Get(ctx context.Context, id int64) (domain.CourseSectionItem, error) {
	item := courseSectionItemModel{
		ID: id,
	}

	result := r.db.WithContext(ctx).Table(courseSectionItemTable).First(&item)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.CourseSectionItem{}, repository.ErrNotFound
		}
		return domain.CourseSectionItem{}, fmt.Errorf("get: %w", err)
	}
	return item.toDomain(), nil
}

func (r *PostgresRepo) GetList(ctx context.Context, params coursesectionitems.GetListRepoParams) ([]domain.CourseSectionItem, error) {
	var models []courseSectionItemModel

	result := r.db.WithContext(ctx).Table(courseSectionItemTable).Find(&models)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("get list: %w", err)
	}

	items := make([]domain.CourseSectionItem, 0, len(models))
	for _, m := range models {
		items = append(items, m.toDomain())
	}

	return items, nil
}
