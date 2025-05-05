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
	service resultService
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

// WithService добавляет [resultService] в обработчик запросов.
func WithService(svc resultService) HandlerOption {
	return func(o *handler) {
		o.service = svc
	}
}

type resultService interface {
	GetDrawResults(ctx context.Context, drawId int) ([]int, error)                             // Получение выигрышной комбинации тиража.
	GenerateDrawResults(ctx context.Context, drawId int) ([]int, error)                        // Генерация результатов тиража.
	CheckTicketResult(ctx context.Context, ticketId, userId int) (*models.TicketResult, error) // Проверка результата по номеру билета.
	CheckTicketsResult(ctx context.Context, userId int) ([]models.TicketResult, error)         // Проверка результата по всем билетам пользователя.
}

type RouteOption func(*handler)

func (h *handler) WithRouter(mux *http.ServeMux) {
	mux.Handle("GET /api/draws/{id}/results", http.HandlerFunc(h.GetDrawResults))
	mux.Handle("PUT /api/draws/{id}/results/generate", auth.Authenticated(h.GenerateDrawResults))
	mux.Handle("GET /api/tickets/{id}/check-result", auth.Authenticated(h.CheckTicketResult))
	mux.Handle("GET /api/tickets/results", auth.Authenticated(h.CheckTicketsResult))
}
