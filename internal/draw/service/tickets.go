package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (s *drawService) Drawing(ctx context.Context, drawId int, combination []int) (*models.DrawingResult, error) {
	s.log.InfoContext(ctx, "start drawing")
	defer s.log.InfoContext(ctx, "end drawing")

	// Получаем данные по тиражу из хранилища
	draw, err := s.repo.GetDraw(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to get draw from storage", "error", err)
		return nil, errors.Errorf("failed load draw info: %w", err)
	}

	// Поскольку считать мы можем только успешно завершённые тиражи, то проверяем это
	if models.DrawStatus(draw.StatusId) != models.DrawStatusCompleted {
		s.log.ErrorContext(ctx, "draw is not completed status")
		return nil, errors.New("draw is not completed status")
	}

	// Считываем купленные билеты из хранилища
	ticketsIn, err := s.repo.LoadBoughtTicketsByDrawId(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed load tickets from storage", "error", err)
		return nil, errors.Errorf("failed load tickets from storage: %w", err)
	}

	// Создаём лотерею по её типу
	lottery, err := s.lottery.LotteryByType(draw.LotteryType)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to detect lottery", "error", err)
		return nil, errors.Errorf("unknown lottery type: %w", err)
	}

	// Добавляем полученные билеты в лотерею
	if err = lottery.AddTickets(ticketsIn); err != nil {
		s.log.ErrorContext(ctx, "failed add stored tickets to lottery", "error", err)
		return nil, errors.Errorf("failed add stored tickets to lottery: %w", err)
	}

	// Проверяем билеты по правилам лотереи
	result, err := lottery.Drawing(combination)
	if err != nil {
		s.log.ErrorContext(ctx, "drawing failed", "error", err)
		return nil, errors.Errorf("drawing failed: %w", err)
	}

	// Подсчитаем статистику, для каждого выигрышного типа
	stat := map[string]int{}
	for key, values := range result {
		stat[key] = len(values)
	}

	// Подготавливаем ответ и возвращаем его
	resp := &models.DrawingResult{
		WinTickets: result,
		Statistic:  stat,
	}

	return resp, nil
}
