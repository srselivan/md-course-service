package banks

import (
	"time"

	"course-service/internal/domain"
)

type bankModel struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	Description *string   `db:"description"`
	CourseId    int64     `db:"course_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (m bankModel) toDomain() domain.Bank {
	return domain.Bank{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		CourseId:    m.CourseId,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
