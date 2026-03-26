package coursesections

import "course-service/internal/domain"

/*
id         serial primary key,
course_id  int          not null references course (id) on delete cascade,
parent_id  int references course_section (id) on delete cascade,
title      varchar(255) not null,
sort_order int default 0
*/
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
