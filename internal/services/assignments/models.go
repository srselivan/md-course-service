package assignments

type (
	CreateRepoParams struct {
		ItemId       int64
		Description  string
		MaxScore     int
		DeadlineDays *int
	}
	UpdateRepoParams struct {
		ItemId       int64
		Description  string
		MaxScore     int
		DeadlineDays *int
	}
	GetListRepoParams struct{}
)

type (
	CreateServiceParams struct {
		ItemId       int64
		Description  string
		MaxScore     int
		DeadlineDays *int
	}
	UpdateServiceParams struct {
		ItemId       int64
		Description  string
		MaxScore     int
		DeadlineDays *int
	}
	GetListServiceParams struct{}
)
