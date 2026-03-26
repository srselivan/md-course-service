package domain

type Lecture struct {
	ItemId            int64  `json:"item_id"`
	Content           string `json:"content"`
	VideoURL          string `json:"video_url"`
	ReadingTimeMinute int    `json:"reading_time_minutes"`
}
