package courses

import (
	"context"
	"fmt"
)

func (s *Service) SetListener(ctx context.Context, params SetListenerServiceParams) error {
	if len(params.GroupIds) > 0 {
		users, err := s.usersApi.GetUsersByGroupIds(ctx, params.GroupIds)
		if err != nil {
			return fmt.Errorf("usersApi.GetUsersByGroupIds: %w", err)
		}
		params.UserIds = append(params.UserIds, users...)
	}
	if err := s.repo.SetListener(ctx, SetListenerRepoParams{
		CourseId: params.CourseId,
		UserIds:  params.UserIds,
	}); err != nil {
		return fmt.Errorf("repo.SetListener: %w", err)
	}
	return nil
}

func (s *Service) GetListenersList(ctx context.Context, params GetListenersListServiceParams) ([]int64, error) {
	list, err := s.repo.GetListenersList(ctx, GetListenersListRepoParams{
		CourseId: params.CourseId,
	})
	if err != nil {
		return nil, fmt.Errorf("repo.GetListenersList: %w", err)
	}
	return list, nil
}

func (s *Service) DeleteListener(ctx context.Context, params DeleteListenerServiceParams) error {
	if len(params.GroupIds) > 0 {
		users, err := s.usersApi.GetUsersByGroupIds(ctx, params.GroupIds)
		if err != nil {
			return fmt.Errorf("usersApi.GetUsersByGroupIds: %w", err)
		}
		params.UserIds = append(params.UserIds, users...)
	}
	if err := s.repo.DeleteListener(ctx, DeleteListenerRepoParams{
		CourseId: params.CourseId,
		UserIds:  params.UserIds,
	}); err != nil {
		return fmt.Errorf("repo.DeleteListener: %w", err)
	}
	return nil
}
