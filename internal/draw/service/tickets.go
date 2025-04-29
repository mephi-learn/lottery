package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (s *drawService) Drawing(ctx context.Context, drawId int, combination []int) (map[string][]*models.Ticket, error) {
	// Получаем информацию по тиражу
	draw, err := s.repo.GetDraw(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed load draw info", "error", err)
		return nil, errors.Errorf("failed load draw info: %w", err)
	}

	//if models.DrawStatus(draw.StatusId) != models.DrawStatusPlanned {
	//	s.log.ErrorContext(ctx, "draw not planned status")
	//	return nil, errors.New("draw not planned status")
	//}

	//// Переводим тираж в статус активного
	//if err = s.repo.ActiveDraw(ctx, drawId); err != nil {
	//	s.log.ErrorContext(ctx, "failed set draw to active", "error", err)
	//	return nil, errors.Errorf("failed set draw to active: %w", err)
	//}

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
		//if err = s.repo.FailedDraw(ctx, drawId); err != nil {
		//	s.log.ErrorContext(ctx, "failed set draw to failed", "error", err)
		//	return nil, errors.Errorf("failed set draw to failed: %w", err)
		//}

		s.log.ErrorContext(ctx, "drawing failed", "error", err)
		return nil, errors.Errorf("drawing failed: %w", err)
	}

	//// Переводим тираж в статус завершённого
	//if err = s.repo.CompletedDraw(ctx, drawId); err != nil {
	//	s.log.ErrorContext(ctx, "failed set draw to completed", "error", err)
	//	return nil, errors.Errorf("failed set draw to completed: %w", err)
	//}

	return result, nil
}
