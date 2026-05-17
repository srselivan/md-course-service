package v1

import (
	"encoding/json"
	"net/http"

	"course-service/internal/services/coursesectionitems"
	"course-service/internal/transport/http/utils"
	"course-service/internal/transport/http/v1/request"

	_ "course-service/internal/domain"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) NewCourseSectionItemsRoutes(router fiber.Router) {
	itemsGroup := router.Group("course-section-items")

	itemsGroup.Post("/", h.createCourseSectionItem)
	itemsGroup.Put("/:id", h.updateCourseSectionItem)
	itemsGroup.Delete("/:id", h.deleteCourseSectionItem)
	itemsGroup.Get("/:id", h.getCourseSectionItem)
	itemsGroup.Get("/", h.getCourseSectionItemsList)
}

// createCourseSectionItem godoc
//
//	@Summary	Create section item
//	@Tags		course-section-items
//	@Accept		json
//	@Produce	json
//	@Param		body	body		request.CreateCourseSectionItem	true	"Item"
//	@Success	201		{object}	domain.CourseSectionItem
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/course-section-items [post]
func (h *Handler) createCourseSectionItem(ctx fiber.Ctx) error {
	var req request.CreateCourseSectionItem
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	item, err := h.courseSectionItemsService.Create(ctx.Context(), coursesectionitems.CreateServiceParams{
		SectionId:   req.SectionId,
		ItemType:    req.ItemType,
		Title:       req.Title,
		SortOrder:   req.SortOrder,
		IsPublished: req.IsPublished,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(item)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusCreated).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// updateCourseSectionItem godoc
//
//	@Summary	Update section item
//	@Tags		course-section-items
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int								true	"Item ID"
//	@Param		body	body		request.UpdateCourseSectionItem	true	"Item"
//	@Success	200		{object}	domain.CourseSectionItem
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/course-section-items/{id} [put]
func (h *Handler) updateCourseSectionItem(ctx fiber.Ctx) error {
	id, err := utils.GetInt64Param(ctx, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req request.UpdateCourseSectionItem
	if err = json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	item, err := h.courseSectionItemsService.Update(ctx.Context(), coursesectionitems.UpdateServiceParams{
		ID:          id,
		SectionId:   req.SectionId,
		ItemType:    req.ItemType,
		Title:       req.Title,
		SortOrder:   req.SortOrder,
		IsPublished: req.IsPublished,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(item)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// deleteCourseSectionItem godoc
//
//	@Summary	Delete section item
//	@Tags		course-section-items
//	@Param		id	path	int	true	"Item ID"
//	@Success	200
//	@Failure	400	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/course-section-items/{id} [delete]
func (h *Handler) deleteCourseSectionItem(ctx fiber.Ctx) error {
	id, err := utils.GetInt64Param(ctx, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = h.courseSectionItemsService.Delete(ctx.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.SendStatus(http.StatusOK); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// getCourseSectionItem godoc
//
//	@Summary	Get section item
//	@Tags		course-section-items
//	@Produce	json
//	@Param		id	path		int	true	"Item ID"
//	@Success	200	{object}	domain.CourseSectionItem
//	@Failure	400	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/course-section-items/{id} [get]
func (h *Handler) getCourseSectionItem(ctx fiber.Ctx) error {
	id, err := utils.GetInt64Param(ctx, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	item, err := h.courseSectionItemsService.Get(ctx.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(item)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// getCourseSectionItemsList godoc
//
//	@Summary	List section items
//	@Tags		course-section-items
//	@Produce	json
//	@Success	200	{array}		domain.CourseSectionItem
//	@Failure	500	{object}	ErrorResponse
//	@Router		/course-section-items [get]
func (h *Handler) getCourseSectionItemsList(ctx fiber.Ctx) error {
	items, err := h.courseSectionItemsService.GetList(ctx.Context(), coursesectionitems.GetListServiceParams{})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(items)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}
