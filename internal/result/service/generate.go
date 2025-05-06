package service

import (
	"context"
	"crypto/rand"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (s *resultService) Drawing(ctx context.Context, drawId int) ([]int, error) {
	draw, err := s.repo.GetDraw(ctx, drawId)
	if err != nil {
		return nil, errors.Errorf("failed to get draw: %w", err)
	}

	if draw == nil {
		return nil, errors.Errorf("draw not found")
	}

	// Если тираж не в состоянии "запланирован", выдаём ошибку
	if draw.DrawStatusId != int(models.DrawStatusPlanned) {
		return nil, errors.Errorf("draw not completed")
	}

	// Переводим тираж в состояние active (розыгрыш в процессе)
	if err = s.draw.ActiveDraw(ctx, drawId); err != nil {
		return nil, errors.Errorf("failed change status to active: %w", err)
	}

	lottery, err := s.lottery.LotteryByType(draw.LotteryType)
	if err != nil {
		_ = s.draw.FailedDraw(ctx, drawId)
		return nil, errors.Errorf("failed to get lottery: %w", err)
	}

	// Generate the winning numbers based on the lottery type
	winningNumbers, err := lottery.GenerateWinningCombination(rand.Reader)
	if err != nil {
		_ = s.draw.FailedDraw(ctx, drawId)
		return nil, errors.Errorf("failed to generate winning numbers: %w", err)
	}

	// save the draw result with the winning numbers
	err = s.repo.SaveWinCombination(ctx, drawId, winningNumbers)
	if err != nil {
		_ = s.draw.FailedDraw(ctx, drawId)
		return nil, errors.Errorf("failed to save winning numbers: %w", err)
	}

	_ = s.draw.CompletedDraw(ctx, drawId)

	//tickets, err := s.draw.Drawing(ctx, drawId, winningNumbers)
	//if err != nil {
	//
	//}

	return winningNumbers, nil
}
