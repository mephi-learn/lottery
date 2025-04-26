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
	service drawService
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

// WithService добавляет [drawService] в обработчик запросов.
func WithService(svc drawService) HandlerOption {
	return func(o *handler) {
		o.service = svc
	}
}

type drawService interface {
	CreateDraw(ctx context.Context, draw *models.DrawInput) (drawId int, err error) // Создание тиража с указанием типа лотереи и времени старта.
	ListActiveDraw(ctx context.Context) ([]*models.DrawOutput, error)               // Получение списка активных тиражей.
	CancelDraw(ctx context.Context, drawId int) error                               // Отмена тиража (изменение статуса на CANCELLED).
	GetDraw(ctx context.Context, drawId int) (*models.DrawStore, error)             // Информация по тиражу.
	LotteryByType(lotteryType string) (models.Lottery, error)                       // Получение лотереи по её типа

	CreateTickets(ctx context.Context, drawId int, num int) ([]*models.Ticket, error)
	ListTickets(ctx context.Context, drawId int) ([]*models.Ticket, error)
	GetTicket(ctx context.Context, ticketId int) (*models.Ticket, error)
	AddTicket(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error)
}

type RouteOption func(*handler)

func (h *handler) WithRouter(mux *http.ServeMux) {
	mux.Handle("POST /api/admin/draws", auth.Authenticated(h.CreateDraw))
	mux.Handle("PUT /api/admin/draws/{id}/cancel", auth.Authenticated(h.CancelDraw))
	mux.Handle("GET /api/draws/active", http.HandlerFunc(h.ListActiveDraw))
	mux.Handle("GET /api/draws/{id}", http.HandlerFunc(h.GetDraw))
	mux.Handle("GET /api/draws/test", http.HandlerFunc(h.Test))
}
