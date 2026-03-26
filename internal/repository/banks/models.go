package banks

import (
	"course-service/internal/domain"
	"time"
)

type bankModel struct {
	ID          int64     `gorm:"primaryKey"`
	Title       string    `gorm:"type:text;not null"`
	Description *string   `gorm:"type:text"`
	CourseId    int64     `gorm:"column:course_id;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (b bankModel) toDomain() domain.Bank {
	return domain.Bank{
		ID:          b.ID,
		Title:       b.Title,
		Description: b.Description,
		CourseId:    b.CourseId,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}
