package assignments

import "course-service/internal/domain"

type assignmentModel struct {
	ItemId       int64  `gorm:"column:item_id;primaryKey"`
	Description  string `gorm:"column:description;not null"`
	MaxScore     int    `gorm:"column:max_score"`
	DeadlineDays *int   `gorm:"column:deadline_days"`
}

func (m assignmentModel) toDomain() domain.Assignment {
	return domain.Assignment{
		ItemId:       m.ItemId,
		Description:  m.Description,
		MaxScore:     m.MaxScore,
		DeadlineDays: m.DeadlineDays,
	}
}
