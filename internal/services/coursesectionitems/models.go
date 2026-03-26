package coursesectionitems

type (
	CreateRepoParams struct {
		SectionId   int64
		ItemType    string
		Title       string
		SortOrder   int
		IsPublished bool
	}
	UpdateRepoParams struct {
		ID          int64
		SectionId   int64
		ItemType    string
		Title       string
		SortOrder   int
		IsPublished bool
	}
	GetListRepoParams struct{}
)

type (
	CreateServiceParams struct {
		SectionId   int64
		ItemType    string
		Title       string
		SortOrder   int
		IsPublished bool
	}
	UpdateServiceParams struct {
		ID          int64
		SectionId   int64
		ItemType    string
		Title       string
		SortOrder   int
		IsPublished bool
	}
	GetListServiceParams struct{}
)
