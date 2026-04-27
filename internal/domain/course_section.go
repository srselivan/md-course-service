package domain

type CourseSection struct {
	ID        int64  `json:"id"`
	CourseId  int64  `json:"course_id"`
	ParentId  *int64 `json:"parent_id"`
	Title     string `json:"title"`
	SortOrder int    `json:"sort_order"`
}
