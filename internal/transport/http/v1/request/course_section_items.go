package request

type (
	CreateCourseSectionItem struct {
		SectionId   int64  `json:"sectionId"`
		ItemType    string `json:"itemType"`
		Title       string `json:"title"`
		SortOrder   int    `json:"sortOrder"`
		IsPublished bool   `json:"isPublished"`
	}
	UpdateCourseSectionItem struct {
		SectionId   int64  `json:"sectionId"`
		ItemType    string `json:"itemType"`
		Title       string `json:"title"`
		SortOrder   int    `json:"sortOrder"`
		IsPublished bool   `json:"isPublished"`
	}
)
