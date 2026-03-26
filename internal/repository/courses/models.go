package courses

import (
	"course-service/internal/domain"
	"time"
)

type courseModel struct {
	ID          int64     `gorm:"primaryKey"`
	Title       string    `gorm:"type:text;not null"`
	Description *string   `gorm:"type:text"`
	OwnerUserId int64     `gorm:"column:owner_user_id;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (c courseModel) toDomain() domain.Course {
	return domain.Course{
		ID:          c.ID,
		Title:       c.Title,
		Description: c.Description,
		OwnerUserId: c.OwnerUserId,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

type courseListenerModel struct {
	CourseId int64 `gorm:"column:course_id;primaryKey"`
	GroupId  int64 `gorm:"column:group_id;primaryKey"`
}
