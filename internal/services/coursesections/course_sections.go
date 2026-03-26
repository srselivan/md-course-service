package coursesections

import (
	"context"
	"course-service/internal/domain"
	"fmt"

	"github.com/rs/zerolog"
)

type courseSectionsRepo interface {
	Create(ctx context.Context, params CreateRepoParams) (domain.CourseSection, error)
	Update(ctx context.Context, params UpdateRepoParams) (domain.CourseSection, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (domain.CourseSection, error)
	GetList(ctx context.Context, params GetListRepoParams) ([]domain.CourseSection, error)
}

type Service struct {
	repo   courseSectionsRepo
	logger *zerolog.Logger
}

func NewService(repo courseSectionsRepo, logger *zerolog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, params CreateServiceParams) (domain.CourseSection, error) {
	section, err := s.repo.Create(ctx, CreateRepoParams{
		CourseId:  params.CourseId,
		ParentId:  params.ParentId,
		Title:     params.Title,
		SortOrder: params.SortOrder,
	})
	if err != nil {
		return domain.CourseSection{}, fmt.Errorf("repo.Create: %w", err)
	}
	return section, nil
}

func (s *Service) Update(ctx context.Context, params UpdateServiceParams) (domain.CourseSection, error) {
	section, err := s.repo.Update(ctx, UpdateRepoParams{
		ID:        params.ID,
		CourseId:  params.CourseId,
		ParentId:  params.ParentId,
		Title:     params.Title,
		SortOrder: params.SortOrder,
	})
	if err != nil {
		return domain.CourseSection{}, fmt.Errorf("repo.Update: %w", err)
	}
	return section, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("repo.Delete: %w", err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, id int64) (domain.CourseSection, error) {
	section, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.CourseSection{}, fmt.Errorf("repo.Get: %w", err)
	}
	return section, nil
}

func (s *Service) GetList(ctx context.Context, params GetListServiceParams) ([]domain.CourseSection, error) {
	sections, err := s.repo.GetList(ctx, GetListRepoParams{})
	if err != nil {
		return nil, fmt.Errorf("repo.GetList: %w", err)
	}
	return sections, nil
}
