package tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"course-service/internal/domain"

	"github.com/rs/zerolog"
)

const defaultAttemptPageLimit = 50

type testsRepo interface {
	Create(ctx context.Context, params CreateRepoParams) (domain.Test, error)
	Update(ctx context.Context, params UpdateRepoParams) (domain.Test, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (domain.Test, error)
	GetList(ctx context.Context, courseID *int64) ([]domain.Test, error)

	GetAttempts(ctx context.Context, params GetAttemptsRepoParams) ([]domain.TestAttempt, error)
	CountAttempts(ctx context.Context, testID, userID int64) (int, error)
	SelectBankQuestionIDs(ctx context.Context, params SelectBankQuestionIDsRepoParams) ([]int64, error)
	GetQuestionsWithAnswers(ctx context.Context, questionIDs []int64) ([]domain.Question, error)

	CreateAttempt(ctx context.Context, params CreateAttemptRepoParams) (domain.TestAttempt, error)
	GetAttempt(ctx context.Context, attemptID int64) (domain.TestAttempt, error)
	GetAttemptQuestion(ctx context.Context, attemptQuestionID int64) (domain.AttemptQuestion, error)
	GetAttemptState(ctx context.Context, params GetAttemptStateRepoParams) (domain.AttemptState, error)
	GetAttemptQuestionsForScoring(ctx context.Context, attemptID int64) ([]domain.AttemptQuestionState, error)
	GetAttemptScores(ctx context.Context, attemptID int64) ([]AttemptQuestionScore, error)

	SaveTextAnswer(ctx context.Context, params SaveTextAnswerRepoParams) error
	SaveChoiceAnswers(ctx context.Context, params SaveChoiceAnswersRepoParams) error

	CompleteAttempt(ctx context.Context, params CompleteAttemptRepoParams) (domain.TestAttempt, error)
	GradeAttempt(ctx context.Context, params GradeAttemptRepoParams) (domain.TestAttempt, error)

	GetExpiredAttemptIDs(ctx context.Context, now time.Time) ([]int64, error)
}

type Service struct {
	repo   testsRepo
	logger *zerolog.Logger
}

func NewService(repo testsRepo, logger *zerolog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, params CreateServiceParams) (domain.Test, error) {
	test, err := s.repo.Create(ctx, CreateRepoParams(params))
	if err != nil {
		return domain.Test{}, fmt.Errorf("repo.Create: %w", err)
	}
	return test, nil
}

func (s *Service) Update(ctx context.Context, params UpdateServiceParams) (domain.Test, error) {
	test, err := s.repo.Update(ctx, UpdateRepoParams(params))
	if err != nil {
		return domain.Test{}, fmt.Errorf("repo.Update: %w", err)
	}
	return test, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("repo.Delete: %w", err)
	}
	return nil
}

func (s *Service) GetList(ctx context.Context, courseID *int64) ([]domain.Test, error) {
	tests, err := s.repo.GetList(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("repo.GetList: %w", err)
	}
	return tests, nil
}

func (s *Service) GetInfo(ctx context.Context, testID, userID int64) (domain.TestWithAttempts, error) {
	test, err := s.repo.Get(ctx, testID)
	if err != nil {
		return domain.TestWithAttempts{}, fmt.Errorf("repo.Get: %w", err)
	}
	attempts, err := s.repo.GetAttempts(ctx, GetAttemptsRepoParams{TestID: &testID, UserID: &userID})
	if err != nil {
		return domain.TestWithAttempts{}, fmt.Errorf("repo.GetAttempts: %w", err)
	}
	return domain.TestWithAttempts{Test: test, Attempts: attempts}, nil
}

func (s *Service) GetAttempts(ctx context.Context, params GetAttemptsServiceParams) ([]domain.TestAttempt, error) {
	attempts, err := s.repo.GetAttempts(ctx, GetAttemptsRepoParams(params))
	if err != nil {
		return nil, fmt.Errorf("repo.GetAttempts: %w", err)
	}
	return attempts, nil
}

func (s *Service) StartAttempt(ctx context.Context, params StartAttemptServiceParams) (domain.AttemptState, error) {
	test, err := s.repo.Get(ctx, params.TestID)
	if err != nil {
		return domain.AttemptState{}, fmt.Errorf("repo.Get: %w", err)
	}

	now := time.Now()
	if now.Before(test.AvailableFrom) {
		return domain.AttemptState{}, errors.New("test is not available yet")
	}
	if now.After(test.AvailableTo) {
		return domain.AttemptState{}, errors.New("test availability period is over")
	}

	attemptsCount, err := s.repo.CountAttempts(ctx, params.TestID, params.UserID)
	if err != nil {
		return domain.AttemptState{}, fmt.Errorf("repo.CountAttempts: %w", err)
	}
	if attemptsCount >= test.MaxAttempts {
		return domain.AttemptState{}, errors.New("max attempts limit exceeded")
	}

	settings, err := decodeGenerationSettings(test.GenerationSettings)
	if err != nil {
		return domain.AttemptState{}, err
	}
	questionsLimit := test.QuestionsCount

	questions, err := s.selectQuestionsForAttempt(ctx, params.TestID, questionsLimit, settings)
	if err != nil {
		return domain.AttemptState{}, err
	}
	if len(questions) < questionsLimit {
		return domain.AttemptState{}, errors.New("not enough questions in linked banks")
	}

	attempt, err := s.repo.CreateAttempt(ctx, CreateAttemptRepoParams{
		TestID:        params.TestID,
		UserID:        params.UserID,
		AttemptNumber: attemptsCount + 1,
		StartedAt:     now,
		ExpiresAt:     now.Add(time.Duration(test.DurationSeconds) * time.Second),
		Questions:     questions,
	})
	if err != nil {
		return domain.AttemptState{}, fmt.Errorf("repo.CreateAttempt: %w", err)
	}

	return s.GetAttemptState(ctx, GetAttemptStateServiceParams{
		AttemptID: attempt.ID,
		Limit:     params.Limit,
		Offset:    params.Offset,
	})
}

func (s *Service) GetAttemptState(ctx context.Context, params GetAttemptStateServiceParams) (domain.AttemptState, error) {
	state, err := s.repo.GetAttemptState(ctx, GetAttemptStateRepoParams{
		AttemptID: params.AttemptID,
		Limit:     normalizeLimit(params.Limit),
		Offset:    params.Offset,
	})
	if err != nil {
		return domain.AttemptState{}, fmt.Errorf("repo.GetAttemptState: %w", err)
	}

	remaining := time.Until(state.Attempt.ExpiresAt).Seconds()
	if remaining < 0 {
		remaining = 0
	}
	state.RemainingSeconds = int64(remaining)
	return state, nil
}

func (s *Service) SaveAnswer(ctx context.Context, params SaveAnswerServiceParams) error {
	attempt, err := s.repo.GetAttempt(ctx, params.AttemptID)
	if err != nil {
		return fmt.Errorf("repo.GetAttempt: %w", err)
	}
	if attempt.Status != domain.TestAttemptStatusInProgress {
		return errors.New("attempt is not in progress")
	}
	if time.Now().After(attempt.ExpiresAt) {
		if _, submitErr := s.SubmitAttempt(ctx, params.AttemptID); submitErr != nil {
			return fmt.Errorf("submit expired attempt: %w", submitErr)
		}
		return errors.New("attempt time is over")
	}

	question, err := s.repo.GetAttemptQuestion(ctx, params.AttemptQuestionID)
	if err != nil {
		return fmt.Errorf("repo.GetAttemptQuestion: %w", err)
	}
	if question.AttemptID != params.AttemptID {
		return errors.New("attempt question does not belong to attempt")
	}

	if question.Type == domain.TestingQuestionTypeText {
		if params.TextResponse == nil {
			return errors.New("text response is required for text question")
		}
		if err = s.repo.SaveTextAnswer(ctx, SaveTextAnswerRepoParams{
			AttemptQuestionID: params.AttemptQuestionID,
			TextResponse:      *params.TextResponse,
		}); err != nil {
			return fmt.Errorf("repo.SaveTextAnswer: %w", err)
		}
		return nil
	}

	if err = s.repo.SaveChoiceAnswers(ctx, SaveChoiceAnswersRepoParams{
		AttemptQuestionID: params.AttemptQuestionID,
		SelectedAnswerIDs: params.SelectedAnswerIDs,
	}); err != nil {
		return fmt.Errorf("repo.SaveChoiceAnswers: %w", err)
	}
	return nil
}

func (s *Service) SubmitWithAnswers(ctx context.Context, params SubmitWithAnswersServiceParams) (domain.TestAttempt, error) {
	attempt, err := s.repo.GetAttempt(ctx, params.AttemptID)
	if err != nil {
		return domain.TestAttempt{}, fmt.Errorf("repo.GetAttempt: %w", err)
	}
	if attempt.Status != domain.TestAttemptStatusInProgress {
		return domain.TestAttempt{}, errors.New("attempt is not in progress")
	}

	for _, answer := range params.Answers {
		answer.AttemptID = params.AttemptID
		question, qErr := s.repo.GetAttemptQuestion(ctx, answer.AttemptQuestionID)
		if qErr != nil {
			return domain.TestAttempt{}, fmt.Errorf("repo.GetAttemptQuestion(%d): %w", answer.AttemptQuestionID, qErr)
		}
		if question.AttemptID != params.AttemptID {
			return domain.TestAttempt{}, errors.New("attempt question does not belong to attempt")
		}

		if question.Type == domain.TestingQuestionTypeText {
			if answer.TextResponse == nil {
				continue
			}
			if err = s.repo.SaveTextAnswer(ctx, SaveTextAnswerRepoParams{
				AttemptQuestionID: answer.AttemptQuestionID,
				TextResponse:      *answer.TextResponse,
			}); err != nil {
				return domain.TestAttempt{}, fmt.Errorf("repo.SaveTextAnswer: %w", err)
			}
		} else {
			if err = s.repo.SaveChoiceAnswers(ctx, SaveChoiceAnswersRepoParams{
				AttemptQuestionID: answer.AttemptQuestionID,
				SelectedAnswerIDs: answer.SelectedAnswerIDs,
			}); err != nil {
				return domain.TestAttempt{}, fmt.Errorf("repo.SaveChoiceAnswers: %w", err)
			}
		}
	}

	return s.SubmitAttempt(ctx, params.AttemptID)
}

func (s *Service) SubmitAttempt(ctx context.Context, attemptID int64) (domain.TestAttempt, error) {
	attempt, err := s.repo.GetAttempt(ctx, attemptID)
	if err != nil {
		return domain.TestAttempt{}, fmt.Errorf("repo.GetAttempt: %w", err)
	}
	if attempt.Status != domain.TestAttemptStatusInProgress {
		return domain.TestAttempt{}, errors.New("attempt is not in progress")
	}

	questions, err := s.repo.GetAttemptQuestionsForScoring(ctx, attemptID)
	if err != nil {
		return domain.TestAttempt{}, fmt.Errorf("repo.GetAttemptQuestionsForScoring: %w", err)
	}

	scores, totalScore, needsGrading := evaluateAttempt(questions)

	gradeStatus := domain.TestGradeStatusReady
	if needsGrading {
		gradeStatus = domain.TestGradeStatusPreGraded
	}

	completed, err := s.repo.CompleteAttempt(ctx, CompleteAttemptRepoParams{
		AttemptID:   attemptID,
		Status:      domain.TestAttemptStatusCompleted,
		GradeStatus: gradeStatus,
		CompletedAt: time.Now(),
		TotalScore:  totalScore,
		Scores:      scores,
	})
	if err != nil {
		return domain.TestAttempt{}, fmt.Errorf("repo.CompleteAttempt: %w", err)
	}
	return completed, nil
}

func (s *Service) GradeAttempt(ctx context.Context, params GradeAttemptServiceParams) (domain.TestAttempt, error) {
	attempt, err := s.repo.GetAttempt(ctx, params.AttemptID)
	if err != nil {
		return domain.TestAttempt{}, fmt.Errorf("repo.GetAttempt: %w", err)
	}

	current, err := s.repo.GetAttemptScores(ctx, params.AttemptID)
	if err != nil {
		return domain.TestAttempt{}, fmt.Errorf("repo.GetAttemptScores: %w", err)
	}

	totalScore := mergeAndSumScores(current, params.Scores)

	completedAt := time.Now()
	if attempt.CompletedAt != nil {
		completedAt = *attempt.CompletedAt
	}

	graded, err := s.repo.GradeAttempt(ctx, GradeAttemptRepoParams{
		AttemptID:   params.AttemptID,
		GradeStatus: domain.TestGradeStatusReady,
		CompletedAt: completedAt,
		TotalScore:  totalScore,
		Scores:      params.Scores,
	})
	if err != nil {
		return domain.TestAttempt{}, fmt.Errorf("repo.GradeAttempt: %w", err)
	}
	return graded, nil
}

func (s *Service) RunExpiredAttemptsWorker(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = time.Minute
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.submitExpiredAttempts(ctx); err != nil && s.logger != nil {
				s.logger.Error().Err(err).Msg("failed to submit expired test attempts")
			}
		}
	}
}

func (s *Service) submitExpiredAttempts(ctx context.Context) error {
	ids, err := s.repo.GetExpiredAttemptIDs(ctx, time.Now())
	if err != nil {
		return fmt.Errorf("repo.GetExpiredAttemptIDs: %w", err)
	}
	for _, id := range ids {
		if _, err = s.SubmitAttempt(ctx, id); err != nil {
			return fmt.Errorf("submit attempt %d: %w", id, err)
		}
	}
	return nil
}

func (s *Service) selectQuestionsForAttempt(
	ctx context.Context,
	testID int64,
	limit int,
	settings domain.TestGenerationSettings,
) ([]domain.Question, error) {
	repoLimit := 0
	if !settings.ShuffleQuestions {
		repoLimit = limit
	}

	ids, err := s.repo.SelectBankQuestionIDs(ctx, SelectBankQuestionIDsRepoParams{
		TestID:           testID,
		DifficultyLevels: settings.DifficultyLevels,
		Limit:            repoLimit,
	})
	if err != nil {
		return nil, fmt.Errorf("repo.SelectBankQuestionIDs: %w", err)
	}
	if len(ids) == 0 {
		return nil, nil
	}

	if settings.ShuffleQuestions {
		rand.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
	}
	if len(ids) > limit {
		ids = ids[:limit]
	}

	questions, err := s.repo.GetQuestionsWithAnswers(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("repo.GetQuestionsWithAnswers: %w", err)
	}

	questionByID := make(map[int64]domain.Question, len(questions))
	for _, q := range questions {
		questionByID[q.ID] = q
	}
	ordered := make([]domain.Question, 0, len(ids))
	for _, id := range ids {
		if q, ok := questionByID[id]; ok {
			ordered = append(ordered, q)
		}
	}
	return ordered, nil
}

func decodeGenerationSettings(raw domain.JSONB) (domain.TestGenerationSettings, error) {
	settings := domain.TestGenerationSettings{ShuffleQuestions: true}
	if len(raw) == 0 {
		return settings, nil
	}
	if err := json.Unmarshal(raw, &settings); err != nil {
		return domain.TestGenerationSettings{}, fmt.Errorf("decode generation settings: %w", err)
	}
	return settings, nil
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return defaultAttemptPageLimit
	}
	return limit
}

func evaluateAttempt(questions []domain.AttemptQuestionState) ([]AttemptQuestionScore, int, bool) {
	scores := make([]AttemptQuestionScore, 0, len(questions))
	total := 0
	needsGrading := false
	for _, question := range questions {
		score := 0
		switch question.Question.Type {
		case domain.TestingQuestionTypeSingle:
			score = scoreSingle(question)
		case domain.TestingQuestionTypeMultiple:
			score = scoreMultiple(question)
		case domain.TestingQuestionTypeText:
			needsGrading = true
		}
		total += score
		scores = append(scores, AttemptQuestionScore{
			AttemptQuestionID: question.ID,
			ScoreAwarded:      score,
		})
	}
	return scores, total, needsGrading
}

func mergeAndSumScores(current, overrides []AttemptQuestionScore) int {
	overrideByID := make(map[int64]AttemptQuestionScore, len(overrides))
	for _, score := range overrides {
		overrideByID[score.AttemptQuestionID] = score
	}
	total := 0
	for _, score := range current {
		if override, ok := overrideByID[score.AttemptQuestionID]; ok {
			total += override.ScoreAwarded
			continue
		}
		total += score.ScoreAwarded
	}
	return total
}

func scoreSingle(question domain.AttemptQuestionState) int {
	if len(question.Response) != 1 || question.Response[0].SelectedAnswerID == nil {
		return 0
	}
	selectedID := *question.Response[0].SelectedAnswerID
	for _, answer := range question.Question.Answers {
		if answer.ID == selectedID && answer.IsCorrect {
			return question.Question.Points
		}
	}
	return 0
}

func scoreMultiple(question domain.AttemptQuestionState) int {
	correctIDs := map[int64]struct{}{}
	for _, answer := range question.Question.Answers {
		if answer.IsCorrect {
			correctIDs[answer.ID] = struct{}{}
		}
	}

	selectedIDs := map[int64]struct{}{}
	for _, response := range question.Response {
		if response.SelectedAnswerID != nil {
			selectedIDs[*response.SelectedAnswerID] = struct{}{}
		}
	}

	if len(correctIDs) != len(selectedIDs) {
		return 0
	}
	for id := range correctIDs {
		if _, ok := selectedIDs[id]; !ok {
			return 0
		}
	}
	return question.Question.Points
}
