package lectures

import "course-service/internal/domain"

type lectureModel struct {
	ItemId            int64  `gorm:"column:item_id;primaryKey"`
	Content           string `gorm:"column:content"`
	VideoURL          string `gorm:"column:video_url"`
	ReadingTimeMinute int    `gorm:"column:reading_time_minutes"`
}

func (m lectureModel) toDomain() domain.Lecture {
	return domain.Lecture{
		ItemId:            m.ItemId,
		Content:           m.Content,
		VideoURL:          m.VideoURL,
		ReadingTimeMinute: m.ReadingTimeMinute,
	}
}
