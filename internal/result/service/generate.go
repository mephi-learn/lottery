package service

import (
	"context"
	"crypto/rand"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (s *resultService) GenerateDrawResults(ctx context.Context, drawId int) ([]int, error) {
	draw, err := s.repo.GetDraw(ctx, drawId)
	if err != nil {
		return make([]int, 0), errors.Errorf("failed to get draw: %w", err)
	}

	if draw == nil {
		return make([]int, 0), errors.Errorf("draw not found")
	}
	if draw.DrawStatusId != int(models.DrawStatusCompleted) {
		return make([]int, 0), errors.Errorf("draw not completed")
	}

	// Check if the draw already has winning numbers
	if draw.WinCombination != nil {
		return GetWinCombSlice(draw.WinCombination), nil
	}

	lottery, err := s.lottery.LotteryByType(draw.LotteryType)
	if err != nil {
		return make([]int, 0), errors.Errorf("failed to get lottery: %w", err)
	}

	// Generate the winning numbers based on the lottery type
	winningNumbers, err := lottery.GenerateWinningCombination(rand.Reader)

	if err != nil {
		return make([]int, 0), errors.Errorf("failed to generate winning numbers: %w", err)
	}

	// save the draw result with the winning numbers
	err = s.repo.SaveWinCombination(ctx, drawId, winningNumbers)
	if err != nil {
		return make([]int, 0), errors.Errorf("failed to save winning numbers: %w", err)
	}

	return winningNumbers, nil
}
