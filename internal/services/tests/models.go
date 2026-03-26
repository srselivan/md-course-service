package tests

import (
	"time"
)

type (
	CreateRepoParams struct {
		ItemId             int32
		Description        string
		EffectiveFrom      time.Time
		EffectiveTill      time.Time
		DurationSec        string
		MaxAttempts        int16
		QuestionsTotal     int16
		GenerationSettings string
	}
	UpdateRepoParams struct {
		ID                 int64
		ItemId             int32
		Description        string
		EffectiveFrom      time.Time
		EffectiveTill      time.Time
		DurationSec        string
		MaxAttempts        int16
		QuestionsTotal     int16
		GenerationSettings string
	}
	GetListRepoParams struct{}
)

type (
	CreateServiceParams struct {
		ItemId             int32
		Description        string
		EffectiveFrom      time.Time
		EffectiveTill      time.Time
		DurationSec        string
		MaxAttempts        int16
		QuestionsTotal     int16
		GenerationSettings string
	}
	UpdateServiceParams struct {
		ID                 int64
		ItemId             int32
		Description        string
		EffectiveFrom      time.Time
		EffectiveTill      time.Time
		DurationSec        string
		MaxAttempts        int16
		QuestionsTotal     int16
		GenerationSettings string
	}
	GetListServiceParams struct{}
)
