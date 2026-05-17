package courses

import (
	"database/sql"
	"time"

	"course-service/internal/domain"
)

type courseModel struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	Description *string   `db:"description"`
	OwnerUserId int64     `db:"owner_user_id"`
	Status      int16     `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (m courseModel) toDomain(stats domain.CourseStats) domain.Course {
	return domain.Course{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		OwnerUserId: m.OwnerUserId,
		Status:      m.Status,
		Stats:       stats,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

type courseWithAllItemsRow struct {
	ID             int64          `db:"id"`
	Title          string         `db:"title"`
	Description    *string        `db:"description"`
	OwnerUserId    int64          `db:"owner_user_id"`
	Status         int16          `db:"status"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
	CsID           sql.NullInt64  `db:"cs_id"`
	CsCourseId     sql.NullInt64  `db:"cs_course_id"`
	CsParentId     sql.NullInt64  `db:"cs_parent_id"`
	CsTitle        sql.NullString `db:"cs_title"`
	CsSortOrder    sql.NullInt64  `db:"cs_sort_order"`
	CsiID          sql.NullInt64  `db:"csi_id"`
	CsiSectionId   sql.NullInt64  `db:"csi_section_id"`
	CsiItemType    sql.NullString `db:"csi_item_type"`
	CsiTitle       sql.NullString `db:"csi_title"`
	CsiSortOrder   sql.NullInt64  `db:"csi_sort_order"`
	CsiIsPublished sql.NullBool   `db:"csi_is_published"`
}

func courseWithAllItemsRowsToDomain(rows []courseWithAllItemsRow) domain.CourseWithItems {
	first := rows[0]
	course := domain.CourseWithItems{
		Course: domain.Course{
			ID:          first.ID,
			Title:       first.Title,
			Description: first.Description,
			OwnerUserId: first.OwnerUserId,
			Status:      first.Status,
			CreatedAt:   first.CreatedAt,
			UpdatedAt:   first.UpdatedAt,
		},
		Sections: make([]domain.CourseSectionWithItems, 0),
	}

	sectionIdx := make(map[int64]int)
	for _, row := range rows {
		if !row.CsID.Valid {
			continue
		}

		sIdx, ok := sectionIdx[row.CsID.Int64]
		if !ok {
			var parentID *int64
			if row.CsParentId.Valid {
				parentID = &row.CsParentId.Int64
			}
			course.Sections = append(course.Sections, domain.CourseSectionWithItems{
				CourseSection: domain.CourseSection{
					ID:        row.CsID.Int64,
					CourseId:  row.CsCourseId.Int64,
					ParentId:  parentID,
					Title:     row.CsTitle.String,
					SortOrder: int(row.CsSortOrder.Int64),
				},
				Items: make([]domain.CourseSectionItem, 0),
			})
			sIdx = len(course.Sections) - 1
			sectionIdx[row.CsID.Int64] = sIdx
		}

		if !row.CsiID.Valid {
			continue
		}

		course.Sections[sIdx].Items = append(course.Sections[sIdx].Items, domain.CourseSectionItem{
			ID:          row.CsiID.Int64,
			SectionId:   row.CsiSectionId.Int64,
			ItemType:    row.CsiItemType.String,
			Title:       row.CsiTitle.String,
			SortOrder:   int(row.CsiSortOrder.Int64),
			IsPublished: row.CsiIsPublished.Bool,
		})
	}

	return course
}
