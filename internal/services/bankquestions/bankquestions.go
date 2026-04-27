package bankquestions

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"course-service/internal/domain"

	"github.com/rs/zerolog"
)

type bankQuestionsRepo interface {
	Create(ctx context.Context, params CreateRepoParams) (domain.BankQuestion, error)
	BulkCreate(ctx context.Context, params BulkCreateRepoParams) error
	Update(ctx context.Context, params UpdateRepoParams) (domain.BankQuestion, error)
	BulkUpdate(ctx context.Context, params BulkUpdateRepoParams) error
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

func (s *Service) BulkCreate(ctx context.Context, params BulkCreateServiceParams) error {
	if err := s.repo.BulkCreate(ctx, BulkCreateRepoParams{
		Questions: lo.Map(params.Questions, func(item CreateServiceParams, _ int) CreateRepoParams {
			return CreateRepoParams{
				QuestionText:  item.QuestionText,
				QuestionType:  item.QuestionType,
				DefaultPoints: item.DefaultPoints,
				BankId:        item.BankId,
				Answers:       item.Answers,
			}
		}),
	}); err != nil {
		return fmt.Errorf("repo.BulkCreate: %w", err)
	}
	return nil
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

func (s *Service) BulkUpdate(ctx context.Context, params BulkUpdateServiceParams) error {
	if err := s.repo.BulkUpdate(ctx, BulkUpdateRepoParams{
		Id: params.Id,
		Questions: lo.Map(params.Questions, func(item UpdateServiceParams, _ int) UpdateRepoParams {
			return UpdateRepoParams{
				QuestionText:  item.QuestionText,
				QuestionType:  item.QuestionType,
				DefaultPoints: item.DefaultPoints,
				BankId:        item.BankId,
				Answers:       item.Answers,
			}
		}),
	}); err != nil {
		return fmt.Errorf("repo.BulkUpdate: %w", err)
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
