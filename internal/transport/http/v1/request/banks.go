package request

type (
	CreateBank struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		UserId      int64   `json:"userId"`
	}
	UpdateBank struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		UserId      int64   `json:"userId"`
	}
)
