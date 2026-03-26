package tests

import (
	"context"
	"course-service/internal/domain"
	"fmt"

	"github.com/rs/zerolog"
)

type testsRepo interface {
	Create(ctx context.Context, params CreateRepoParams) (domain.Test, error)
	Update(ctx context.Context, params UpdateRepoParams) (domain.Test, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (domain.Test, error)
	GetList(ctx context.Context, params GetListRepoParams) ([]domain.Test, error)
}

type bankQuestionsService interface {
	GetList(ctx context.Context, bankIds []int64) ([]domain.BankQuestion, error)
}

type Service struct {
	repo                 testsRepo
	bankQuestionsService bankQuestionsService
	logger               *zerolog.Logger
}

func NewService(repo testsRepo, logger *zerolog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, params CreateServiceParams) (domain.Test, error) {
	test, err := s.repo.Create(ctx, CreateRepoParams{
		ItemId:             params.ItemId,
		Description:        params.Description,
		EffectiveFrom:      params.EffectiveFrom,
		EffectiveTill:      params.EffectiveTill,
		DurationSec:        params.DurationSec,
		MaxAttempts:        params.MaxAttempts,
		QuestionsTotal:     params.QuestionsTotal,
		GenerationSettings: params.GenerationSettings,
	})
	if err != nil {
		return domain.Test{}, fmt.Errorf("repo.Create: %w", err)
	}
	return test, nil
}

func (s *Service) Update(ctx context.Context, params UpdateServiceParams) (domain.Test, error) {
	test, err := s.repo.Update(ctx, UpdateRepoParams{
		ID:                 params.ID,
		ItemId:             params.ItemId,
		Description:        params.Description,
		EffectiveFrom:      params.EffectiveFrom,
		EffectiveTill:      params.EffectiveTill,
		DurationSec:        params.DurationSec,
		MaxAttempts:        params.MaxAttempts,
		QuestionsTotal:     params.QuestionsTotal,
		GenerationSettings: params.GenerationSettings,
	})
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

func (s *Service) Get(ctx context.Context, id int64) (domain.Test, error) {
	test, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Test{}, fmt.Errorf("repo.Get: %w", err)
	}
	return test, nil
}

func (s *Service) GetList(ctx context.Context, params GetListServiceParams) ([]domain.Test, error) {
	tests, err := s.repo.GetList(ctx, GetListRepoParams{})
	if err != nil {
		return nil, fmt.Errorf("repo.GetList: %w", err)
	}
	return tests, nil
}

func (s *Service) GenerateTest() {

}
