package domain

// CourseSectionItem is an element inside a course section (file, test, or assignment).
type CourseSectionItem struct {
	ID          int64  `json:"id" example:"10"`
	SectionId   int64  `json:"sectionId" example:"5"`
	ItemType    string `json:"itemType" example:"test" enums:"lecture,assignment,test"`
	ItemId      *int64 `json:"itemId,omitempty" example:"100"`
	Title       string `json:"title" example:"Midterm"`
	SortOrder   int    `json:"sortOrder" example:"1"`
	IsPublished bool   `json:"isPublished" example:"true"`
}
