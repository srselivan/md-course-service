package domain

import "time"

type Test struct {
	ID                 int64                  `json:"id"`
	ItemId             int32                  `json:"item_id"`
	Title              string                 `json:"title"`
	Description        string                 `json:"description"`
	EffectiveFrom      time.Time              `json:"effective_from"`
	EffectiveTill      time.Time              `json:"effective_till"`
	DurationSec        time.Duration          `json:"duration_sec"`
	MaxAttempts        int16                  `json:"max_attempts"`
	QuestionsTotal     int16                  `json:"questions_total"`
	GenerationSettings TestGenerationSettings `json:"generation_settings"`
	BankIds            []int64                `json:"bank_ids"`
}

type TestGenerationSettings struct {
	QuestionsDistribution map[int]int16 `json:"questions_distribution"`
}

func (t Test) IsDistributionNeeded() bool {
	return len(t.GenerationSettings.QuestionsDistribution) != 0
}
