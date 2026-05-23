package request

type (
	CreateCourseSection struct {
		CourseId  int64  `json:"courseId"`
		ParentId  *int64 `json:"parentId"`
		Title     string `json:"title"`
		SortOrder int    `json:"sortOrder"`
	}
	UpdateCourseSection struct {
		CourseId  int64  `json:"courseId"`
		ParentId  *int64 `json:"parentId"`
		Title     string `json:"title"`
		SortOrder int    `json:"sortOrder"`
	}
)
