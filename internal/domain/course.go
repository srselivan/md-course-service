package domain

import "time"

// Course is a course entity.
type Course struct {
	ID           int64       `json:"id" example:"1"`
	Title        string      `json:"title" example:"Advanced Calculus"`
	Description  *string     `json:"description"`
	OwnerUserId  int64       `json:"ownerUserId" example:"42"`
	Status       int16       `json:"status" example:"0" enums:"0,1"` // 0 — Draft, 1 — Active
	CoverImageId *string     `json:"coverImageId,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	Stats        CourseStats `json:"stats"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}

// CourseStats holds aggregated course metrics.
type CourseStats struct {
	TotalStudents int64 `json:"totalStudents" example:"45"`
}

// CoursesStatsResponse is the response for GET /courses/stats.
type CoursesStatsResponse struct {
	ActiveCoursesCount int64 `json:"activeCoursesCount" example:"3" db:"active_courses_count"`
	TotalStudents      int64 `json:"totalStudents" example:"120" db:"total_students"`
}

// CourseWithItems is a course with nested sections and section items.
type CourseWithItems struct {
	Course
	Sections []CourseSectionWithItems `json:"sections"`
}

// CourseSectionWithItems is a course section with its items.
type CourseSectionWithItems struct {
	CourseSection
	Items []CourseSectionItem `json:"items"`
}
