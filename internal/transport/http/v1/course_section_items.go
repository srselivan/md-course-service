package v1

import (
	"course-service/internal/services/coursesectionitems"
	"course-service/internal/transport/http/utils"
	"course-service/internal/transport/http/v1/request"
	"encoding/json"
	"net/http"

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
