package coursesections

import "course-service/internal/domain"

type courseSectionModel struct {
	ID        int64  `db:"id"`
	CourseId  int64  `db:"course_id"`
	ParentId  *int64 `db:"parent_id"`
	Title     string `db:"title"`
	SortOrder int    `db:"sort_order"`
}

func (m courseSectionModel) toDomain() domain.CourseSection {
	return domain.CourseSection{
		ID:        m.ID,
		CourseId:  m.CourseId,
		ParentId:  m.ParentId,
		Title:     m.Title,
		SortOrder: m.SortOrder,
	}
}
