package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (s *drawService) ListActiveDraws(ctx context.Context) ([]*models.DrawOutput, error) {
	s.log.InfoContext(ctx, "start list active draws")
	defer s.log.InfoContext(ctx, "end list active draws")

	// Получаем список активных тиражей из хранилища
	list, err := s.repo.ListActiveDraw(ctx)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to list active draws from storage", "error", err)
		return nil, errors.Errorf("failed to list active draws: %w", err)
	}

	// Перекладываем данные по тиражам из формата хранения в формат выдачи результата
	resp := make([]*models.DrawOutput, len(list))
	for i, draw := range list {
		lottery, err := s.lottery.LotteryByType(draw.LotteryType)
		if err != nil {
			s.log.ErrorContext(ctx, "failed to detect lottery", "error", err)
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
	s.log.InfoContext(ctx, "start list draws")
	defer s.log.InfoContext(ctx, "end list draws")

	// Получаем список завершённых тиражей из хранилища
	list, err := s.repo.ListCompletedDraw(ctx)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to list completed draws from storage", "error", err)
		return nil, errors.Errorf("failed to list completed draws: %w", err)
	}

	// Перекладываем данные по тиражам из формата хранения в формат выдачи результата
	resp := make([]*models.DrawOutput, len(list))
	for i, draw := range list {
		lottery, err := s.lottery.LotteryByType(draw.LotteryType)
		if err != nil {
			s.log.ErrorContext(ctx, "failed to detect lottery", "error", err)
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
