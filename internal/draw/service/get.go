package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (s *drawService) GetDraw(ctx context.Context, drawId int) (*models.DrawStore, error) {
	s.log.InfoContext(ctx, "start get draw")
	defer s.log.InfoContext(ctx, "end get draw")

	// Получаем данные по тиражу из хранилища
	draw, err := s.repo.GetDraw(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to get draw from storage", "error", err)
		return nil, errors.Errorf("failed to get draw: %w", err)
	}

	return draw, nil
}

func (s *drawService) LotteryByType(lotteryType string) (models.Lottery, error) {
	return s.lottery.LotteryByType(lotteryType)
}

// Получение тиража по идентификатору билета
func (s *drawService) GetDrawByTicketId(ctx context.Context, ticketId int) (*models.DrawStore, error) {
	return s.repo.GetDrawByTicketId(ctx, ticketId)
}
