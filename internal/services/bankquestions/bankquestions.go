package bankquestions

import (
	"context"
	"course-service/internal/domain"
	"fmt"

	"github.com/rs/zerolog"
)

type bankQuestionsRepo interface {
	Create(ctx context.Context, params CreateRepoParams) (domain.BankQuestion, error)
	Update(ctx context.Context, params UpdateRepoParams) (domain.BankQuestion, error)
	Delete(ctx context.Context, id int64) error
	GetList(ctx context.Context, params GetListRepoParams) ([]domain.BankQuestion, error)
}

type Service struct {
	repo   bankQuestionsRepo
	logger *zerolog.Logger
}

func NewService(repo bankQuestionsRepo, logger *zerolog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, params CreateServiceParams) (domain.BankQuestion, error) {
	question, err := s.repo.Create(ctx, CreateRepoParams{
		QuestionText:  params.QuestionText,
		QuestionType:  params.QuestionType,
		DefaultPoints: params.DefaultPoints,
		BankId:        params.BankId,
		Answers:       params.Answers,
	})
	if err != nil {
		return domain.BankQuestion{}, fmt.Errorf("repo.Create: %w", err)
	}
	return question, nil
}

func (s *Service) Update(ctx context.Context, params UpdateServiceParams) (domain.BankQuestion, error) {
	question, err := s.repo.Update(ctx, UpdateRepoParams{
		ID:            params.ID,
		QuestionText:  params.QuestionText,
		QuestionType:  params.QuestionType,
		DefaultPoints: params.DefaultPoints,
		BankId:        params.BankId,
		Answers:       params.Answers,
	})
	if err != nil {
		return domain.BankQuestion{}, fmt.Errorf("repo.Update: %w", err)
	}
	return question, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("repo.Delete: %w", err)
	}
	return nil
}

func (s *Service) GetList(ctx context.Context, params GetListServiceParams) ([]domain.BankQuestion, error) {
	list, err := s.repo.GetList(ctx, GetListRepoParams{
		BankId: params.BankId,
	})
	if err != nil {
		return nil, fmt.Errorf("repo.GetList: %w", err)
	}
	return list, nil
}
