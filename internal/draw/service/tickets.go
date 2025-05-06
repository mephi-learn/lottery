package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (s *drawService) Drawing(ctx context.Context, drawId int, combination []int) (*models.DrawingResult, error) {
	// Получаем информацию по тиражу
	draw, err := s.repo.GetDraw(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed load draw info", "error", err)
		return nil, errors.Errorf("failed load draw info: %w", err)
	}

	if models.DrawStatus(draw.StatusId) != models.DrawStatusCompleted {
		s.log.ErrorContext(ctx, "draw not planned status")
		return nil, errors.New("draw not planned status")
	}

	// Считываем существующие билеты из БД
	ticketsIn, err := s.repo.LoadTicketsByDrawId(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed load tickets from repository", "error", err)
		return nil, errors.Errorf("failed load tickets from repository: %w", err)
	}

	// Создаём лотерею по её типу
	lotteryType, err := s.lottery.LotteryByType(draw.LotteryType)
	if err != nil {
		s.log.ErrorContext(ctx, "unknown lottery type", "error", err)
		return nil, errors.Errorf("unknown lottery type: %w", err)
	}
	lottery := lotteryType.Create()

	// Добавляем полученные билеты в лотерею
	if err = lottery.AddTickets(ticketsIn); err != nil {
		s.log.ErrorContext(ctx, "failed add stored tickets to lottery", "error", err)
		return nil, errors.Errorf("failed add stored tickets to lottery: %w", err)
	}

	// Проводим тираж и, если возникли ошибки, переводим тираж в статус ошибочного
	result, err := lottery.Drawing(combination)
	if err != nil {
		s.log.ErrorContext(ctx, "drawing failed", "error", err)
		return nil, errors.Errorf("drawing failed: %w", err)
	}

	stat := map[string]int{}
	for key, values := range result {
		stat[key] = len(values)
	}

	resp := &models.DrawingResult{
		WinTickets: result,
		Statistic:  stat,
	}

	return resp, nil
}
