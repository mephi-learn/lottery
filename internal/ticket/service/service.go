package service

import (
	"homework/internal/models"
	"homework/pkg/errors"
	"homework/pkg/log"
)

type TicketService interface {
	CreateTickets(drawId int, lotteryType string, num int) ([]*models.Ticket, error)
}

type TicketOption func(*ticketService) error

type ticketService struct {
	lotteryService LotteryService
	log            log.Logger
}

type LotteryService interface {
	LotteryByType(name string) (models.Lottery, error)
}

// NewTicketService возвращает имплементацию сервиса для билетов.
func NewTicketService(opts ...TicketOption) (TicketService, error) {
	var svc ticketService

	for _, opt := range opts {
		if err := opt(&svc); err != nil {
			return nil, err
		}
	}

	if svc.log == nil {
		return nil, errors.Errorf("no logger provided")
	}

	if svc.lotteryService == nil {
		return nil, errors.Errorf("no lottery service provided")
	}

	return &svc, nil
}

func WithLogger(logger log.Logger) TicketOption {
	return func(s *ticketService) error {
		s.log = logger
		return nil
	}
}

func WithLotteryService(lotteryService LotteryService) TicketOption {
	return func(s *ticketService) error {
		s.lotteryService = lotteryService
		return nil
	}
}

func (s *ticketService) CreateTickets(drawId int, lotteryType string, num int) ([]*models.Ticket, error) {
	lottery, err := s.lotteryService.LotteryByType(lotteryType)
	if err != nil {
		return nil, errors.Errorf("failed to get lottery: %w", err)
	}

	tickets, err := lottery.CreateTickets(drawId, num)
	if err != nil {
		return nil, errors.Errorf("failed to create tickets: %w", err)
	}

	return tickets, nil
}
