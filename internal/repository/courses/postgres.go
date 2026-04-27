package courses

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"gorm.io/gorm/clause"

	"course-service/internal/domain"
	"course-service/internal/repository"
	"course-service/internal/services/courses"

	"gorm.io/gorm"
)

type PostgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepo(db *gorm.DB) *PostgresRepo {
	return &PostgresRepo{
		db: db,
	}
}

func (r *PostgresRepo) Create(ctx context.Context, params courses.CreateRepoParams) (domain.Course, error) {
	course := courseModel{
		Title:       params.Title,
		Description: params.Description,
		OwnerUserId: params.OwnerUserId,
	}

	result := r.db.WithContext(ctx).Table("course").Create(&course)

	if err := result.Error; err != nil {
		return domain.Course{}, fmt.Errorf("create: %w", err)
	}
	return course.toDomain(), nil
}

func (r *PostgresRepo) Update(ctx context.Context, params courses.UpdateRepoParams) (domain.Course, error) {
	course := courseModel{
		ID: params.ID,
	}

	result := r.db.WithContext(ctx).Table("course").Model(&course).Updates(
		map[string]interface{}{
			"title":         params.Title,
			"description":   params.Description,
			"owner_user_id": params.OwnerUserId,
			"status":        params.Status,
		},
	)

	if err := result.Error; err != nil {
		return domain.Course{}, fmt.Errorf("update: %w", err)
	}
	return course.toDomain(), nil
}

func (r *PostgresRepo) Delete(ctx context.Context, id int64) error {
	course := courseModel{
		ID: id,
	}

	result := r.db.WithContext(ctx).Table("course").Delete(&course)

	if err := result.Error; err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (r *PostgresRepo) Get(ctx context.Context, id int64) (domain.Course, error) {
	course := courseModel{
		ID: id,
	}

	result := r.db.WithContext(ctx).Table("course").First(&course)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Course{}, repository.ErrNotFound
		}
		return domain.Course{}, fmt.Errorf("get: %w", err)
	}

	const statsQuery = `
select count(*)
from course_listener
where course_listener.course_id = ?
`

	var listenersCount int64
	result = r.db.WithContext(ctx).Table("course").Raw(statsQuery, id).Scan(&listenersCount)
	if err := result.Error; err != nil {
		return domain.Course{}, fmt.Errorf("get: %w", err)
	}

	courseDomain := course.toDomain()
	courseDomain.Stats = domain.CourseStats{
		TotalStudents: listenersCount,
	}
	return courseDomain, nil
}

func (r *PostgresRepo) GetList(ctx context.Context, params courses.GetListRepoParams) ([]domain.Course, error) {
	var models []courseModel

	tx := r.db.WithContext(ctx).Table("course")

	if params.Status != nil {
		tx.Where("status = ?", *params.Status)
	}
	if params.OwnerUserId != nil {
		tx.Where("owner_user_id = ?", *params.OwnerUserId)
	}
	if params.UserID != nil {
		tx.
			Joins("JOIN course_listener cs ON cs.course_id = course.id").
			Where("cs.user_id = ?", *params.UserID)
	}

	result := tx.Find(&models)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("get list: %w", err)
	}

	coursesList := make([]domain.Course, 0, len(models))
	for _, m := range models {
		coursesList = append(coursesList, m.toDomain())
	}

	return coursesList, nil
}

func (r *PostgresRepo) GetWithAllItems(ctx context.Context, params courses.GetWithAllItemsRepoParams) (domain.CourseWithItems, error) {
	const query = `
select
	c.id, 
	c.title, 
	c.description, 
	c.owner_user_id, 
	c.status,
	c.created_at, 
	c.updated_at,
	cs.id cs_id, 
	cs.course_id cs_course_id, 
	cs.parent_id cs_parent_id, 
	cs.title cs_title, 
	cs.sort_order cs_sort_order,
	csi.id csi_id, 
	csi.section_id csi_section_id, 
	csi.item_type csi_item_type, 
	csi.title csi_title, 
	csi.sort_order csi_sort_order, 
	csi.is_published csi_is_published
from course c
left join course_section cs on cs.course_id = c.id
left join course_section_item csi on csi.section_id = cs.id
where c.id = ?
order by cs.sort_order, csi.sort_order
limit ?
offset ?
`

	rows, err := r.db.WithContext(ctx).Raw(query, params.Id, params.Limit, params.Offset).Rows()
	if err != nil {
		return domain.CourseWithItems{}, fmt.Errorf("raw query: %w", err)
	}
	defer rows.Close()

	var coursesList []courseWithAllItemsModel
	if rows.Next() {
		if err = r.db.WithContext(ctx).ScanRows(rows, &coursesList); err != nil {
			return domain.CourseWithItems{}, fmt.Errorf("scan rows: %w", err)
		}
	}

	return courseWithAllItemsModelListToDomain(coursesList), nil
}

func (r *PostgresRepo) SetListener(ctx context.Context, params courses.SetListenerRepoParams) error {
	listeners := lo.Map(params.UserIds, func(userId int64, _ int) courseListenerModel {
		return courseListenerModel{
			CourseId: params.CourseId,
			UserId:   userId,
		}
	})

	result := r.db.
		WithContext(ctx).
		Table("course_listener").
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&listeners, 100)

	if err := result.Error; err != nil {
		return fmt.Errorf("set listener: %w", err)
	}
	return nil
}

func (r *PostgresRepo) DeleteListener(ctx context.Context, params courses.DeleteListenerRepoParams) error {
	listener := courseListenerModel{
		CourseId: params.CourseId,
	}

	result := r.db.WithContext(ctx).Table("course_listener").Delete(&listener, params.UserIds)
	if err := result.Error; err != nil {
		return fmt.Errorf("delete listener: %w", err)
	}

	return nil
}

func (r *PostgresRepo) GetListenersList(ctx context.Context, params courses.GetListenersListRepoParams) ([]int64, error) {
	var listeners []courseListenerModel

	result := r.db.WithContext(ctx).Table("course_listener").Where("course_id = ?", params.CourseId).Find(&listeners)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("get listeners: %w", err)
	}

	userIds := make([]int64, 0, len(listeners))
	for _, l := range listeners {
		userIds = append(userIds, l.UserId)
	}

	return userIds, nil
}
