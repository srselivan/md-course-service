package v1

import (
	"encoding/json"
	"net/http"

	"course-service/internal/services/banks"
	"course-service/internal/transport/http/utils"
	"course-service/internal/transport/http/v1/request"

	_ "course-service/internal/domain"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) NewBanksRoutes(router fiber.Router) {
	banksGroup := router.Group("banks")

	banksGroup.Post("/", h.createBank)
	banksGroup.Put("/:id", h.updateBank)
	banksGroup.Delete("/:id", h.deleteBank)
	banksGroup.Get("/", h.getBanksList)
}

// createBank godoc
//
//	@Summary	Create question bank
//	@Tags		banks
//	@Accept		json
//	@Produce	json
//	@Param		body	body		request.CreateBank	true	"Bank"
//	@Success	201		{object}	domain.Bank
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/banks [post]
func (h *Handler) createBank(ctx fiber.Ctx) error {
	var req request.CreateBank
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	bank, err := h.banksService.Create(ctx.Context(), banks.CreateServiceParams{
		Title:       req.Title,
		Description: req.Description,
		UserId:      req.UserId,
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

// updateBank godoc
//
//	@Summary		Update question bank
//	@Description	Updates title and description. userId in body must match the bank owner; userId is not changed.
//	@Tags			banks
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int					true	"Bank ID"
//	@Param		body	body		request.UpdateBank	true	"Bank"
//	@Success	200		{object}	domain.Bank
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/banks/{id} [put]
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
		UserId:      req.UserId,
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

// deleteBank godoc
//
//	@Summary	Delete question bank
//	@Tags		banks
//	@Param		id	path	int	true	"Bank ID"
//	@Success	200
//	@Failure	400	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/banks/{id} [delete]
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

// getBanksList godoc
//
//	@Summary	List question banks
//	@Tags		banks
//	@Produce	json
//	@Param		user_id	query		int	true	"Owner user ID"
//	@Success	200		{array}		domain.Bank
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/banks [get]
func (h *Handler) getBanksList(ctx fiber.Ctx) error {
	userId, err := utils.GetInt64Query(ctx, "user_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	banksList, err := h.banksService.GetList(ctx.Context(), banks.GetListServiceParams{
		UserId: userId,
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
