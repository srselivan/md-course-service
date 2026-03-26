package coursesectionitems

import (
	"context"
	"course-service/internal/domain"
	"fmt"

	"github.com/rs/zerolog"
)

type courseSectionItemsRepo interface {
	Create(ctx context.Context, params CreateRepoParams) (domain.CourseSectionItem, error)
	Update(ctx context.Context, params UpdateRepoParams) (domain.CourseSectionItem, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (domain.CourseSectionItem, error)
	GetList(ctx context.Context, params GetListRepoParams) ([]domain.CourseSectionItem, error)
}

type Service struct {
	repo   courseSectionItemsRepo
	logger *zerolog.Logger
}

func NewService(repo courseSectionItemsRepo, logger *zerolog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, params CreateServiceParams) (domain.CourseSectionItem, error) {
	item, err := s.repo.Create(ctx, CreateRepoParams{
		SectionId:   params.SectionId,
		ItemType:    params.ItemType,
		Title:       params.Title,
		SortOrder:   params.SortOrder,
		IsPublished: params.IsPublished,
	})
	if err != nil {
		return domain.CourseSectionItem{}, fmt.Errorf("repo.Create: %w", err)
	}
	return item, nil
}

func (s *Service) Update(ctx context.Context, params UpdateServiceParams) (domain.CourseSectionItem, error) {
	item, err := s.repo.Update(ctx, UpdateRepoParams{
		ID:          params.ID,
		SectionId:   params.SectionId,
		ItemType:    params.ItemType,
		Title:       params.Title,
		SortOrder:   params.SortOrder,
		IsPublished: params.IsPublished,
	})
	if err != nil {
		return domain.CourseSectionItem{}, fmt.Errorf("repo.Update: %w", err)
	}
	return item, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("repo.Delete: %w", err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, id int64) (domain.CourseSectionItem, error) {
	item, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.CourseSectionItem{}, fmt.Errorf("repo.Get: %w", err)
	}
	return item, nil
}

func (s *Service) GetList(ctx context.Context, params GetListServiceParams) ([]domain.CourseSectionItem, error) {
	items, err := s.repo.GetList(ctx, GetListRepoParams{})
	if err != nil {
		return nil, fmt.Errorf("repo.GetList: %w", err)
	}
	return items, nil
}
