package v1

import (
	"course-service/internal/services/assignments"
	"course-service/internal/transport/http/utils"
	"course-service/internal/transport/http/v1/request"
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) NewAssignmentsRoutes(router fiber.Router) {
	assignmentsGroup := router.Group("assignments")

	assignmentsGroup.Post("/", h.createAssignment)
	assignmentsGroup.Put("/:item_id", h.updateAssignment)
	assignmentsGroup.Delete("/:item_id", h.deleteAssignment)
	assignmentsGroup.Get("/:item_id", h.getAssignment)
	assignmentsGroup.Get("/", h.getAssignmentsList)
}

func (h *Handler) createAssignment(ctx fiber.Ctx) error {
	var req request.CreateAssignment
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	assignment, err := h.assignmentsService.Create(ctx.Context(), assignments.CreateServiceParams{
		ItemId:       req.ItemId,
		Description:  req.Description,
		MaxScore:     req.MaxScore,
		DeadlineDays: req.DeadlineDays,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(assignment)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusCreated).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *Handler) updateAssignment(ctx fiber.Ctx) error {
	itemId, err := utils.GetInt64Param(ctx, "item_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req request.UpdateAssignment
	if err = json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	assignment, err := h.assignmentsService.Update(ctx.Context(), assignments.UpdateServiceParams{
		ItemId:       itemId,
		Description:  req.Description,
		MaxScore:     req.MaxScore,
		DeadlineDays: req.DeadlineDays,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(assignment)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *Handler) deleteAssignment(ctx fiber.Ctx) error {
	itemId, err := utils.GetInt64Param(ctx, "item_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = h.assignmentsService.Delete(ctx.Context(), itemId); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.SendStatus(http.StatusOK); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *Handler) getAssignment(ctx fiber.Ctx) error {
	itemId, err := utils.GetInt64Param(ctx, "item_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	assignment, err := h.assignmentsService.Get(ctx.Context(), itemId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(assignment)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *Handler) getAssignmentsList(ctx fiber.Ctx) error {
	assignmentsList, err := h.assignmentsService.GetList(ctx.Context(), assignments.GetListServiceParams{})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(assignmentsList)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}
