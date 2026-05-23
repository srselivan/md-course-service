package v1

import (
	"encoding/json"
	"net/http"

	"course-service/internal/services/courses"
	"course-service/internal/transport/http/utils"
	"course-service/internal/transport/http/v1/request"

	_ "course-service/internal/domain"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) NewCoursesRoutes(router fiber.Router) {
	coursesGroup := router.Group("courses")

	coursesGroup.Get("/", h.getCoursesList)
	coursesGroup.Get("/stats", h.getCoursesStats)
	coursesGroup.Post("/", h.createCourse)
	coursesGroup.Put("/:id", h.updateCourse)
	coursesGroup.Get("/:id", h.getCourse)

	coursesGroup.Post("/:course_id/listeners", h.setCourseListeners)
	coursesGroup.Delete("/:course_id/listeners", h.deleteCourseListeners)
	coursesGroup.Get("/:course_id/listeners", h.getCourseListenersList)

	coursesGroup.Get("/:course_id/with-all-items", h.getCourseWithAllItems)
}

// createCourse godoc
//
//	@Summary		Create course
//	@Description	Creates a course with status Draft (0). If coverImageId is set, publishes FileLoadedEvent to Kafka (file-topic).
//	@Tags			courses
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.CreateCourse	true	"Course"
//	@Success		201		{object}	domain.Course
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/courses [post]
func (h *Handler) createCourse(ctx fiber.Ctx) error {
	var req request.CreateCourse
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	course, err := h.coursesService.Create(ctx.Context(), courses.CreateServiceParams{
		Title:        req.Title,
		Description:  req.Description,
		OwnerUserId:  req.OwnerUserId,
		CoverImageId: req.CoverImageId,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(course)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusCreated).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// updateCourse godoc
//
//	@Summary		Update course
//	@Description	Updates course fields. When coverImageId changes, emits FileLoadedEvent and/or FileDeletedEvent to Kafka.
//	@Tags			courses
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Course ID"
//	@Param			body	body		request.UpdateCourse	true	"Course"
//	@Success		200		{object}	domain.Course
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/courses/{id} [put]
func (h *Handler) updateCourse(ctx fiber.Ctx) error {
	id, err := utils.GetInt64Param(ctx, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req request.UpdateCourse
	if err = json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	course, err := h.coursesService.Update(ctx.Context(), courses.UpdateServiceParams{
		ID:           id,
		Title:        req.Title,
		Description:  req.Description,
		OwnerUserId:  req.OwnerUserId,
		Status:       req.Status,
		CoverImageId: req.CoverImageId,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(course)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// deleteCourse godoc
//
//	@Summary	Delete course
//	@Tags		courses
//	@Param		id	path	int	true	"Course ID"
//	@Success	200
//	@Failure	400	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/courses/{id} [delete]
func (h *Handler) deleteCourse(ctx fiber.Ctx) error {
	id, err := utils.GetInt64Param(ctx, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = h.coursesService.Delete(ctx.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.SendStatus(http.StatusOK); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// getCourse godoc
//
//	@Summary		Get course by ID
//	@Tags			courses
//	@Produce		json
//	@Param			id	path		int	true	"Course ID"
//	@Success		200	{object}	domain.Course
//	@Failure		400	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/courses/{id} [get]
func (h *Handler) getCourse(ctx fiber.Ctx) error {
	id, err := utils.GetInt64Param(ctx, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	course, err := h.coursesService.Get(ctx.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(course)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// getCoursesList godoc
//
//	@Summary		List instructor courses
//	@Description	Returns courses owned by user_id. Filter: if the value is numeric, matches course id; otherwise matches title (case-insensitive substring). Student count is in stats.totalStudents.
//	@Tags			courses
//	@Produce		json
//	@Param			user_id	query		int		true	"Instructor user ID (until JWT auth)"
//	@Param			role	query		string	false	"User role (reserved for future auth)"
//	@Param			filter	query		string	false	"Filter by course id or title"
//	@Success		200		{array}		domain.Course
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/courses [get]
func (h *Handler) getCoursesList(ctx fiber.Ctx) error {
	userId, err := utils.GetInt64Query(ctx, "user_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	params := courses.GetListServiceParams{
		OwnerUserId: &userId,
	}

	filter := ctx.Query("filter")
	if filter != "" {
		params.Filter = &filter
	}

	coursesList, err := h.coursesService.GetList(ctx.Context(), params)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(coursesList)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// getCoursesStats godoc
//
//	@Summary		Course statistics
//	@Description	Returns count of active courses and total students enrolled on active courses for the instructor.
//	@Tags			courses
//	@Produce		json
//	@Param			user_id	query		int		true	"Instructor user ID (until JWT auth)"
//	@Param			role	query		string	false	"User role (reserved for future auth)"
//	@Success		200		{object}	domain.CoursesStatsResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/courses/stats [get]
func (h *Handler) getCoursesStats(ctx fiber.Ctx) error {
	userId, err := utils.GetInt64Query(ctx, "user_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	stats, err := h.coursesService.GetStats(ctx.Context(), courses.GetStatsServiceParams{
		OwnerUserId: userId,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(stats)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// getCourseListenersList godoc
//
//	@Summary	List course listeners
//	@Tags		courses
//	@Produce	json
//	@Param		course_id	path		int	true	"Course ID"
//	@Success	200			{array}		int64
//	@Failure	400			{object}	ErrorResponse
//	@Failure	500			{object}	ErrorResponse
//	@Router		/courses/{course_id}/listeners [get]
func (h *Handler) getCourseListenersList(ctx fiber.Ctx) error {
	courseID, err := utils.GetInt64Param(ctx, "course_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	listenersList, err := h.courseListenersService.GetListenersList(
		ctx.Context(),
		courses.GetListenersListServiceParams{CourseId: courseID},
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(listenersList)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// deleteCourseListeners godoc
//
//	@Summary	Remove course listeners
//	@Tags		courses
//	@Accept		json
//	@Param		course_id	path	int								true	"Course ID"
//	@Param		body		body	request.DeleteCourseListeners	true	"Listeners"
//	@Success	200
//	@Failure	400	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/courses/{course_id}/listeners [delete]
func (h *Handler) deleteCourseListeners(ctx fiber.Ctx) error {
	courseID, err := utils.GetInt64Param(ctx, "course_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req request.DeleteCourseListeners
	if err = json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = h.courseListenersService.DeleteListener(
		ctx.Context(),
		courses.DeleteListenerServiceParams{
			CourseId: courseID,
			UserIds:  req.UserIds,
			GroupIds: req.GroupIds,
		},
	); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(nil); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// setCourseListeners godoc
//
//	@Summary	Add course listeners
//	@Tags		courses
//	@Accept		json
//	@Param		course_id	path	int							true	"Course ID"
//	@Param		body		body	request.SetCourseListener	true	"Listeners"
//	@Success	200
//	@Failure	400	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/courses/{course_id}/listeners [post]
func (h *Handler) setCourseListeners(ctx fiber.Ctx) error {
	courseID, err := utils.GetInt64Param(ctx, "course_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req request.SetCourseListener
	if err = json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = h.courseListenersService.SetListener(
		ctx.Context(),
		courses.SetListenerServiceParams{
			CourseId: courseID,
			GroupIds: req.GroupIds,
			UserIds:  req.UserIds,
		},
	); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(nil); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// getCourseWithAllItems godoc
//
//	@Summary		Get course with all items
//	@Description	Returns course metadata (including coverImageId) and nested sections with items (item_type, item_id for file/test/assignment reference).
//	@Tags			courses
//	@Produce		json
//	@Param			course_id	path		int	true	"Course ID"
//	@Param			limit		query		int	true	"Page size for joined rows"
//	@Param			offset		query		int	true	"Page offset for joined rows"
//	@Success		200			{object}	domain.CourseWithItems
//	@Failure		400			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/courses/{course_id}/with-all-items [get]
func (h *Handler) getCourseWithAllItems(ctx fiber.Ctx) error {
	courseID, err := utils.GetInt64Param(ctx, "course_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	limit, err := utils.GetInt64Query(ctx, "limit")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	offset, err := utils.GetInt64Query(ctx, "offset")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	courseWithItems, err := h.coursesService.GetWithAllItems(ctx, courses.GetWithAllItemsParams{
		Id:     courseID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(courseWithItems)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}
