package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
	"homework/pkg/log"
)

// Repository реализует интерфейс репозитория сервиса покупки билетов.
type Repository interface {
	CreateInvoice(ctx context.Context, invoice models.InvoiceStore) (invoiceId int, err error) // Создание инвойса
	GetInvoice(ctx context.Context, invoiceId int) (*models.InvoiceStore, error)               // Получние инвойса по идентификатору
	GetInvoiceByTicketId(ctx context.Context, ticketId int) (*models.InvoiceStore, error)      // Получение инвойса по идентификатору билета

	DebitingFundsFromWallet(ctx context.Context, invoice float64) error // Списание средств с кошелька пользователя
	GetAmountInUserWallet(ctx context.Context) (float64, error)         // Получение суммы на кошельке пользователя
}

// TicketService реализует интерфейс сервиса лотереи.
type TicketService interface {
	ListAvailableTicketsByDrawId(ctx context.Context, drawId int) ([]*models.Ticket, error)
	CreateReservedTicket(ctx context.Context, drawId int, data string) (*models.Ticket, error)
	ReserveTicket(ctx context.Context, ticketId int, userId int) error
	BoughtTicket(ctx context.Context, ticketId int) error
	CancelTicket(ctx context.Context, ticketId int) error
	GetTicketById(ctx context.Context, ticketId int) (*models.Ticket, error)
}

type DrawService interface {
	GetDrawByTicketId(ctx context.Context, ticketId int) (*models.DrawStore, error)
}

type PaymentOption func(*paymentService) error

type paymentService struct {
	repo Repository

	ticket TicketService
	draw   DrawService

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

	if svc.draw == nil {
		return nil, errors.Errorf("no draw provided")
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

func WithTicketService(ticket TicketService) PaymentOption {
	return func(r *paymentService) error {
		r.ticket = ticket
		return nil
	}
}

func WithDrawService(draw DrawService) PaymentOption {
	return func(r *paymentService) error {
		r.draw = draw
		return nil
	}
}
