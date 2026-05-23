package banks

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

func (r *Repository) Create(ctx context.Context, params CreateRepoParams) (domain.Bank, error) {
	const query = `
		insert into bank (title, description, user_id)
		values ($1, $2, $3)
		returning id, title, description, user_id, created_at, updated_at`

	var bank bankModel
	err := r.db.GetContext(ctx, &bank, query, params.Title, params.Description, params.UserId)
	if err != nil {
		return domain.Bank{}, fmt.Errorf("insert bank: %w", err)
	}
	return bank.toDomain(), nil
}

func (r *Repository) Update(ctx context.Context, params UpdateRepoParams) (domain.Bank, error) {
	const query = `
		update bank
		set title = $2,
			description = $3,
			updated_at = current_timestamp
		where id = $1 and user_id = $4
		returning id, title, description, user_id, created_at, updated_at`

	var bank bankModel
	err := r.db.GetContext(ctx, &bank, query, params.ID, params.Title, params.Description, params.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Bank{}, repository.ErrNotFound
		}
		return domain.Bank{}, fmt.Errorf("update bank: %w", err)
	}
	return bank.toDomain(), nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `delete from bank where id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete bank: %w", err)
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

func (r *Repository) GetList(ctx context.Context, params GetListRepoParams) ([]domain.Bank, error) {
	const query = `
		select b.id, b.title, b.description, b.user_id, b.created_at, b.updated_at,
			coalesce(q_count.cnt, 0) as questions_count
		from bank b
		left join (
			select bank_id, count(*) as cnt
			from bank_question
			group by bank_id
		) q_count on q_count.bank_id = b.id
		where b.user_id = $1
		order by b.id`

	var models []bankWithStatsModel
	if err := r.db.SelectContext(ctx, &models, query, params.UserId); err != nil {
		return nil, fmt.Errorf("list banks: %w", err)
	}

	banksList := make([]domain.Bank, 0, len(models))
	for _, m := range models {
		banksList = append(banksList, m.toDomain())
	}
	return banksList, nil
}
