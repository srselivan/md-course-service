package assignments

import "course-service/internal/domain"

type assignmentModel struct {
	ItemId       int64  `db:"item_id"`
	Description  string `db:"description"`
	MaxScore     int    `db:"max_score"`
	DeadlineDays *int   `db:"deadline_days"`
}

func (m assignmentModel) toDomain() domain.Assignment {
	return domain.Assignment{
		ItemId:       m.ItemId,
		Description:  m.Description,
		MaxScore:     m.MaxScore,
		DeadlineDays: m.DeadlineDays,
	}
}
