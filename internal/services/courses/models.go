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
		Status      int16
	}
	GetListRepoParams struct {
		Status      *int16
		UserID      *int64
		OwnerUserId *int64
	}
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
		Status      int16
	}
	GetListServiceParams struct {
		Status      *int16
		UserID      *int64
		OwnerUserId *int64
	}
	GetWithAllItemsParams struct {
		Id     int64
		Limit  int64
		Offset int64
	}
)

type (
	SetListenerRepoParams struct {
		CourseId int64
		UserIds  []int64
	}
	DeleteListenerRepoParams struct {
		CourseId int64
		UserIds  []int64
	}
	GetListenersListRepoParams struct {
		CourseId int64
	}
	GetWithAllItemsRepoParams struct {
		Id     int64
		Limit  int64
		Offset int64
	}
)

type (
	SetListenerServiceParams struct {
		CourseId int64
		GroupIds []int64
		UserIds  []int64
	}
	DeleteListenerServiceParams struct {
		CourseId int64
		UserIds  []int64
		GroupIds []int64
	}
	GetListenersListServiceParams struct {
		CourseId int64
	}
)
