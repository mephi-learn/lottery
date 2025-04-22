package controller

import (
	"context"
	"homework/internal/auth"
	"homework/internal/models"
	"homework/pkg/errors"
	"homework/pkg/log"
	"net/http"
	"time"
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
	CreateDraw(ctx context.Context, begin time.Time, start time.Time) (drawId int, err error) // Создание тиража с указанием типа лотереи и времени старта.
	ListActiveDraw(ctx context.Context) ([]models.Draw, error)                                // Получение списка активных тиражей.
	CancelDraw(ctx context.Context, drawId int) error                                         // Отмена тиража (изменение статуса на CANCELLED).
}

type RouteOption func(*handler)

func (h *handler) WithRouter(mux *http.ServeMux) {
	mux.Handle("POST /api/admin/draws", auth.Authenticated(h.CreateDraw))
	mux.Handle("/api/admin/draws/{id}/cancel", auth.Authenticated(h.CancelDraw))
	mux.Handle("GET /api/draws/active", http.HandlerFunc(h.ListActiveDraw))
}
