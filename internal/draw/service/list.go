package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (s *drawService) ListActiveDraw(ctx context.Context) ([]*models.DrawOutput, error) {
	list, err := s.repo.ListActiveDraw(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to list draws: %w", err)
	}

	resp := make([]*models.DrawOutput, len(list))
	for i, draw := range list {
		lottery, err := s.lottery.LotteryByType(draw.LotteryType)
		if err != nil {
			return nil, errors.Errorf("lottery unknown type: %w", err)
		}
		status := models.DrawStatus(draw.StatusId)
		resp[i] = &models.DrawOutput{
			Id:        draw.Id,
			Status:    status.String(),
			Lottery:   lottery.Type(),
			SaleDate:  draw.SaleDate,
			StartDate: draw.StartDate,
		}
	}

	return resp, nil
}

func (s *drawService) ListCompletedDraw(ctx context.Context) ([]*models.DrawOutput, error) {
	list, err := s.repo.ListCompletedDraw(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to list draws: %w", err)
	}

	resp := make([]*models.DrawOutput, len(list))
	for i, draw := range list {
		lottery, err := s.lottery.LotteryByType(draw.LotteryType)
		if err != nil {
			return nil, errors.Errorf("lottery unknown type: %w", err)
		}
		status := models.DrawStatus(draw.StatusId)
		resp[i] = &models.DrawOutput{
			Id:        draw.Id,
			Status:    status.String(),
			Lottery:   lottery.Type(),
			SaleDate:  draw.SaleDate,
			StartDate: draw.StartDate,
		}
	}

	return resp, nil
}
