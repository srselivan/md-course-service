package courses

import (
	"context"
	"fmt"
)

func (s *Service) SetListener(ctx context.Context, params SetListenerServiceParams) error {
	if err := s.repo.SetListener(ctx, SetListenerRepoParams{
		CourseId: params.CourseId,
		GroupId:  params.GroupId,
	}); err != nil {
		return fmt.Errorf("repo.SetListener: %w", err)
	}
	return nil
}

func (s *Service) GetListenersList(ctx context.Context, params GetListenersListRepoParams) ([]int64, error) {
	list, err := s.repo.GetListenersList(ctx, GetListenersListRepoParams{
		CourseId: params.CourseId,
	})
	if err != nil {
		return nil, fmt.Errorf("repo.GetListenersList: %w", err)
	}
	return list, nil
}
