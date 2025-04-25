package controller

import (
	"homework/internal/auth"
	"homework/internal/ticket/service"
	"homework/pkg/errors"
	"homework/pkg/log"
	"net/http"
)

type handler struct {
	service service.Service
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

	if h.service == nil {
		return nil, errors.New("service is missing")
	}

	return h, nil
}

func WithLogger(logger log.Logger) HandlerOption {
	return func(o *handler) {
		o.log = logger
	}
}

func WithService(svc service.Service) HandlerOption {
	return func(o *handler) {
		o.service = svc
	}
}

func (h *handler) WithRouter(mux *http.ServeMux) {
	mux.Handle("POST /api/tickets", auth.Authenticated(http.HandlerFunc(h.CreateTicket)))
}
