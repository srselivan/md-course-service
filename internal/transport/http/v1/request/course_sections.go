package request

type (
	CreateCourseSection struct {
		CourseId  int64  `json:"course_id"`
		ParentId  *int64 `json:"parent_id"`
		Title     string `json:"title"`
		SortOrder int    `json:"sort_order"`
	}
	UpdateCourseSection struct {
		CourseId  int64  `json:"course_id"`
		ParentId  *int64 `json:"parent_id"`
		Title     string `json:"title"`
		SortOrder int    `json:"sort_order"`
	}
)
