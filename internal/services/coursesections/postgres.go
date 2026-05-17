package coursesections

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

func (r *Repository) Create(ctx context.Context, params CreateRepoParams) (domain.CourseSection, error) {
	const query = `
		insert into course_section (course_id, parent_id, title, sort_order)
		values ($1, $2, $3, $4)
		returning id, course_id, parent_id, title, sort_order`

	var section courseSectionModel
	err := r.db.GetContext(ctx, &section, query,
		params.CourseId, params.ParentId, params.Title, params.SortOrder)
	if err != nil {
		return domain.CourseSection{}, fmt.Errorf("insert section: %w", err)
	}
	return section.toDomain(), nil
}

func (r *Repository) Update(ctx context.Context, params UpdateRepoParams) (domain.CourseSection, error) {
	const query = `
		update course_section
		set course_id = $2,
			parent_id = $3,
			title = $4,
			sort_order = $5
		where id = $1
		returning id, course_id, parent_id, title, sort_order`

	var section courseSectionModel
	err := r.db.GetContext(ctx, &section, query,
		params.ID, params.CourseId, params.ParentId, params.Title, params.SortOrder)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.CourseSection{}, repository.ErrNotFound
		}
		return domain.CourseSection{}, fmt.Errorf("update section: %w", err)
	}
	return section.toDomain(), nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `delete from course_section where id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete section: %w", err)
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

func (r *Repository) Get(ctx context.Context, id int64) (domain.CourseSection, error) {
	const query = `
		select id, course_id, parent_id, title, sort_order
		from course_section
		where id = $1`

	var section courseSectionModel
	err := r.db.GetContext(ctx, &section, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.CourseSection{}, repository.ErrNotFound
		}
		return domain.CourseSection{}, fmt.Errorf("get section: %w", err)
	}
	return section.toDomain(), nil
}

func (r *Repository) GetList(ctx context.Context, _ GetListRepoParams) ([]domain.CourseSection, error) {
	const query = `
		select id, course_id, parent_id, title, sort_order
		from course_section
		order by sort_order, id`

	var models []courseSectionModel
	err := r.db.SelectContext(ctx, &models, query)
	if err != nil {
		return nil, fmt.Errorf("list sections: %w", err)
	}

	sections := make([]domain.CourseSection, 0, len(models))
	for _, m := range models {
		sections = append(sections, m.toDomain())
	}
	return sections, nil
}
