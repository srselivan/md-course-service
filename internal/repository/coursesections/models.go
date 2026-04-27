package coursesections

import "course-service/internal/domain"

type courseSectionModel struct {
	ID        int64
	CourseId  int64
	ParentId  *int64
	Title     string
	SortOrder int
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
