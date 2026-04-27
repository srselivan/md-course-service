package domain

import "time"

type Course struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	Description *string     `json:"description"`
	OwnerUserId int64       `json:"ownerUserId"`
	Status      int16       `json:"status"`
	Stats       CourseStats `json:"stats"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

type CourseStats struct {
	TotalStudents int64 `json:"totalStudents"`
}

type CourseWithItems struct {
	Course
	Sections []CourseSectionWithItems `json:"sections"`
}

type CourseSectionWithItems struct {
	CourseSection
	Items []CourseSectionItem `json:"items"`
}
