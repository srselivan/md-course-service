package lectures

type (
	CreateRepoParams struct {
		ItemId            int64
		Content           string
		VideoURL          string
		ReadingTimeMinute int
	}
	UpdateRepoParams struct {
		ItemId            int64
		Content           string
		VideoURL          string
		ReadingTimeMinute int
	}
	GetListRepoParams struct{}
)

type (
	CreateServiceParams struct {
		ItemId            int64
		Content           string
		VideoURL          string
		ReadingTimeMinute int
	}
	UpdateServiceParams struct {
		ItemId            int64
		Content           string
		VideoURL          string
		ReadingTimeMinute int
	}
	GetListServiceParams struct{}
)
