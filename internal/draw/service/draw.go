package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
	"homework/pkg/log"
	"time"
)

// Repository реализует интерфейс репозитория тиража.
type Repository interface {
	Create(ctx context.Context, begin time.Time, start time.Time, lotteryType string) (drawId int, err error) // Создание тиража, указывается дата начала и окончания приёма билетов
	Cancel(ctx context.Context, drawId int) error                                                             // Отмена тиража, все деньги возвращаются клиентам
	SetBeginTime(ctx context.Context, drawId int, begin time.Time) error                                      // Установка времени начала продажи билетов
	SetStartTime(ctx context.Context, drawId int, start time.Time) error                                      // Установка времени начала тиража
	ListActive(ctx context.Context) ([]models.Draw, error)                                                    // Получение списка
}

// AuthService реализует интерфейс репозитория тиража.
type AuthService interface {
	GetById(ctx context.Context, userId int) (*models.User, error) // Получение информации по пользователю
}

type DrawOption func(*drawService) error

type drawService struct {
	repo Repository
	auth AuthService

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

func WithAuthService(auth AuthService) DrawOption {
	return func(r *drawService) error {
		r.auth = auth
		return nil
	}
}
