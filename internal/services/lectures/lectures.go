package lectures

import (
	"context"
	"course-service/internal/domain"
	"fmt"

	"github.com/rs/zerolog"
)

type lecturesRepo interface {
	Create(ctx context.Context, params CreateRepoParams) (domain.Lecture, error)
	Update(ctx context.Context, params UpdateRepoParams) (domain.Lecture, error)
	Delete(ctx context.Context, itemId int64) error
	Get(ctx context.Context, itemId int64) (domain.Lecture, error)
	GetList(ctx context.Context, params GetListRepoParams) ([]domain.Lecture, error)
}

type Service struct {
	repo   lecturesRepo
	logger *zerolog.Logger
}

func NewService(repo lecturesRepo, logger *zerolog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, params CreateServiceParams) (domain.Lecture, error) {
	lecture, err := s.repo.Create(ctx, CreateRepoParams{
		ItemId:            params.ItemId,
		Content:           params.Content,
		VideoURL:          params.VideoURL,
		ReadingTimeMinute: params.ReadingTimeMinute,
	})
	if err != nil {
		return domain.Lecture{}, fmt.Errorf("repo.Create: %w", err)
	}
	return lecture, nil
}

func (s *Service) Update(ctx context.Context, params UpdateServiceParams) (domain.Lecture, error) {
	lecture, err := s.repo.Update(ctx, UpdateRepoParams{
		ItemId:            params.ItemId,
		Content:           params.Content,
		VideoURL:          params.VideoURL,
		ReadingTimeMinute: params.ReadingTimeMinute,
	})
	if err != nil {
		return domain.Lecture{}, fmt.Errorf("repo.Update: %w", err)
	}
	return lecture, nil
}

func (s *Service) Delete(ctx context.Context, itemId int64) error {
	if err := s.repo.Delete(ctx, itemId); err != nil {
		return fmt.Errorf("repo.Delete: %w", err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, itemId int64) (domain.Lecture, error) {
	lecture, err := s.repo.Get(ctx, itemId)
	if err != nil {
		return domain.Lecture{}, fmt.Errorf("repo.Get: %w", err)
	}
	return lecture, nil
}

func (s *Service) GetList(ctx context.Context, params GetListServiceParams) ([]domain.Lecture, error) {
	lecturesList, err := s.repo.GetList(ctx, GetListRepoParams{})
	if err != nil {
		return nil, fmt.Errorf("repo.GetList: %w", err)
	}
	return lecturesList, nil
}
