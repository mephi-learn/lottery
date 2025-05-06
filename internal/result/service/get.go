package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (s *resultService) GetDrawResults(ctx context.Context, drawId int) ([]int, error) {
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

	return make([]int, 0), errors.Errorf("No winning numbers found for this draw")
}
