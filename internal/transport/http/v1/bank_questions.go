package v1

import (
	"encoding/json"
	"net/http"

	"github.com/samber/lo"

	"course-service/internal/services/bankquestions"
	"course-service/internal/transport/http/utils"
	"course-service/internal/transport/http/v1/request"

	_ "course-service/internal/domain"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) NewBankQuestionsRoutes(router fiber.Router) {
	qGroup := router.Group("bank_questions")

	qGroup.Post("/bulk", h.bulkCreateBankQuestion)
	qGroup.Put("/:id/bulk", h.bulkUpdateBankQuestion)
	qGroup.Get("/", h.getBankQuestionsList)
	qGroup.Post("/", h.createBankQuestion)
	qGroup.Put("/:id", h.updateBankQuestion)
	qGroup.Delete("/:id", h.deleteBankQuestion)
}

// createBankQuestion godoc
//
//	@Summary	Create bank question
//	@Tags		bank-questions
//	@Accept		json
//	@Produce	json
//	@Param		body	body		request.CreateBankQuestion	true	"Question"
//	@Success	201		{object}	domain.BankQuestion
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/bank_questions [post]
func (h *Handler) createBankQuestion(ctx fiber.Ctx) error {
	var req request.CreateBankQuestion
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	question, err := h.bankQuestionsService.Create(ctx.Context(), req.ToService())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(question)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusCreated).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// bulkCreateBankQuestion godoc
//
//	@Summary	Bulk create bank questions
//	@Tags		bank-questions
//	@Accept		json
//	@Param		body	body	request.BulkCreateBankQuestion	true	"Questions"
//	@Success	201
//	@Failure	400	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/bank_questions/bulk [post]
func (h *Handler) bulkCreateBankQuestion(ctx fiber.Ctx) error {
	var req request.BulkCreateBankQuestion
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err := h.bankQuestionsService.BulkCreate(
		ctx.Context(),
		bankquestions.BulkCreateServiceParams{
			Questions: lo.Map(req.Questions, func(item request.CreateBankQuestion, _ int) bankquestions.CreateServiceParams {
				return item.ToService()
			}),
		},
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusCreated).Send(nil); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// updateBankQuestion godoc
//
//	@Summary	Update bank question
//	@Tags		bank-questions
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int							true	"Question ID"
//	@Param		body	body		request.UpdateBankQuestion	true	"Question"
//	@Success	200		{object}	domain.BankQuestion
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/bank_questions/{id} [put]
func (h *Handler) updateBankQuestion(ctx fiber.Ctx) error {
	id, err := utils.GetInt64Param(ctx, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req request.UpdateBankQuestion
	if err = json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	question, err := h.bankQuestionsService.Update(ctx.Context(), req.ToService(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(question)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// bulkUpdateBankQuestion godoc
//
//	@Summary	Replace all questions in a bank
//	@Tags		bank-questions
//	@Accept		json
//	@Param		id		path	int								true	"Bank ID"
//	@Param		body	body	request.BulkUpdateBankQuestion	true	"Questions"
//	@Success	200
//	@Failure	400	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/bank_questions/{id}/bulk [put]
func (h *Handler) bulkUpdateBankQuestion(ctx fiber.Ctx) error {
	id, err := utils.GetInt64Param(ctx, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req request.BulkUpdateBankQuestion
	if err = json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = h.bankQuestionsService.BulkUpdate(
		ctx.Context(),
		bankquestions.BulkUpdateServiceParams{
			Id: id,
			Questions: lo.Map(req.Questions, func(item request.UpdateBankQuestion, _ int) bankquestions.UpdateServiceParams {
				return item.ToService(0)
			}),
		},
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(nil); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// deleteBankQuestion godoc
//
//	@Summary	Delete bank question
//	@Tags		bank-questions
//	@Param		id	path	int	true	"Question ID"
//	@Success	200
//	@Failure	400	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/bank_questions/{id} [delete]
func (h *Handler) deleteBankQuestion(ctx fiber.Ctx) error {
	id, err := utils.GetInt64Param(ctx, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = h.bankQuestionsService.Delete(ctx.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.SendStatus(http.StatusOK); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

// getBankQuestionsList godoc
//
//	@Summary	List bank questions
//	@Tags		bank-questions
//	@Produce	json
//	@Param		bank_id			query		int		true	"Bank ID"
//	@Param		limit			query		int		true	"Page size"
//	@Param		offset			query		int		true	"Page offset"
//	@Param		question_type	query		string	false	"Filter by question type"	Enums(SINGLE,MULTIPLE,TEXT)
//	@Param		filter			query		string	false	"Filter by question text (substring, case-insensitive)"
//	@Success	200				{object}	domain.BankQuestionsListResponse
//	@Failure	400				{object}	ErrorResponse
//	@Failure	500				{object}	ErrorResponse
//	@Router		/bank_questions [get]
func (h *Handler) getBankQuestionsList(ctx fiber.Ctx) error {
	bankId, err := utils.GetInt64Query(ctx, "bank_id")
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

	params := bankquestions.GetListServiceParams{
		BankId: bankId,
		Limit:  limit,
		Offset: offset,
	}

	if questionType := ctx.Query("question_type"); questionType != "" {
		params.Type = &questionType
	}

	filter := ctx.Query("filter")
	if filter != "" {
		params.Filter = &filter
	}

	response, err := h.bankQuestionsService.GetList(ctx.Context(), params)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}
