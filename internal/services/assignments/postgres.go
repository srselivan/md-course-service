package assignments

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

func (r *Repository) Create(ctx context.Context, params CreateRepoParams) (domain.Assignment, error) {
	const query = `
		insert into assignments (item_id, description, max_score, deadline_days)
		values ($1, $2, $3, $4)
		returning item_id, description, max_score, deadline_days`

	var assignment assignmentModel
	err := r.db.GetContext(ctx, &assignment, query,
		params.ItemId, params.Description, params.MaxScore, params.DeadlineDays)
	if err != nil {
		return domain.Assignment{}, fmt.Errorf("insert assignment: %w", err)
	}
	return assignment.toDomain(), nil
}

func (r *Repository) Update(ctx context.Context, params UpdateRepoParams) (domain.Assignment, error) {
	const query = `
		update assignments
		set description = $2,
			max_score = $3,
			deadline_days = $4
		where item_id = $1
		returning item_id, description, max_score, deadline_days`

	var assignment assignmentModel
	err := r.db.GetContext(ctx, &assignment, query,
		params.ItemId, params.Description, params.MaxScore, params.DeadlineDays)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Assignment{}, repository.ErrNotFound
		}
		return domain.Assignment{}, fmt.Errorf("update assignment: %w", err)
	}
	return assignment.toDomain(), nil
}

func (r *Repository) Delete(ctx context.Context, itemID int64) error {
	const query = `delete from assignments where item_id = $1`

	result, err := r.db.ExecContext(ctx, query, itemID)
	if err != nil {
		return fmt.Errorf("delete assignment: %w", err)
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

func (r *Repository) Get(ctx context.Context, itemID int64) (domain.Assignment, error) {
	const query = `
		select item_id, description, max_score, deadline_days
		from assignments
		where item_id = $1`

	var assignment assignmentModel
	err := r.db.GetContext(ctx, &assignment, query, itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Assignment{}, repository.ErrNotFound
		}
		return domain.Assignment{}, fmt.Errorf("get assignment: %w", err)
	}
	return assignment.toDomain(), nil
}

func (r *Repository) GetList(ctx context.Context, _ GetListRepoParams) ([]domain.Assignment, error) {
	const query = `
		select item_id, description, max_score, deadline_days
		from assignments
		order by item_id`

	var models []assignmentModel
	err := r.db.SelectContext(ctx, &models, query)
	if err != nil {
		return nil, fmt.Errorf("list assignments: %w", err)
	}

	assignmentsList := make([]domain.Assignment, 0, len(models))
	for _, m := range models {
		assignmentsList = append(assignmentsList, m.toDomain())
	}
	return assignmentsList, nil
}
