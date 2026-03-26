package request

type (
	CreateCourse struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		OwnerUserId int64   `json:"owner_user_id"`
	}
	UpdateCourse struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		OwnerUserId int64   `json:"owner_user_id"`
	}
)
