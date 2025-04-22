package service

import (
	"context"
	"homework/internal/draw"
	"homework/pkg/errors"
	"homework/pkg/log"
	"time"
)

// Lottery реализует интерфейс лотереи
type Lottery interface {
	Do() bool
}

// Repository реализует интерфейс репозитория тиража
type Repository interface {
	CreateDraw(ctx context.Context, begin time.Time, start time.Time, lotteryType string) (drawId int, err error) // Создание тиража, указывается дата начала и окончания приёма билетов
	CancelDraw(ctx context.Context, drawId int) error                                                             // Отмена тиража, все деньги возвращаются клиентам
	DrawSetBeginTime(ctx context.Context, drawId int, begin time.Time) error                                      // Установка времени начала продажи билетов
	DrawSetStartTime(ctx context.Context, drawId int, start time.Time) error                                      // Установка времени начала тиража
	ListActiveDraw(ctx context.Context) ([]draw.Draw, error)                                                      // Получение списка
}

type DrawOption func(*drawService) error

type drawService struct {
	repo Repository

	log log.Logger
}

// NewDrawService возвращает имплементацию сервиса для тиража.
func NewDrawService(opts ...DrawOption) (*drawService, error) {
	var svc drawService

	for _, opt := range opts {
		opt(&svc)
	}

	if svc.log == nil {
		return nil, errors.Errorf("no logger provided")
	}

	if svc.repo == nil {
		return nil, errors.Errorf("no repository provided")
	}

	return &svc, nil
}

func WithDrawLogger(logger log.Logger) DrawOption {
	return func(r *drawService) error {
		r.log = logger
		return nil
	}
}

func WithDrawRepository(repo Repository) DrawOption {
	return func(r *drawService) error {
		r.repo = repo
		return nil
	}
}
