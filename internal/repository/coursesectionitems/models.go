package coursesectionitems

import "course-service/internal/domain"

type courseSectionItemModel struct {
	ID          int64  `gorm:"primaryKey"`
	SectionId   int64  `gorm:"column:section_id;not null"`
	ItemType    string `gorm:"column:item_type;type:varchar(50);not null"`
	Title       string `gorm:"column:title;type:varchar(255);not null"`
	SortOrder   int    `gorm:"column:sort_order"`
	IsPublished bool   `gorm:"column:is_published"`
}

func (m courseSectionItemModel) toDomain() domain.CourseSectionItem {
	return domain.CourseSectionItem{
		ID:          m.ID,
		SectionId:   m.SectionId,
		ItemType:    m.ItemType,
		Title:       m.Title,
		SortOrder:   m.SortOrder,
		IsPublished: m.IsPublished,
	}
}
