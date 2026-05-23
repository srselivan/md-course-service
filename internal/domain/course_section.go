package domain

type CourseSection struct {
	ID        int64  `json:"id"`
	CourseId  int64  `json:"courseId"`
	ParentId  *int64 `json:"parentId"`
	Title     string `json:"title"`
	SortOrder int    `json:"sortOrder"`
}
