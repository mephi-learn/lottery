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
	service tickerService
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

func WithService(svc tickerService) HandlerOption {
	return func(o *handler) {
		o.service = svc
	}
}

type tickerService interface {
	CreateTicket(ctx context.Context, ticket *models.TicketRequest) (ticketId string, err error)
}

func (h *handler) WithRouter(mux *http.ServeMux) {
	mux.Handle("POST /api/tickets", auth.Authenticated(http.HandlerFunc(h.CreateTicket)))
}
