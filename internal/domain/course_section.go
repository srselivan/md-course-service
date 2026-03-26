package domain

type CourseSection struct {
	ID        int64
	CourseId  int64
	ParentId  *int64
	Title     string
	SortOrder int
}
