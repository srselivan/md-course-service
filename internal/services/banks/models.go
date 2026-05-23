package banks

type (
	CreateRepoParams struct {
		Title       string
		Description *string
		UserId      int64
	}
	UpdateRepoParams struct {
		ID          int64
		Title       string
		Description *string
		UserId      int64
	}
	GetListRepoParams struct {
		UserId int64
	}
)

type (
	CreateServiceParams struct {
		Title       string
		Description *string
		UserId      int64
	}
	UpdateServiceParams struct {
		ID          int64
		Title       string
		Description *string
		UserId      int64
	}
	GetListServiceParams struct {
		UserId int64
	}
)
