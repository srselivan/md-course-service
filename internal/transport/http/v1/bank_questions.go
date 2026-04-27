package v1

import (
	"encoding/json"
	"net/http"

	"github.com/samber/lo"

	"course-service/internal/domain"
	"course-service/internal/services/bankquestions"
	"course-service/internal/transport/http/utils"
	"course-service/internal/transport/http/v1/request"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) NewBankQuestionsRoutes(router fiber.Router) {
	qGroup := router.Group("bank_questions")

	// qGroup.Post("/", h.createBankQuestion)
	qGroup.Post("/bulk", h.bulkCreateBankQuestion)
	qGroup.Put("/:id/bulk", h.bulkUpdateBankQuestion)
	// qGroup.Put("/:id", h.updateBankQuestion)
	qGroup.Get("/", h.getBankQuestionsList)
	qGroup.Post("/generate", h.generateBankQuestions)
}

func (h *Handler) createBankQuestion(ctx fiber.Ctx) error {
	var req request.CreateBankQuestion
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	answers := make([]bankquestions.BankAnswerDTO, 0, len(req.Answers))
	for _, a := range req.Answers {
		answers = append(answers, bankquestions.BankAnswerDTO{
			AnswerText: a.AnswerText,
			IsCorrect:  a.IsCorrect,
		})
	}

	question, err := h.bankQuestionsService.Create(ctx.Context(), bankquestions.CreateServiceParams{
		QuestionText:  req.QuestionText,
		QuestionType:  domain.QuestionType(req.QuestionType),
		DefaultPoints: req.DefaultPoints,
		BankId:        req.BankId,
		Answers:       answers,
	})
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

func (h *Handler) updateBankQuestion(ctx fiber.Ctx) error {
	id, err := utils.GetInt64Param(ctx, "id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req request.UpdateBankQuestion
	if err = json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	answers := make([]bankquestions.BankAnswerDTO, 0, len(req.Answers))
	for _, a := range req.Answers {
		answers = append(answers, bankquestions.BankAnswerDTO{
			AnswerText: a.AnswerText,
			IsCorrect:  a.IsCorrect,
		})
	}

	question, err := h.bankQuestionsService.Update(ctx.Context(), bankquestions.UpdateServiceParams{
		ID:            id,
		QuestionText:  req.QuestionText,
		QuestionType:  domain.QuestionType(req.QuestionType),
		DefaultPoints: req.DefaultPoints,
		BankId:        req.BankId,
		Answers:       answers,
	})
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
				return item.ToService()
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

func (h *Handler) getBankQuestionsList(ctx fiber.Ctx) error {
	bankId, err := utils.GetInt64Query(ctx, "bank_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	list, err := h.bankQuestionsService.GetList(ctx.Context(), bankquestions.GetListServiceParams{
		BankId: &bankId,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	bytes, err := json.Marshal(list)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err = ctx.Status(http.StatusOK).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (h *Handler) generateBankQuestions(ctx fiber.Ctx) error {

	return nil
}
