package banks

import (
	"time"

	"course-service/internal/domain"
)

type bankModel struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	Description *string   `db:"description"`
	UserId      int64     `db:"user_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (m bankModel) toDomain() domain.Bank {
	return domain.Bank{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		UserId:      m.UserId,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

type bankWithStatsModel struct {
	ID             int64     `db:"id"`
	Title          string    `db:"title"`
	Description    *string   `db:"description"`
	UserId         int64     `db:"user_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	QuestionsCount int64     `db:"questions_count"`
}

func (m bankWithStatsModel) toDomain() domain.Bank {
	return domain.Bank{
		ID:             m.ID,
		Title:          m.Title,
		Description:    m.Description,
		UserId:         m.UserId,
		QuestionsCount: m.QuestionsCount,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
