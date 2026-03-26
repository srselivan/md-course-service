package banks

type (
	CreateRepoParams struct {
		Title       string
		Description *string
		CourseId    int64
	}
	UpdateRepoParams struct {
		ID          int64
		Title       string
		Description *string
		CourseId    int64
	}
	GetListRepoParams struct {
		CourseId *int64
	}
)

type (
	CreateServiceParams struct {
		Title       string
		Description *string
		CourseId    int64
	}
	UpdateServiceParams struct {
		ID          int64
		Title       string
		Description *string
		CourseId    int64
	}
	GetListServiceParams struct {
		CourseId *int64
	}
)
