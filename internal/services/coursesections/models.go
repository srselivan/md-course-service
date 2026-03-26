package coursesections

type (
	CreateRepoParams struct {
		CourseId  int64
		ParentId  *int64
		Title     string
		SortOrder int
	}
	UpdateRepoParams struct {
		ID        int64
		CourseId  int64
		ParentId  *int64
		Title     string
		SortOrder int
	}
	GetListRepoParams struct{}
)

type (
	CreateServiceParams struct {
		CourseId  int64
		ParentId  *int64
		Title     string
		SortOrder int
	}
	UpdateServiceParams struct {
		ID        int64
		CourseId  int64
		ParentId  *int64
		Title     string
		SortOrder int
	}
	GetListServiceParams struct{}
)
