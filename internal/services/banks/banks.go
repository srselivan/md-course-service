package banks

import (
	"context"
	"course-service/internal/domain"
	"fmt"

	"github.com/rs/zerolog"
)

type banksRepo interface {
	Create(ctx context.Context, params CreateRepoParams) (domain.Bank, error)
	Update(ctx context.Context, params UpdateRepoParams) (domain.Bank, error)
	Delete(ctx context.Context, id int64) error
	GetList(ctx context.Context, params GetListRepoParams) ([]domain.Bank, error)
}

type Service struct {
	repo   banksRepo
	logger *zerolog.Logger
}

func NewService(repo banksRepo, logger *zerolog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, params CreateServiceParams) (domain.Bank, error) {
	bank, err := s.repo.Create(ctx, CreateRepoParams{
		Title:       params.Title,
		Description: params.Description,
		CourseId:    params.CourseId,
	})
	if err != nil {
		return domain.Bank{}, fmt.Errorf("repo.Create: %w", err)
	}
	return bank, nil
}

func (s *Service) Update(ctx context.Context, params UpdateServiceParams) (domain.Bank, error) {
	bank, err := s.repo.Update(ctx, UpdateRepoParams{
		ID:          params.ID,
		Title:       params.Title,
		Description: params.Description,
		CourseId:    params.CourseId,
	})
	if err != nil {
		return domain.Bank{}, fmt.Errorf("repo.Update: %w", err)
	}
	return bank, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("repo.Delete: %w", err)
	}
	return nil
}

func (s *Service) GetList(ctx context.Context, params GetListServiceParams) ([]domain.Bank, error) {
	banksList, err := s.repo.GetList(ctx, GetListRepoParams{
		CourseId: params.CourseId,
	})
	if err != nil {
		return nil, fmt.Errorf("repo.GetList: %w", err)
	}
	return banksList, nil
}
