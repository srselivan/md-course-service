package v1

import (
	"encoding/json"
	"net/http"

	"course-service/internal/services/coursesections"
	"course-service/internal/transport/http/utils"
	"course-service/internal/transport/http/v1/request"

	_ "course-service/internal/domain"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) NewCourseSectionsRoutes(router fiber.Router) {
	sectionsGroup := router.Group("course-sections")

	sectionsGroup.Post("/", h.createCourseSection)
	sectionsGroup.Put("/:id", h.updateCourseSection)
	sectionsGroup.Delete("/:id", h.deleteCourseSection)
	sectionsGroup.Get("/:id", h.getCourseSection)
	sectionsGroup.Get("/", h.getCourseSectionsList)
}

// createCourseSection godoc
//
//	@Summary	Create course section
//	@Tags		course-sections
//	@Accept		json
//	@Produce	json
//	@Param		body	body		request.CreateCourseSection	true	"Section"
//	@Success	201		{object}	domain.CourseSection
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/course-sections [post]
func (h *Handler) createCourseSection(ctx fiber.Ctx) error {
	var req request.CreateCourseSection
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	section, err := h.courseSectionsService.Create(ctx.Context(), coursesections.CreateServiceParams{
		CourseId:  req.CourseId,
		ParentId:  req.ParentId,
		Title:     req.Title,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(section)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusCreated).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// updateCourseSection godoc
//
//	@Summary	Update course section
//	@Tags		course-sections
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int							true	"Section ID"
//	@Param		body	body		request.UpdateCourseSection	true	"Section"
//	@Success	200		{object}	domain.CourseSection
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/course-sections/{id} [put]
func (h *Handler) updateCourseSection(ctx fiber.Ctx) error {
	id, err := utils.GetInt64Param(ctx, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req request.UpdateCourseSection
	if err = json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	section, err := h.courseSectionsService.Update(ctx.Context(), coursesections.UpdateServiceParams{
		ID:        id,
		CourseId:  req.CourseId,
		ParentId:  req.ParentId,
		Title:     req.Title,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(section)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// deleteCourseSection godoc
//
//	@Summary	Delete course section
//	@Tags		course-sections
//	@Param		id	path	int	true	"Section ID"
//	@Success	200
//	@Failure	400	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/course-sections/{id} [delete]
func (h *Handler) deleteCourseSection(ctx fiber.Ctx) error {
	id, err := utils.GetInt64Param(ctx, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = h.courseSectionsService.Delete(ctx.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.SendStatus(http.StatusOK); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// getCourseSection godoc
//
//	@Summary	Get course section
//	@Tags		course-sections
//	@Produce	json
//	@Param		id	path		int	true	"Section ID"
//	@Success	200	{object}	domain.CourseSection
//	@Failure	400	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/course-sections/{id} [get]
func (h *Handler) getCourseSection(ctx fiber.Ctx) error {
	id, err := utils.GetInt64Param(ctx, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	section, err := h.courseSectionsService.Get(ctx.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(section)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// getCourseSectionsList godoc
//
//	@Summary	List course sections
//	@Tags		course-sections
//	@Produce	json
//	@Success	200	{array}		domain.CourseSection
//	@Failure	500	{object}	ErrorResponse
//	@Router		/course-sections [get]
func (h *Handler) getCourseSectionsList(ctx fiber.Ctx) error {
	sections, err := h.courseSectionsService.GetList(ctx.Context(), coursesections.GetListServiceParams{})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(sections)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}
