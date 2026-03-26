package domain

type Assignment struct {
	ItemId       int64  `json:"item_id"`
	Description  string `json:"description"`
	MaxScore     int    `json:"max_score"`
	DeadlineDays *int   `json:"deadline_days"`
}
