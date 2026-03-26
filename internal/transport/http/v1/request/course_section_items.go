package request

type (
	CreateCourseSectionItem struct {
		SectionId   int64  `json:"section_id"`
		ItemType    string `json:"item_type"`
		Title       string `json:"title"`
		SortOrder   int    `json:"sort_order"`
		IsPublished bool   `json:"is_published"`
	}
	UpdateCourseSectionItem struct {
		SectionId   int64  `json:"section_id"`
		ItemType    string `json:"item_type"`
		Title       string `json:"title"`
		SortOrder   int    `json:"sort_order"`
		IsPublished bool   `json:"is_published"`
	}
)
