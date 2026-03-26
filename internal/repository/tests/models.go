package tests

import (
	"course-service/internal/domain"
	"time"
)

type testModel struct {
	ID                 int64     `gorm:"primaryKey"`
	ItemId             int32     `gorm:"column:item_id;not null"`
	Description        string    `gorm:"type:text;not null"`
	EffectiveFrom      time.Time `gorm:"column:effective_from;not null"`
	EffectiveTill      time.Time `gorm:"column:effective_till;not null"`
	DurationSec        string    `gorm:"column:duration_sec;type:interval;not null"`
	MaxAttempts        int16     `gorm:"column:max_attempts;default:1"`
	QuestionsTotal     int16     `gorm:"column:questions_total;not null"`
	GenerationSettings []byte    `gorm:"column:generation_settings;type:jsonb;not null"`
}

func (t testModel) toDomain() domain.Test {
	return domain.Test{
		ID:                 t.ID,
		ItemId:             t.ItemId,
		Description:        t.Description,
		EffectiveFrom:      t.EffectiveFrom,
		EffectiveTill:      t.EffectiveTill,
		DurationSec:        t.DurationSec,
		MaxAttempts:        t.MaxAttempts,
		QuestionsTotal:     t.QuestionsTotal,
		GenerationSettings: t.GenerationSettings,
	}
}
