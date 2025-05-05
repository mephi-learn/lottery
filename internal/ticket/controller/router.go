package controller

import (
	"context"
	"homework/internal/auth"
	"homework/internal/models"
	"homework/pkg/errors"
	"homework/pkg/log"
	"net/http"
)

type handler struct {
	service ticketService
	log     log.Logger
}

type HandlerOption func(*handler)

func NewHandler(opts ...HandlerOption) (*handler, error) {
	h := &handler{}

	for _, opt := range opts {
		opt(h)
	}

	if h.log == nil {
		return nil, errors.New("logger is missing")
	}

	return h, nil
}

func WithLogger(logger log.Logger) HandlerOption {
	return func(o *handler) {
		o.log = logger
	}
}

// WithService добавляет [ticetService] в обработчик запросов.
func WithService(svc ticketService) HandlerOption {
	return func(o *handler) {
		o.service = svc
	}
}

type ticketService interface {
	CreateTickets(ctx context.Context, drawId int, num int) ([]*models.Ticket, error)
	ListDrawTickets(ctx context.Context, drawId int) ([]*models.Ticket, error)
	GetTicketById(ctx context.Context, ticketId int) (*models.Ticket, error)
	AddTicket(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error)
	// ListTickets(ctx context.Context, userId int) ([]*models.Ticket, error)
}

type RouteOption func(*handler)

func (h *handler) WithRouter(mux *http.ServeMux) {
	// админ создает множество билетов
	mux.Handle("POST /api/admin/tickets/draws/{drawId}/generate/{num}", auth.Authenticated(h.CreateTickets))

	// TODO: USER получает информацию по билету
	mux.Handle("GET /api/tickets/{ticketId}", auth.Authenticated(h.GetTicketById))

	// TODO: USER получает список своих билетов
	// mux.Handle("GET /api/tickets", auth.Authenticated(h.ListAvailableTickets))
}
