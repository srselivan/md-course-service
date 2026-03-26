package assignments

import (
	"context"
	"course-service/internal/domain"
	"fmt"

	"github.com/rs/zerolog"
)

type assignmentsRepo interface {
	Create(ctx context.Context, params CreateRepoParams) (domain.Assignment, error)
	Update(ctx context.Context, params UpdateRepoParams) (domain.Assignment, error)
	Delete(ctx context.Context, itemId int64) error
	Get(ctx context.Context, itemId int64) (domain.Assignment, error)
	GetList(ctx context.Context, params GetListRepoParams) ([]domain.Assignment, error)
}

type Service struct {
	repo   assignmentsRepo
	logger *zerolog.Logger
}

func NewService(repo assignmentsRepo, logger *zerolog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, params CreateServiceParams) (domain.Assignment, error) {
	assignment, err := s.repo.Create(ctx, CreateRepoParams{
		ItemId:       params.ItemId,
		Description:  params.Description,
		MaxScore:     params.MaxScore,
		DeadlineDays: params.DeadlineDays,
	})
	if err != nil {
		return domain.Assignment{}, fmt.Errorf("repo.Create: %w", err)
	}
	return assignment, nil
}

func (s *Service) Update(ctx context.Context, params UpdateServiceParams) (domain.Assignment, error) {
	assignment, err := s.repo.Update(ctx, UpdateRepoParams{
		ItemId:       params.ItemId,
		Description:  params.Description,
		MaxScore:     params.MaxScore,
		DeadlineDays: params.DeadlineDays,
	})
	if err != nil {
		return domain.Assignment{}, fmt.Errorf("repo.Update: %w", err)
	}
	return assignment, nil
}

func (s *Service) Delete(ctx context.Context, itemId int64) error {
	if err := s.repo.Delete(ctx, itemId); err != nil {
		return fmt.Errorf("repo.Delete: %w", err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, itemId int64) (domain.Assignment, error) {
	assignment, err := s.repo.Get(ctx, itemId)
	if err != nil {
		return domain.Assignment{}, fmt.Errorf("repo.Get: %w", err)
	}
	return assignment, nil
}

func (s *Service) GetList(ctx context.Context, params GetListServiceParams) ([]domain.Assignment, error) {
	assignmentsList, err := s.repo.GetList(ctx, GetListRepoParams{})
	if err != nil {
		return nil, fmt.Errorf("repo.GetList: %w", err)
	}
	return assignmentsList, nil
}
