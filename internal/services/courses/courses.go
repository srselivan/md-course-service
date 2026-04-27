package courses

import (
	"context"
	"fmt"

	"course-service/internal/domain"

	"github.com/rs/zerolog"
)

type coursesRepo interface {
	Create(ctx context.Context, params CreateRepoParams) (domain.Course, error)
	Update(ctx context.Context, params UpdateRepoParams) (domain.Course, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (domain.Course, error)
	GetList(ctx context.Context, params GetListRepoParams) ([]domain.Course, error)
	GetWithAllItems(ctx context.Context, params GetWithAllItemsRepoParams) (domain.CourseWithItems, error)

	listenersRepo
}

type listenersRepo interface {
	SetListener(ctx context.Context, params SetListenerRepoParams) error
	GetListenersList(ctx context.Context, params GetListenersListRepoParams) ([]int64, error)
	DeleteListener(ctx context.Context, params DeleteListenerRepoParams) error
}

type usersApi interface {
	GetUsersByGroupIds(ctx context.Context, groupIds []int64) ([]int64, error)
}

type Service struct {
	repo     coursesRepo
	usersApi usersApi
	logger   *zerolog.Logger
}

func NewService(repo coursesRepo, logger *zerolog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, params CreateServiceParams) (domain.Course, error) {
	course, err := s.repo.Create(ctx, CreateRepoParams{
		Title:       params.Title,
		Description: params.Description,
		OwnerUserId: params.OwnerUserId,
	})
	if err != nil {
		return domain.Course{}, fmt.Errorf("repo.Create: %w", err)
	}
	return course, nil
}

func (s *Service) Update(ctx context.Context, params UpdateServiceParams) (domain.Course, error) {
	course, err := s.repo.Update(ctx, UpdateRepoParams{
		ID:          params.ID,
		Title:       params.Title,
		Description: params.Description,
		OwnerUserId: params.OwnerUserId,
		Status:      params.Status,
	})
	if err != nil {
		return domain.Course{}, fmt.Errorf("repo.Update: %w", err)
	}
	return course, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("repo.Delete: %w", err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, id int64) (domain.Course, error) {
	course, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Course{}, fmt.Errorf("repo.Get: %w", err)
	}
	return course, nil
}

func (s *Service) GetList(ctx context.Context, params GetListServiceParams) ([]domain.Course, error) {
	courses, err := s.repo.GetList(ctx, GetListRepoParams{
		Status:      params.Status,
		UserID:      params.UserID,
		OwnerUserId: params.OwnerUserId,
	})
	if err != nil {
		return nil, fmt.Errorf("repo.GetList: %w", err)
	}
	return courses, nil
}

func (s *Service) GetWithAllItems(ctx context.Context, params GetWithAllItemsParams) (domain.CourseWithItems, error) {
	courseWithItems, err := s.repo.GetWithAllItems(ctx, GetWithAllItemsRepoParams{
		Id:     params.Id,
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		return domain.CourseWithItems{}, fmt.Errorf("repo.GetWithAllItems: %w", err)
	}
	return courseWithItems, nil
}
