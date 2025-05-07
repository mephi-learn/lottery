package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
	"time"
)

const ticketLockTime = 2

func (s *ticketService) CreateTickets(ctx context.Context, drawId int, num int) ([]*models.Ticket, error) {
	// Читаем существующие билеты конкретного тиража из БД, генерирую на основе правил новые билеты, сохраняю их и возвращаю списком
	// Билеты не должны повторять существующие комбинации

	// Считываем существующие билеты из БД
	ticketsIn, err := s.repo.LoadTicketsByDrawId(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed load tickets from repository", "error", err)
		return nil, errors.Errorf("failed load tickets from repository: %w", err)
	}

	// Получаем информацию по тиражу
	draw, err := s.draw.GetDraw(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed load draw info", "error", err)
		return nil, errors.Errorf("failed load draw info: %w", err)
	}

	// Создаём лотерею по её типу
	lottery, err := s.lottery.LotteryByType(draw.LotteryType)
	if err != nil {
		s.log.ErrorContext(ctx, "unknown lottery type", "error", err)
		return nil, errors.Errorf("unknown lottery type: %w", err)
	}

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

func (s *ticketService) ListDrawTickets(ctx context.Context, drawId int) ([]*models.Ticket, error) {
	// Читаем существующие билеты конкретного тиража из БД и возвращаем списком

	// Считываем существующие билеты из БД
	tickets, err := s.repo.LoadTicketsByDrawId(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed load tickets from repository", "error", err)
		return nil, errors.Errorf("failed load tickets from repository: %w", err)
	}

	return tickets, nil
}

func (s *ticketService) GetTicketById(ctx context.Context, ticketId int) (*models.Ticket, error) {
	// Читаем существующий билет из БД и возвращаем его

	ticket, err := s.repo.GetTicketById(ctx, ticketId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed load ticket from repository", "error", err)
		return nil, errors.Errorf("failed load ticket from repository: %w", err)
	}

	return ticket, nil
}

func (s *ticketService) AddTicket(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {
	// Добавляем билет в тираж, проверяя на соответствие правилам, если
	// Билеты могут повторять существующие комбинации

	// Получаем информацию по тиражу
	draw, err := s.draw.GetDraw(ctx, ticket.DrawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed load draw info", "error", err)
		return nil, errors.Errorf("failed load draw info: %w", err)
	}

	// Создаём лотерею по её типу
	lottery, err := s.lottery.LotteryByType(draw.LotteryType)
	if err != nil {
		s.log.ErrorContext(ctx, "unknown lottery type", "error", err)
		return nil, errors.Errorf("unknown lottery type: %w", err)
	}

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

// ListAvailableTicketsByDrawId выдаёт список допустимых билетов для покупки (конкретный тираж).
func (s *ticketService) ListAvailableTicketsByDrawId(ctx context.Context, drawId int) ([]*models.Ticket, error) {
	tickets, err := s.repo.ListAvailableTicketsByDrawId(ctx, drawId)
	if err != nil {
		return nil, errors.Errorf("failed to list tickets: %w", err)
	}

	return tickets, nil
}

// CreateReservedTicket создаёт билет для лотереи из данных (номера, перечисленные через запятую) и сразу резервирует его.
func (s *ticketService) CreateReservedTicket(ctx context.Context, drawId int, combination []int) (*models.Ticket, error) {
	// Получаем аутентифицированного пользователя
	user, err := models.UserFromContext(ctx)
	if err != nil {
		return nil, errors.Errorf("authentificate need: %w", err)
	}

	// Получаем информацию по тиражу
	draw, err := s.draw.GetDraw(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed load draw info", "error", err)
		return nil, errors.Errorf("failed load draw info: %w", err)
	}

	// Создаём лотерею по её типу
	lottery, err := s.lottery.LotteryByType(draw.LotteryType)
	if err != nil {
		s.log.ErrorContext(ctx, "unknown lottery type", "error", err)
		return nil, errors.Errorf("unknown lottery type: %w", err)
	}

	ticket, err := lottery.AddTicketWithCombination(drawId, combination)
	if err != nil {
		return nil, errors.Errorf("failed create ticket: %w", err)
	}

	ticket.Status = models.TicketStatusReady
	ticket.UserId = user.ID
	ticket.LockTime = time.Now().Add(ticketLockTime * time.Minute)

	if err = s.repo.StoreTicket(ctx, ticket); err != nil {
		return nil, errors.Errorf("failed store ticket: %w", err)
	}

	return ticket, nil
}

// ReserveTicket маркирует билет зарезервированным (выставляет время окончания в поле lock_time).
func (s *ticketService) ReserveTicket(ctx context.Context, ticketId int, userId int) error {
	lockTime := time.Now().Add(ticketLockTime * time.Minute)
	return s.repo.ReserveTicket(ctx, ticketId, userId, lockTime)
}

// BoughtTicket маркирует билет купленным (стирает время окончания в поле lock_time и меняет статус на КУПЛЕН).
func (s *ticketService) BoughtTicket(ctx context.Context, ticketId int) error {
	return s.repo.MarkTicketAsBought(ctx, ticketId)
}

// CancelTicket делает билет снова доступным для покупки (стирает время окончания в поле lock_time).
func (s *ticketService) CancelTicket(ctx context.Context, ticketId int) error {
	return s.repo.CancelTicket(ctx, ticketId)
}

func (s *ticketService) StartExpiredTicketsCleaner(ctx context.Context) {
	s.log.InfoContext(ctx, "starting expired tickets cleaner")
	ticker := time.NewTicker(time.Minute)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				s.log.InfoContext(ctx, "expired tickets cleaner stopped")

				return
			case <-ticker.C:
				s.log.InfoContext(ctx, "checking for expired tickets")
				ticketIds, err := s.repo.GetExpiredTickets(ctx)
				if err != nil {
					s.log.ErrorContext(ctx, "failed to get expired tickets", "error", err)
					continue
				}

				if len(ticketIds) > 0 {
					s.log.InfoContext(ctx, "found expired tickets", "count", len(ticketIds))
				}

				for _, id := range ticketIds {
					if err := s.repo.CancelTicket(ctx, id); err != nil {
						s.log.ErrorContext(ctx, "failed to cancel expired ticket", "ticket_id", id, "error", err)
					} else {
						s.log.InfoContext(ctx, "cancelled expired ticket", "ticket_id", id)
					}
				}
			}
		}
	}()
}
