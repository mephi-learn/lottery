package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
	"homework/pkg/log"
)

// Repository реализует интерфейс репозитория сервиса покупки билетов.
type Repository interface {
	CreateInvoice(ctx context.Context, invoice models.Invoice) (invoiceId int, err error) // - Обработчки бд
}

// TicketService реализует интерфейс сервиса лотереи.
type TicketService interface {
	ListAvailableTicketsByDrawId(ctx context.Context, drawId int) ([]*models.Ticket, error)
	CreateReservedTicket(ctx context.Context, drawId int, data string) (*models.Ticket, error)
	ReserveTicket(ctx context.Context, ticketId int, userId int) error
	BoughtTicket(ctx context.Context, ticketId int) error
	CancelTicket(ctx context.Context, ticketId int) error
}

type PaymentOption func(*paymentService) error

type paymentService struct {
	repo Repository

	ticket TicketService

	log log.Logger
}

// NewPaymentService возвращает имплементацию сервиса для оплаты платежей.
func NewPaymentService(opts ...PaymentOption) (*paymentService, error) {
	var svc paymentService

	for _, opt := range opts {
		opt(&svc)
	}

	if svc.log == nil {
		return nil, errors.Errorf("no logger provided")
	}

	if svc.ticket == nil {
		return nil, errors.Errorf("no ticket provided")
	}

	return &svc, nil
}

func WithPaymentLogger(logger log.Logger) PaymentOption {
	return func(r *paymentService) error {
		r.log = logger
		return nil
	}
}

func WithPaymentRepository(repo Repository) PaymentOption {
	return func(r *paymentService) error {
		r.repo = repo
		return nil
	}
}

func WithTicketService(draw TicketService) PaymentOption {
	return func(r *paymentService) error {
		r.ticket = draw
		return nil
	}
}
