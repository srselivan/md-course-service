package coursesectionitems

import "course-service/internal/domain"

type courseSectionItemModel struct {
	ID          int64  `db:"id"`
	SectionId   int64  `db:"section_id"`
	ItemType    string `db:"item_type"`
	Title       string `db:"title"`
	SortOrder   int    `db:"sort_order"`
	IsPublished bool   `db:"is_published"`
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
