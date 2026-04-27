package request

type (
	CreateCourse struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		OwnerUserId int64   `json:"ownerUserId"`
	}
	UpdateCourse struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		OwnerUserId int64   `json:"ownerUserId"`
		Status      int16   `json:"status"`
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
