package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

// CreateDraw создание тиража.
func (s *drawService) CreateDraw(ctx context.Context, draw *models.DrawInput) (int, error) {
	// Проверки корректности статуса
	status := models.DrawStatusFromString(draw.Status)
	if status == models.DrawStatusUnknown {
		return -1, errors.Errorf("unknown status: %s", draw.Status)
	}

	// Получаем тип лотереи по её имени
	lottery, err := s.lottery.LotteryByName(draw.Lottery)
	if err != nil {
		return -1, errors.Errorf("unknown lottery: %w", err)
	}

	drawQuery := &models.DrawStore{
		StatusId:    int(status),
		LotteryType: lottery.Type(),
		SaleDate:    draw.SaleDate,
		StartDate:   draw.StartDate,
	}

	// Создаём тираж в хранилище, получая идентификатор
	drawId, err := s.repo.CreateDraw(ctx, drawQuery)
	if err != nil {
		return -1, errors.Errorf("cannot create draw: %w", err)
	}

	return drawId, nil
}
