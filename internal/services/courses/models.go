package courses

type (
	CreateRepoParams struct {
		Title       string
		Description *string
		OwnerUserId int64
	}
	UpdateRepoParams struct {
		ID          int64
		Title       string
		Description *string
		OwnerUserId int64
	}
	GetListRepoParams struct{}
)

type (
	CreateServiceParams struct {
		Title       string
		Description *string
		OwnerUserId int64
	}
	UpdateServiceParams struct {
		ID          int64
		Title       string
		Description *string
		OwnerUserId int64
	}
	GetListServiceParams struct{}
)

type (
	SetListenerRepoParams struct {
		CourseId int64
		GroupId  int64
	}
	GetListenersListRepoParams struct {
		CourseId int64
	}
)

type (
	SetListenerServiceParams struct {
		CourseId int64
		GroupId  int64
	}
	GetListenersListServiceParams struct {
		CourseId int64
	}
)
