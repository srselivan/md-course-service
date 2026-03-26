package tests

import (
	"context"
	"course-service/internal/domain"
	"errors"
	"fmt"

	"github.com/samber/lo/mutable"
)

func (s *Service) generateTestQuestions(ctx context.Context, test domain.Test) ([]domain.BankQuestion, error) {
	bankQuestions, err := s.bankQuestionsService.GetList(ctx, test.BankIds)
	if err != nil {
		return nil, fmt.Errorf("bankQuestionsService.GetList: %w", err)
	}

	if len(bankQuestions) == 0 {
		return nil, errors.New("bankQuestionsService.GetList: no bank questions found")
	}

	shuffle := func(src []domain.BankQuestion) {
		// Uses the Fisher-Yates shuffle algorithm
		mutable.Shuffle(src)
	}

	// TODO: Strategy pattern

	if test.IsDistributionNeeded() {
		var questionsByWeight map[int][]domain.BankQuestion

		for _, question := range bankQuestions {
			questionsByWeight[question.DefaultPoints] = append(questionsByWeight[question.DefaultPoints], question)
		}

		out := make([]domain.BankQuestion, 0, test.QuestionsTotal)

		for weight, questions := range questionsByWeight {
			count := test.GenerationSettings.QuestionsDistribution[weight]
			if count == 0 {
				continue
			}
			if int(count) > len(questions) {
				count = int16(len(questions))
			}

			shuffle(questions)
			out = append(out, questions[:count-1]...)
		}

		shuffle(out)
		return out, nil
	}

	shuffle(bankQuestions)
	return bankQuestions[:test.QuestionsTotal-1], nil
}
