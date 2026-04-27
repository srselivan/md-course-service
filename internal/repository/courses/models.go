package courses

import (
	"time"

	"course-service/internal/domain"
)

type courseModel struct {
	ID          int64     `gorm:"primaryKey"`
	Title       string    `gorm:"type:text;not null"`
	Description *string   `gorm:"type:text"`
	OwnerUserId int64     `gorm:"column:owner_user_id;not null"`
	Status      int16     `gorm:"column:status"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (c courseModel) toDomain() domain.Course {
	return domain.Course{
		ID:          c.ID,
		Title:       c.Title,
		Description: c.Description,
		OwnerUserId: c.OwnerUserId,
		Status:      c.Status,
		Stats:       domain.CourseStats{},
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

type courseWithStatsModel struct {
	ID             int64     `gorm:"primaryKey"`
	Title          string    `gorm:"type:text;not null"`
	Description    *string   `gorm:"type:text"`
	OwnerUserId    int64     `gorm:"column:owner_user_id;not null"`
	Status         int16     `gorm:"column:status"`
	TotalListeners int64     `gorm:"column:total_listeners"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type courseListenerModel struct {
	CourseId int64 `gorm:"column:course_id"`
	UserId   int64 `gorm:"column:user_id"`
}

type courseWithAllItemsModel struct {
	ID             int64     `gorm:"column:id"`
	Title          string    `gorm:"column:title"`
	Description    *string   `gorm:"column:description"`
	OwnerUserId    int64     `gorm:"column:owner_user_id;not null"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	CsID           int64     `gorm:"column:cs_id"`
	CsCourseId     int64     `gorm:"column:cs_course_id"`
	CsParentId     *int64    `gorm:"column:cs_parent_id"`
	CsTitle        string    `gorm:"column:cs_title"`
	CsSortOrder    int       `gorm:"column:cs_sort_order"`
	CsiID          int64     `gorm:"column:csi_id"`
	CsiSectionId   int64     `gorm:"column:csi_section_id;not null"`
	CsiItemType    string    `gorm:"column:csi_item_type;type:varchar(50);not null"`
	CsiTitle       string    `gorm:"column:csi_title;type:varchar(255);not null"`
	CsiSortOrder   int       `gorm:"column:csi_sort_order"`
	CsiIsPublished bool      `gorm:"column:csi_is_published"`
}

func courseWithAllItemsModelListToDomain(rows []courseWithAllItemsModel) domain.CourseWithItems {
	first := rows[0]
	course := domain.CourseWithItems{
		Course: domain.Course{
			ID:          first.ID,
			Title:       first.Title,
			Description: first.Description,
			OwnerUserId: first.OwnerUserId,
			CreatedAt:   first.CreatedAt,
			UpdatedAt:   first.UpdatedAt,
		},
		Sections: make([]domain.CourseSectionWithItems, 0),
	}

	sectionIdxMap := make(map[int64]int)

	for _, row := range rows {
		if row.CsID == 0 {
			continue
		}

		sIdx, sectionExists := sectionIdxMap[row.CsID]
		if !sectionExists {
			section := domain.CourseSectionWithItems{
				CourseSection: domain.CourseSection{
					ID:        row.CsID,
					CourseId:  row.CsCourseId,
					ParentId:  row.CsParentId,
					Title:     row.CsTitle,
					SortOrder: row.CsSortOrder,
				},
				Items: make([]domain.CourseSectionItem, 0),
			}
			course.Sections = append(course.Sections, section)
			sIdx = len(course.Sections) - 1
			sectionIdxMap[row.CsID] = sIdx
		}

		if row.CsiID == 0 {
			continue
		}

		item := domain.CourseSectionItem{
			ID:          row.CsiID,
			SectionId:   row.CsiSectionId,
			ItemType:    row.CsiItemType,
			Title:       row.CsiTitle,
			SortOrder:   row.CsiSortOrder,
			IsPublished: row.CsiIsPublished,
		}
		course.Sections[sIdx].Items = append(course.Sections[sIdx].Items, item)
	}

	return course
}
