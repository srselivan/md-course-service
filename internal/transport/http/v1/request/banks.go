package request

type (
	CreateBank struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		CourseId    int64   `json:"course_id"`
	}
	UpdateBank struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		CourseId    int64   `json:"course_id"`
	}
)
