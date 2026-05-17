package v1

import (
	"encoding/json"
	"net/http"

	"course-service/internal/services/assignments"
	"course-service/internal/transport/http/utils"
	"course-service/internal/transport/http/v1/request"

	_ "course-service/internal/domain"

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

// createAssignment godoc
//
//	@Summary	Create assignment
//	@Tags		assignments
//	@Accept		json
//	@Produce	json
//	@Param		body	body		request.CreateAssignment	true	"Assignment"
//	@Success	201		{object}	domain.Assignment
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/assignments [post]
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

// updateAssignment godoc
//
//	@Summary	Update assignment
//	@Tags		assignments
//	@Accept		json
//	@Produce	json
//	@Param		item_id	path		int							true	"Section item ID"
//	@Param		body	body		request.UpdateAssignment	true	"Assignment"
//	@Success	200		{object}	domain.Assignment
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/assignments/{item_id} [put]
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

// deleteAssignment godoc
//
//	@Summary	Delete assignment
//	@Tags		assignments
//	@Param		item_id	path	int	true	"Section item ID"
//	@Success	200
//	@Failure	400	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/assignments/{item_id} [delete]
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

// getAssignment godoc
//
//	@Summary	Get assignment
//	@Tags		assignments
//	@Produce	json
//	@Param		item_id	path		int	true	"Section item ID"
//	@Success	200		{object}	domain.Assignment
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/assignments/{item_id} [get]
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

// getAssignmentsList godoc
//
//	@Summary	List assignments
//	@Tags		assignments
//	@Produce	json
//	@Success	200	{array}		domain.Assignment
//	@Failure	500	{object}	ErrorResponse
//	@Router		/assignments [get]
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
