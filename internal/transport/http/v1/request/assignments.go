package request

type (
	CreateAssignment struct {
		ItemId       int64  `json:"itemId"`
		Description  string `json:"description"`
		MaxScore     int    `json:"maxScore"`
		DeadlineDays *int   `json:"deadlineDays"`
	}
	UpdateAssignment struct {
		ItemId       int64  `json:"itemId"`
		Description  string `json:"description"`
		MaxScore     int    `json:"maxScore"`
		DeadlineDays *int   `json:"deadlineDays"`
	}
)
