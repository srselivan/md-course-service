package courses

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"course-service/internal/domain"
	"course-service/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, params CreateRepoParams) (domain.Course, error) {
	const query = `
		insert into course (title, description, owner_user_id, cover_image_id)
		values ($1, $2, $3, $4)
		returning id, title, description, owner_user_id, status, cover_image_id, created_at, updated_at`

	var course courseModel
	err := r.db.GetContext(ctx, &course, query, params.Title, params.Description, params.OwnerUserId, params.CoverImageId)
	if err != nil {
		return domain.Course{}, fmt.Errorf("insert course: %w", err)
	}
	return course.toDomain(domain.CourseStats{}), nil
}

func (r *Repository) Update(ctx context.Context, params UpdateRepoParams) (domain.Course, error) {
	const query = `
		update course
		set title = $2,
			description = $3,
			owner_user_id = $4,
			status = $5,
			cover_image_id = $6,
			updated_at = current_timestamp
		where id = $1
		returning id, title, description, owner_user_id, status, cover_image_id, created_at, updated_at`

	var course courseModel
	err := r.db.GetContext(ctx, &course, query,
		params.ID, params.Title, params.Description, params.OwnerUserId, params.Status, params.CoverImageId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Course{}, repository.ErrNotFound
		}
		return domain.Course{}, fmt.Errorf("update course: %w", err)
	}
	return course.toDomain(domain.CourseStats{}), nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `delete from course where id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete course: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete rows affected: %w", err)
	}
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, id int64) (domain.Course, error) {
	const (
		getCourseQuery = `
			select id, title, description, owner_user_id, status, cover_image_id, created_at, updated_at
			from course
			where id = $1`
		countListenersQuery = `
			select count(*) from course_listener where course_id = $1`
	)

	var course courseModel
	err := r.db.GetContext(ctx, &course, getCourseQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Course{}, repository.ErrNotFound
		}
		return domain.Course{}, fmt.Errorf("get course: %w", err)
	}

	var listenersCount int64
	if err = r.db.GetContext(ctx, &listenersCount, countListenersQuery, id); err != nil {
		return domain.Course{}, fmt.Errorf("count listeners: %w", err)
	}

	return course.toDomain(domain.CourseStats{TotalStudents: listenersCount}), nil
}

func (r *Repository) GetList(ctx context.Context, params GetListRepoParams) ([]domain.Course, error) {
	const (
		listCoursesBaseQuery = `
			select c.id, c.title, c.description, c.owner_user_id, c.status, c.cover_image_id, c.created_at, c.updated_at,
				coalesce(cl_count.cnt, 0) as students_count
			from course c
			left join (select course_id, count(*) as cnt from course_listener group by course_id) cl_count on cl_count.course_id = c.id`
		listCoursesJoinListener = ` join course_listener cl on cl.course_id = c.id`
		listCoursesOrderBy      = ` order by c.id`
	)

	query := listCoursesBaseQuery
	var (
		conds  []string
		args   []any
		argNum int
	)

	if params.UserID != nil {
		query += listCoursesJoinListener
		argNum++
		conds = append(conds, fmt.Sprintf("cl.user_id = $%d", argNum))
		args = append(args, *params.UserID)
	}
	if params.Status != nil {
		argNum++
		conds = append(conds, fmt.Sprintf("c.status = $%d", argNum))
		args = append(args, *params.Status)
	}
	if params.OwnerUserId != nil {
		argNum++
		conds = append(conds, fmt.Sprintf("c.owner_user_id = $%d", argNum))
		args = append(args, *params.OwnerUserId)
	}
	if params.Filter != nil && *params.Filter != "" {
		_, parseErr := strconv.ParseInt(*params.Filter, 10, 64)
		if parseErr == nil {
			argNum++
			conds = append(conds, fmt.Sprintf("c.id = $%d", argNum))
			args = append(args, *params.Filter)
		} else {
			argNum++
			conds = append(conds, fmt.Sprintf("c.title ILIKE $%d", argNum))
			args = append(args, "%"+*params.Filter+"%")
		}
	}
	if len(conds) > 0 {
		query += " where " + strings.Join(conds, " and ")
	}
	query += listCoursesOrderBy

	var models []courseWithStatsModel
	if err := r.db.SelectContext(ctx, &models, query, args...); err != nil {
		return nil, fmt.Errorf("list courses: %w", err)
	}

	courses := make([]domain.Course, 0, len(models))
	for _, m := range models {
		courses = append(courses, m.toDomain())
	}
	return courses, nil
}

func (r *Repository) GetWithAllItems(ctx context.Context, params GetWithAllItemsRepoParams) (domain.CourseWithItems, error) {
	const query = `
		select
			c.id,
			c.title,
			c.description,
			c.owner_user_id,
			c.status,
			c.cover_image_id,
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
			csi.item_id csi_item_id,
			csi.title csi_title,
			csi.sort_order csi_sort_order,
			csi.is_published csi_is_published
		from course c
		left join course_section cs on cs.course_id = c.id
		left join course_section_item csi on csi.section_id = cs.id
		where c.id = $1
		order by cs.sort_order, csi.sort_order
		limit $2 offset $3`

	var rows []courseWithAllItemsRow
	err := r.db.SelectContext(ctx, &rows, query, params.Id, params.Limit, params.Offset)
	if err != nil {
		return domain.CourseWithItems{}, fmt.Errorf("get course with items: %w", err)
	}
	if len(rows) == 0 {
		return domain.CourseWithItems{}, repository.ErrNotFound
	}
	return courseWithAllItemsRowsToDomain(rows), nil
}

func (r *Repository) GetStats(ctx context.Context, params GetStatsRepoParams) (domain.CoursesStatsResponse, error) {
	const query = `
		select
			count(*) as active_courses_count,
			coalesce((select count(*) from course_listener cl
				join course c2 on c2.id = cl.course_id
				where c2.owner_user_id = $1 and c2.status = 1), 0) as total_students
		from course
		where owner_user_id = $1 and status = 1`

	var stats domain.CoursesStatsResponse
	err := r.db.GetContext(ctx, &stats, query, params.OwnerUserId)
	if err != nil {
		return domain.CoursesStatsResponse{}, fmt.Errorf("get stats: %w", err)
	}
	return stats, nil
}

func (r *Repository) SetListener(ctx context.Context, params SetListenerRepoParams) error {
	const query = `
		insert into course_listener (course_id, user_id)
		select $1, unnest($2::bigint[])
		on conflict do nothing`

	if len(params.UserIds) == 0 {
		return nil
	}
	_, err := r.db.ExecContext(ctx, query, params.CourseId, pq.Array(params.UserIds))
	if err != nil {
		return fmt.Errorf("set listeners: %w", err)
	}
	return nil
}

func (r *Repository) DeleteListener(ctx context.Context, params DeleteListenerRepoParams) error {
	const query = `
		delete from course_listener
		where course_id = $1 and user_id = any($2)`

	if len(params.UserIds) == 0 {
		return nil
	}
	_, err := r.db.ExecContext(ctx, query, params.CourseId, pq.Array(params.UserIds))
	if err != nil {
		return fmt.Errorf("delete listeners: %w", err)
	}
	return nil
}

func (r *Repository) GetListenersList(ctx context.Context, params GetListenersListRepoParams) ([]int64, error) {
	const query = `
		select user_id from course_listener where course_id = $1 order by user_id`

	var userIDs []int64
	err := r.db.SelectContext(ctx, &userIDs, query, params.CourseId)
	if err != nil {
		return nil, fmt.Errorf("list listeners: %w", err)
	}
	return userIDs, nil
}
