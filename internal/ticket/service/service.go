package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
	"homework/pkg/log"
	"time"
)

// Repository реализует интерфейс репозитория сервиса билетов.
type Repository interface {
	StoreTicket(ctx context.Context, ticket *models.Ticket) error                           // Сохранить билет в хранилище
	StoreTickets(ctx context.Context, tickets []*models.Ticket) error                       // Сохранить список билетов в хранилище
	LoadTicketsByDrawId(ctx context.Context, drawId int) ([]*models.Ticket, error)          // Получение списка билетов по идентификатору тиража
	LoadTicketsByUserId(ctx context.Context, userId int) ([]*models.Ticket, error)          // Получение списка билетов по идентификатору тиража
	GetTicketById(ctx context.Context, ticketId int) (*models.Ticket, error)                // Получение билета по его идентификатору
	ListAvailableTicketsByDrawId(ctx context.Context, drawId int) ([]*models.Ticket, error) // Получение списка билетов, доступных для покупки
	MarkTicketAsBought(ctx context.Context, ticketId int) error                             // Маркировка билета как купленного
	ReserveTicket(ctx context.Context, ticketId int, userId int, lockTime time.Time) error  // Резервирование билета
	CancelTicket(ctx context.Context, ticketId int) error                                   // Отмена резервирования билета
	GetExpiredTickets(ctx context.Context) ([]int, error)                                   // Получение списка просроченных билетов
}

// LotteryService реализует интерфейс сервиса лотереи.
type LotteryService interface {
	LotteryByName(name string) (models.Lottery, error)
	LotteryByType(name string) (models.Lottery, error)
}

// DrawService реализует интерфейс сервиса тиража.
type DrawService interface {
	GetDraw(ctx context.Context, drawId int) (*models.DrawStore, error) // Получение информации по тиражу
}

type TicketOption func(*ticketService) error

type ticketService struct {
	repo    Repository
	lottery LotteryService
	draw    DrawService

	log log.Logger
}

// NewTicketService возвращает имплементацию сервиса для генерации билетов.
func NewTicketService(opts ...TicketOption) (*ticketService, error) {
	var svc ticketService

	for _, opt := range opts {
		opt(&svc)
	}

	if svc.log == nil {
		return nil, errors.Errorf("no logger provided")
	}

	if svc.repo == nil {
		return nil, errors.Errorf("no repository provided")
	}

	if svc.lottery == nil {
		return nil, errors.Errorf("no lottery provided")
	}

	if svc.draw == nil {
		return nil, errors.Errorf("no draw provided")
	}

	return &svc, nil
}

func WithTicketLogger(logger log.Logger) TicketOption {
	return func(r *ticketService) error {
		r.log = logger
		return nil
	}
}

func WithTicketRepository(repo Repository) TicketOption {
	return func(r *ticketService) error {
		r.repo = repo
		return nil
	}
}

func WithLotteryService(lottery LotteryService) TicketOption {
	return func(r *ticketService) error {
		r.lottery = lottery
		return nil
	}
}

func WithDrawService(draw DrawService) TicketOption {
	return func(r *ticketService) error {
		r.draw = draw
		return nil
	}
}
