package controller

import (
	"context"
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
	// CheckTicketResult(ctx context.Context, ticketId int) error //(только USER): Проверка результата билета.
	GetDrawResults(ctx context.Context, drawId int) (int, error)                        // Получение выигрышной комбинации тиража.
}

type RouteOption func(*handler)

func (h *handler) WithRouter(mux *http.ServeMux) {
	mux.Handle("GET /api/draws/{id}/results", http.HandlerFunc(h.GetDrawResults))
	// mux.Handle("GET /api/tickets/{id}/check-result", auth.Authenticated(h.CheckTicketResult))
}
