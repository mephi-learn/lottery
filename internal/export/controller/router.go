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
	service exportService
	log     log.Logger
}

type exportService interface {
	ExportDraws(ctx context.Context) (*models.DrawExportResults, error)
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

// WithService добавляет [exportService] в обработчик запросов.
func WithService(svc exportService) HandlerOption {
	return func(o *handler) {
		o.service = svc
	}
}

type RouteOption func(*handler)

func (h *handler) WithRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/admin/export/draws", auth.Authenticated(h.ExportDraws))
}
