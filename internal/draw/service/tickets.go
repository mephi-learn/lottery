package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (s *drawService) CreateTickets(ctx context.Context, drawId int, num int) ([]*models.Ticket, error) {
	// Читаем существующие билеты конкретного тиража из БД, генерирую на основе правил новые билеты, сохраняю их и возвращаю списком
	// Билеты не должны повторять существующие комбинации

	// Считываем существующие билеты из БД
	ticketsIn, err := s.repo.LoadTickets(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed load tickets from repository", "error", err)
		return nil, errors.Errorf("failed load tickets from repository: %w", err)
	}

	// Получаем информацию по тиражу
	draw, err := s.repo.GetDraw(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed load draw info", "error", err)
		return nil, errors.Errorf("failed load draw info: %w", err)
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

	// Генерируем необходимое количество билетов
	tickets, err := lottery.CreateTickets(drawId, num)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to create new tickets in lottery", "error", err)
		return nil, errors.Errorf("failed to create new tickets in lottery: %w", err)
	}

	if err = s.repo.StoreTickets(ctx, tickets); err != nil {
		s.log.ErrorContext(ctx, "failed to store new tickets", "error", err)
		return nil, errors.Errorf("failed to store new tickets: %w", err)
	}

	return tickets, nil
}

func (s *drawService) ListDrawTickets(ctx context.Context, drawId int) ([]*models.Ticket, error) {
	// Читаем существующие билеты конкретного тиража из БД и возвращаем списком

	// Считываем существующие билеты из БД
	tickets, err := s.repo.LoadTickets(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed load tickets from repository", "error", err)
		return nil, errors.Errorf("failed load tickets from repository: %w", err)
	}

	return tickets, nil
}

func (s *drawService) GetTicketById(ctx context.Context, ticketId int) (*models.Ticket, error) {
	// Читаем существующий билет из БД и возвращаем его

	ticket, err := s.repo.GetTicket(ctx, ticketId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed load ticket from repository", "error", err)
		return nil, errors.Errorf("failed load ticket from repository: %w", err)
	}

	return ticket, nil
}

func (s *drawService) AddTicket(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {
	// Добавляем билет в тираж, проверяя на соответствие правилам, если
	// Билеты могут повторять существующие комбинации

	// Получаем информацию по тиражу
	draw, err := s.repo.GetDraw(ctx, ticket.DrawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed load draw info", "error", err)
		return nil, errors.Errorf("failed load draw info: %w", err)
	}

	// Создаём лотерею по её типу
	lotteryType, err := s.lottery.LotteryByType(draw.LotteryType)
	if err != nil {
		s.log.ErrorContext(ctx, "unknown lottery type", "error", err)
		return nil, errors.Errorf("unknown lottery type: %w", err)
	}
	lottery := lotteryType.Create()

	// Добавляем полученный билет в лотерею, проверяя его корректность его
	if err = lottery.AddTickets([]*models.Ticket{ticket}); err != nil {
		s.log.ErrorContext(ctx, "failed add stored tickets to lottery", "error", err)
		return nil, errors.Errorf("failed add stored tickets to lottery: %w", err)
	}

	if err = s.repo.StoreTicket(ctx, ticket); err != nil {
		s.log.ErrorContext(ctx, "failed to store new tickets", "error", err)
		return nil, errors.Errorf("failed to store new tickets: %w", err)
	}

	return ticket, nil
}

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
	ticketsIn, err := s.repo.LoadTickets(ctx, drawId)
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
