package coursesectionitems

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"course-service/internal/domain"
	"course-service/internal/repository"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, params CreateRepoParams) (domain.CourseSectionItem, error) {
	const query = `
		insert into course_section_item (section_id, item_type, title, sort_order, is_published)
		values ($1, $2, $3, $4, $5)
		returning id, section_id, item_type, title, sort_order, is_published`

	var item courseSectionItemModel
	err := r.db.GetContext(ctx, &item, query,
		params.SectionId, params.ItemType, params.Title, params.SortOrder, params.IsPublished)
	if err != nil {
		return domain.CourseSectionItem{}, fmt.Errorf("insert item: %w", err)
	}
	return item.toDomain(), nil
}

func (r *Repository) Update(ctx context.Context, params UpdateRepoParams) (domain.CourseSectionItem, error) {
	const query = `
		update course_section_item
		set section_id = $2,
			item_type = $3,
			title = $4,
			sort_order = $5,
			is_published = $6
		where id = $1
		returning id, section_id, item_type, title, sort_order, is_published`

	var item courseSectionItemModel
	err := r.db.GetContext(ctx, &item, query,
		params.ID, params.SectionId, params.ItemType, params.Title, params.SortOrder, params.IsPublished)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.CourseSectionItem{}, repository.ErrNotFound
		}
		return domain.CourseSectionItem{}, fmt.Errorf("update item: %w", err)
	}
	return item.toDomain(), nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `delete from course_section_item where id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete rows affected: %w", err)
	}
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, id int64) (domain.CourseSectionItem, error) {
	const query = `
		select id, section_id, item_type, title, sort_order, is_published
		from course_section_item
		where id = $1`

	var item courseSectionItemModel
	err := r.db.GetContext(ctx, &item, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.CourseSectionItem{}, repository.ErrNotFound
		}
		return domain.CourseSectionItem{}, fmt.Errorf("get item: %w", err)
	}
	return item.toDomain(), nil
}

func (r *Repository) GetList(ctx context.Context, _ GetListRepoParams) ([]domain.CourseSectionItem, error) {
	const query = `
		select id, section_id, item_type, title, sort_order, is_published
		from course_section_item
		order by sort_order, id`

	var models []courseSectionItemModel
	err := r.db.SelectContext(ctx, &models, query)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}

	items := make([]domain.CourseSectionItem, 0, len(models))
	for _, m := range models {
		items = append(items, m.toDomain())
	}
	return items, nil
}
