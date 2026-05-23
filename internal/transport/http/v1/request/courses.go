package request

type (
	// CreateCourse request body for POST /courses.
	CreateCourse struct {
		Title        string  `json:"title" example:"Advanced Calculus"`
		Description  *string `json:"description"`
		OwnerUserId  int64   `json:"ownerUserId" example:"42"`
		CoverImageId *string `json:"coverImageId" example:"550e8400-e29b-41d4-a716-446655440000"`
	}
	// UpdateCourse request body for PUT /courses/{id}.
	UpdateCourse struct {
		Title        string  `json:"title" example:"Advanced Calculus"`
		Description  *string `json:"description"`
		OwnerUserId  int64   `json:"ownerUserId" example:"42"`
		Status       int16   `json:"status" example:"1" enums:"0,1"` // 0 — Draft, 1 — Active
		CoverImageId *string `json:"coverImageId" example:"550e8400-e29b-41d4-a716-446655440000"`
	}
)

type (
	SetCourseListener struct {
		GroupIds []int64 `json:"groupIds"`
		UserIds  []int64 `json:"userIds"`
	}
	DeleteCourseListeners struct {
		GroupIds []int64 `json:"groupIds"`
		UserIds  []int64 `json:"userIds"`
	}
)
