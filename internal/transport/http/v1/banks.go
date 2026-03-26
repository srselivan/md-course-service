package v1

import (
	"course-service/internal/services/banks"
	"course-service/internal/transport/http/utils"
	"course-service/internal/transport/http/v1/request"
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) NewBanksRoutes(router fiber.Router) {
	banksGroup := router.Group("banks")

	banksGroup.Post("/", h.createBank)
	banksGroup.Put("/:id", h.updateBank)
	banksGroup.Delete("/:id", h.deleteBank)
	banksGroup.Get("/", h.getBanksList)
}

func (h *Handler) createBank(ctx fiber.Ctx) error {
	var req request.CreateBank
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	bank, err := h.banksService.Create(ctx.Context(), banks.CreateServiceParams{
		Title:       req.Title,
		Description: req.Description,
		CourseId:    req.CourseId,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(bank)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusCreated).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *Handler) updateBank(ctx fiber.Ctx) error {
	id, err := utils.GetInt64Param(ctx, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req request.UpdateBank
	if err = json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	bank, err := h.banksService.Update(ctx.Context(), banks.UpdateServiceParams{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		CourseId:    req.CourseId,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(bank)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *Handler) deleteBank(ctx fiber.Ctx) error {
	id, err := utils.GetInt64Param(ctx, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = h.banksService.Delete(ctx.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.SendStatus(http.StatusOK); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *Handler) getBanksList(ctx fiber.Ctx) error {
	courseId, err := utils.GetInt64Query(ctx, "course_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	banksList, err := h.banksService.GetList(ctx.Context(), banks.GetListServiceParams{
		CourseId: &courseId,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(banksList)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}
