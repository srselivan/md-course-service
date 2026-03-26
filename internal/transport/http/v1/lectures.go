package v1

import (
	"course-service/internal/services/lectures"
	"course-service/internal/transport/http/utils"
	"course-service/internal/transport/http/v1/request"
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) NewLecturesRoutes(router fiber.Router) {
	lecturesGroup := router.Group("lectures")

	lecturesGroup.Post("/", h.createLecture)
	lecturesGroup.Put("/:item_id", h.updateLecture)
	lecturesGroup.Delete("/:item_id", h.deleteLecture)
	lecturesGroup.Get("/:item_id", h.getLecture)
	lecturesGroup.Get("/", h.getLecturesList)
}

func (h *Handler) createLecture(ctx fiber.Ctx) error {
	var req request.CreateLecture
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	lecture, err := h.lecturesService.Create(ctx.Context(), lectures.CreateServiceParams{
		ItemId:            req.ItemId,
		Content:           req.Content,
		VideoURL:          req.VideoURL,
		ReadingTimeMinute: req.ReadingTimeMinute,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(lecture)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusCreated).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *Handler) updateLecture(ctx fiber.Ctx) error {
	itemId, err := utils.GetInt64Param(ctx, "item_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req request.UpdateLecture
	if err = json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	lecture, err := h.lecturesService.Update(ctx.Context(), lectures.UpdateServiceParams{
		ItemId:            itemId,
		Content:           req.Content,
		VideoURL:          req.VideoURL,
		ReadingTimeMinute: req.ReadingTimeMinute,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(lecture)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *Handler) deleteLecture(ctx fiber.Ctx) error {
	itemId, err := utils.GetInt64Param(ctx, "item_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = h.lecturesService.Delete(ctx.Context(), itemId); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.SendStatus(http.StatusOK); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *Handler) getLecture(ctx fiber.Ctx) error {
	itemId, err := utils.GetInt64Param(ctx, "item_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	lecture, err := h.lecturesService.Get(ctx.Context(), itemId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(lecture)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *Handler) getLecturesList(ctx fiber.Ctx) error {
	lecturesList, err := h.lecturesService.GetList(ctx.Context(), lectures.GetListServiceParams{})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(lecturesList)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}
