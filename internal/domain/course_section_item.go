package domain

type CourseSectionItem struct {
	ID          int64  `json:"id"`
	SectionId   int64  `json:"section_id"`
	ItemType    string `json:"item_type"`
	Title       string `json:"title"`
	SortOrder   int    `json:"sort_order"`
	IsPublished bool   `json:"is_published"`
}
