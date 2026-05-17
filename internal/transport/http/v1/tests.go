package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"course-service/internal/domain"
	"course-service/internal/services/tests"
	"course-service/internal/transport/http/utils"
	"course-service/internal/transport/http/v1/request"

	"github.com/gofiber/fiber/v3"
)

func (h *Handler) NewTestsRoutes(router fiber.Router) {
	router.Get("/tests", h.getTestsList)
	router.Post("/tests", h.createTest)
	router.Put("/tests", h.updateTest)
	router.Put("/tests/:test_id", h.updateTest)
	router.Delete("/tests", h.deleteTest)
	router.Delete("/tests/:test_id", h.deleteTest)
	router.Get("/tests/:test_id", h.getTestInfo)
	router.Get("/tests/:test_id/attempts", h.getTestAttempts)
	router.Get("/tests/:test_id/attempts/:attempt_id", h.getTestAttemptState)
	router.Post("/tests/:test_id/attempts/:attempt_id/grades", h.gradeAttempt)
	router.Get("/attempts", h.getAttempts)
	router.Post("/attempts", h.startAttempt)
	router.Get("/attempts/:attempt_id", h.getAttemptState)
	router.Put("/attempts/:attempt_id/answers", h.saveAttemptAnswer)
	router.Post("/attempts/:attempt_id/submit", h.submitAttempt)
}

// getTestsList godoc
//
//	@Summary	List tests
//	@Tags		tests
//	@Produce	json
//	@Param		course_id	query		int	false	"Filter by course ID"
//	@Success	200			{array}		domain.Test
//	@Failure	400			{object}	ErrorResponse
//	@Failure	500			{object}	ErrorResponse
//	@Router		/tests [get]
func (h *Handler) getTestsList(ctx fiber.Ctx) error {
	courseID, err := utils.GetInt64pQuery(ctx, "course_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	list, err := h.testsService.GetList(ctx.Context(), courseID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return sendJSON(ctx, http.StatusOK, list)
}

// createTest godoc
//
//	@Summary	Create test
//	@Tags		tests
//	@Accept		json
//	@Produce	json
//	@Param		body	body		request.CreateTest	true	"Test"
//	@Success	201		{object}	domain.Test
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/tests [post]
func (h *Handler) createTest(ctx fiber.Ctx) error {
	var req request.CreateTest
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	test, err := h.testsService.Create(ctx.Context(), req.ToService())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return sendJSON(ctx, http.StatusCreated, test)
}

// updateTest godoc
//
//	@Summary	Update test
//	@Tags		tests
//	@Accept		json
//	@Produce	json
//	@Param		test_id	path		int					false	"Test ID (path)"
//	@Param		body	body		request.UpdateTest	true	"Test"
//	@Success	200		{object}	domain.Test
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/tests [put]
//	@Router		/tests/{test_id} [put]
func (h *Handler) updateTest(ctx fiber.Ctx) error {
	var req request.UpdateTest
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	testID, err := getTestID(ctx, req.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	test, err := h.testsService.Update(ctx.Context(), req.ToService(testID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return sendJSON(ctx, http.StatusOK, test)
}

// deleteTest godoc
//
//	@Summary	Delete test
//	@Tags		tests
//	@Param		test_id	query	int	false	"Test ID (query)"
//	@Param		test_id	path	int	false	"Test ID (path)"
//	@Success	204
//	@Failure	400	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/tests [delete]
//	@Router		/tests/{test_id} [delete]
func (h *Handler) deleteTest(ctx fiber.Ctx) error {
	testID, err := getTestID(ctx, 0)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = h.testsService.Delete(ctx.Context(), testID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// getTestInfo godoc
//
//	@Summary	Get test info with user attempts
//	@Tags		tests
//	@Produce	json
//	@Param		test_id	path		int	true	"Test ID"
//	@Param		user_id	query		int	true	"User ID"
//	@Success	200		{object}	domain.TestWithAttempts
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/tests/{test_id} [get]
func (h *Handler) getTestInfo(ctx fiber.Ctx) error {
	testID, err := utils.GetInt64Param(ctx, "test_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	userID, err := utils.GetInt64Query(ctx, "user_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	info, err := h.testsService.GetInfo(ctx.Context(), testID, userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return sendJSON(ctx, http.StatusOK, info)
}

// startAttempt godoc
//
//	@Summary	Start test attempt
//	@Tags		attempts
//	@Accept		json
//	@Produce	json
//	@Param		limit	query		int						false	"Questions page size"
//	@Param		offset	query		int						false	"Questions page offset"
//	@Param		body	body		request.StartAttempt	true	"Attempt"
//	@Success	201		{object}	domain.AttemptState
//	@Failure	400		{object}	ErrorResponse
//	@Router		/attempts [post]
func (h *Handler) startAttempt(ctx fiber.Ctx) error {
	var req request.StartAttempt
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	limit, err := getOptionalIntQuery(ctx, "limit")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	offset, err := getOptionalIntQuery(ctx, "offset")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	state, err := h.testsService.StartAttempt(ctx.Context(), tests.StartAttemptServiceParams{
		TestID: req.TestID,
		UserID: req.UserID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return sendJSON(ctx, http.StatusCreated, state)
}

// getAttemptState godoc
//
//	@Summary	Get attempt state
//	@Tags		attempts
//	@Produce	json
//	@Param		attempt_id	path		int	true	"Attempt ID"
//	@Param		limit		query		int	false	"Questions page size"
//	@Param		offset		query		int	false	"Questions page offset"
//	@Success	200			{object}	domain.AttemptState
//	@Failure	400			{object}	ErrorResponse
//	@Failure	500			{object}	ErrorResponse
//	@Router		/attempts/{attempt_id} [get]
func (h *Handler) getAttemptState(ctx fiber.Ctx) error {
	attemptID, err := utils.GetInt64Param(ctx, "attempt_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	limit, err := getOptionalIntQuery(ctx, "limit")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	offset, err := getOptionalIntQuery(ctx, "offset")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	state, err := h.testsService.GetAttemptState(ctx.Context(), tests.GetAttemptStateServiceParams{
		AttemptID: attemptID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return sendJSON(ctx, http.StatusOK, state)
}

// saveAttemptAnswer godoc
//
//	@Summary	Save attempt answer
//	@Tags		attempts
//	@Accept		json
//	@Param		attempt_id	path	int							true	"Attempt ID"
//	@Param		body		body	request.SaveAttemptAnswer	true	"Answer"
//	@Success	204
//	@Failure	400	{object}	ErrorResponse
//	@Router		/attempts/{attempt_id}/answers [put]
func (h *Handler) saveAttemptAnswer(ctx fiber.Ctx) error {
	attemptID, err := utils.GetInt64Param(ctx, "attempt_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req request.SaveAttemptAnswer
	if err = json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = h.testsService.SaveAnswer(ctx.Context(), req.ToService(attemptID)); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// submitAttempt godoc
//
//	@Summary	Submit attempt for grading
//	@Tags		attempts
//	@Produce	json
//	@Param		attempt_id	path		int	true	"Attempt ID"
//	@Success	200			{object}	domain.TestAttempt
//	@Failure	400			{object}	ErrorResponse
//	@Router		/attempts/{attempt_id}/submit [post]
func (h *Handler) submitAttempt(ctx fiber.Ctx) error {
	attemptID, err := utils.GetInt64Param(ctx, "attempt_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	attempt, err := h.testsService.SubmitAttempt(ctx.Context(), attemptID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return sendJSON(ctx, http.StatusOK, attempt)
}

// getAttempts godoc
//
//	@Summary	List attempts
//	@Tags		attempts
//	@Produce	json
//	@Param		test_id	query		int		false	"Filter by test ID"
//	@Param		user_id	query		int		false	"Filter by user ID"
//	@Param		status	query		string	false	"Filter by status"
//	@Success	200		{array}		domain.TestAttempt
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/attempts [get]
func (h *Handler) getAttempts(ctx fiber.Ctx) error {
	testID, err := utils.GetInt64pQuery(ctx, "test_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	userID, err := utils.GetInt64pQuery(ctx, "user_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var status *domain.TestAttemptStatus
	if statusQuery := ctx.Query("status"); statusQuery != "" {
		s := domain.TestAttemptStatus(statusQuery)
		status = &s
	}

	attempts, err := h.testsService.GetAttempts(ctx.Context(), tests.GetAttemptsServiceParams{
		TestID: testID,
		UserID: userID,
		Status: status,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return sendJSON(ctx, http.StatusOK, attempts)
}

// getTestAttempts godoc
//
//	@Summary	List attempts for a test
//	@Tags		tests
//	@Produce	json
//	@Param		test_id	path		int		true	"Test ID"
//	@Param		user_id	query		int		false	"Filter by user ID"
//	@Param		status	query		string	false	"Filter by status"
//	@Success	200		{array}		domain.TestAttempt
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/tests/{test_id}/attempts [get]
func (h *Handler) getTestAttempts(ctx fiber.Ctx) error {
	testID, err := utils.GetInt64Param(ctx, "test_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	userID, err := utils.GetInt64pQuery(ctx, "user_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var status *domain.TestAttemptStatus
	if statusQuery := ctx.Query("status"); statusQuery != "" {
		s := domain.TestAttemptStatus(statusQuery)
		status = &s
	}

	attempts, err := h.testsService.GetAttempts(ctx.Context(), tests.GetAttemptsServiceParams{
		TestID: &testID,
		UserID: userID,
		Status: status,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return sendJSON(ctx, http.StatusOK, attempts)
}

// getTestAttemptState godoc
//
//	@Summary	Get attempt state within a test
//	@Tags		tests
//	@Produce	json
//	@Param		test_id		path		int	true	"Test ID"
//	@Param		attempt_id	path		int	true	"Attempt ID"
//	@Param		limit		query		int	false	"Questions page size"
//	@Param		offset		query		int	false	"Questions page offset"
//	@Success	200			{object}	domain.AttemptState
//	@Failure	400			{object}	ErrorResponse
//	@Failure	404			{object}	ErrorResponse
//	@Failure	500			{object}	ErrorResponse
//	@Router		/tests/{test_id}/attempts/{attempt_id} [get]
func (h *Handler) getTestAttemptState(ctx fiber.Ctx) error {
	testID, err := utils.GetInt64Param(ctx, "test_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	attemptID, err := utils.GetInt64Param(ctx, "attempt_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	limit, err := getOptionalIntQuery(ctx, "limit")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	offset, err := getOptionalIntQuery(ctx, "offset")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	state, err := h.testsService.GetAttemptState(ctx.Context(), tests.GetAttemptStateServiceParams{
		AttemptID: attemptID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if state.Attempt.TestID != testID {
		return fiber.NewError(fiber.StatusNotFound, "attempt not found for test")
	}

	return sendJSON(ctx, http.StatusOK, state)
}

// gradeAttempt godoc
//
//	@Summary	Grade attempt manually
//	@Tags		tests
//	@Accept		json
//	@Produce	json
//	@Param		test_id		path		int						true	"Test ID"
//	@Param		attempt_id	path		int						true	"Attempt ID"
//	@Param		body		body		request.GradeAttempt	true	"Scores"
//	@Success	200			{object}	domain.TestAttempt
//	@Failure	400			{object}	ErrorResponse
//	@Failure	404			{object}	ErrorResponse
//	@Failure	500			{object}	ErrorResponse
//	@Router		/tests/{test_id}/attempts/{attempt_id}/grades [post]
func (h *Handler) gradeAttempt(ctx fiber.Ctx) error {
	testID, err := utils.GetInt64Param(ctx, "test_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	attemptID, err := utils.GetInt64Param(ctx, "attempt_id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var req request.GradeAttempt
	if err = json.Unmarshal(ctx.Body(), &req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	state, err := h.testsService.GetAttemptState(ctx.Context(), tests.GetAttemptStateServiceParams{
		AttemptID: attemptID,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if state.Attempt.TestID != testID {
		return fiber.NewError(fiber.StatusNotFound, "attempt not found for test")
	}

	attempt, err := h.testsService.GradeAttempt(ctx.Context(), req.ToService(attemptID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return sendJSON(ctx, http.StatusOK, attempt)
}

func sendJSON(ctx fiber.Ctx, status int, value any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if err = ctx.Status(status).Send(bytes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

func getOptionalIntQuery(ctx fiber.Ctx, key string) (int, error) {
	raw := ctx.Query(key)
	if raw == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func getTestID(ctx fiber.Ctx, fallback int64) (int64, error) {
	if raw := ctx.Params("test_id"); raw != "" {
		return strconv.ParseInt(raw, 10, 64)
	}
	if raw := ctx.Query("test_id"); raw != "" {
		return strconv.ParseInt(raw, 10, 64)
	}
	if fallback > 0 {
		return fallback, nil
	}
	return 0, fiber.NewError(fiber.StatusBadRequest, "test_id is required")
}
