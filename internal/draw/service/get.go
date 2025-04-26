package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (s *drawService) GetDraw(ctx context.Context, drawId int) (*models.DrawStore, error) {
	draw, err := s.repo.GetDraw(ctx, drawId)
	if err != nil {
		return nil, errors.Errorf("failed to get draw: %w", err)
	}

	return draw, nil
}

func (s *drawService) LotteryByType(lotteryType string) (models.Lottery, error) {
	return s.lottery.LotteryByType(lotteryType)
}
