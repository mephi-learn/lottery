package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (s *resultService) GetDrawResults(ctx context.Context, drawId int) ([]int, error) {
	draw, err := s.repo.GetDraw(ctx, drawId)
	if err != nil {
		return nil, errors.Errorf("failed to get draw: %w", err)
	}

	if draw == nil {
		return nil, errors.Errorf("draw not found")
	}
	if draw.DrawStatusId != int(models.DrawStatusCompleted) {
		return nil, errors.Errorf("draw not completed")
	}

	// Check if the draw already has winning numbers
	if draw.WinCombination != nil {
		return GetWinCombSlice(draw.WinCombination), nil
	}

	return nil, errors.Errorf("No winning numbers found for this draw")
}

func (s *resultService) GetDrawWinResults(ctx context.Context, drawId int) (*models.DrawingResult, error) {
	draw, err := s.repo.GetDraw(ctx, drawId)
	if err != nil {
		return nil, errors.Errorf("failed to get draw: %w", err)
	}

	if draw == nil {
		return nil, errors.Errorf("draw not found")
	}
	if draw.DrawStatusId != int(models.DrawStatusCompleted) {
		return nil, errors.Errorf("draw not completed")
	}

	combination := make([]int, len(draw.WinCombination))
	for i, digit := range draw.WinCombination {
		combination[i] = int(digit)
	}

	tickets, err := s.draw.Drawing(ctx, drawId, combination)
	if err != nil {
		return nil, errors.Errorf("failed to get draw statistic: %w", err)
	}

	return tickets, nil
}

func (s *resultService) GetCompletedDraws(ctx context.Context) ([]*models.DrawResultStore, error) {
	list, err := s.repo.GetCompletedDraws(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to list draws: %w", err)
	}

	//resp := make([]*models.DrawOutput, len(list))
	//for i, draw := range list {
	//	lottery, err := s.lottery.LotteryByType(draw.LotteryType)
	//	if err != nil {
	//		return nil, errors.Errorf("lottery unknown type: %w", err)
	//	}
	//	status := models.DrawStatus(draw.StatusId)
	//	resp[i] = &models.DrawOutput{
	//		Id:        draw.Id,
	//		Status:    status.String(),
	//		Lottery:   lottery.Type(),
	//		SaleDate:  draw.SaleDate,
	//		StartDate: draw.StartDate,
	//	}
	//}

	return list, nil

}
