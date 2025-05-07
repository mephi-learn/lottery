package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

// CreateDraw создание тиража.
func (s *drawService) CreateDraw(ctx context.Context, draw *models.DrawInput) (int, error) {
	s.log.InfoContext(ctx, "start create draw")
	defer s.log.InfoContext(ctx, "end create draw")

	// Проверки корректности статуса
	status := models.DrawStatusFromString(draw.Status)
	if status == models.DrawStatusUnknown {
		s.log.ErrorContext(ctx, "unknown status", "status", draw.Status)
		return -1, errors.Errorf("unknown status: %s", draw.Status)
	}

	// Получаем тип лотереи по её идентификатору
	lottery, err := s.lottery.LotteryByType(draw.Lottery)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to detect lottery", "error", err)
		return -1, errors.Errorf("unknown lottery: %w", err)
	}

	// Формируем данные для тиража и создаём тираж в хранилище, получая идентификатор
	drawQuery := &models.DrawStore{
		StatusId:    int(status),
		LotteryType: lottery.Type(),
		SaleDate:    draw.SaleDate,
		StartDate:   draw.StartDate,
	}
	drawId, err := s.repo.CreateDraw(ctx, drawQuery)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to create draw in storage", "error", err)
		return -1, errors.Errorf("cannot create draw: %w", err)
	}

	return drawId, nil
}
