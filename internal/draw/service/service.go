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
	CreateDraw(ctx context.Context, draw *models.DrawStore) (drawId int, err error) // Создание тиража
	PlannedDraw(ctx context.Context, drawId int) error                              // Перевод тиража в статус планирования, можно покупать билеты
	ActiveDraw(ctx context.Context, drawId int) error                               // Перевод тиража в статус активного, Билеты покупать нельзя, но можно проводить розыгрыши
	CompletedDraw(ctx context.Context, drawId int) error                            // Перевод тиража в статус завершённого, можно раздавать призы и рассылать выигрыши
	CancelDraw(ctx context.Context, drawId int) error                               // Отмена тиража, все деньги возвращаются клиентам
	FailedDraw(ctx context.Context, drawId int) error                               // Перевод тиража в статус испорченного, все деньги возвращаются клиентам
	SetDrawSaleDate(ctx context.Context, drawId int, begin time.Time) error         // Установка времени начала продажи билетов
	SetDrawStartDate(ctx context.Context, drawId int, start time.Time) error        // Установка времени начала тиража
	ListActiveDraw(ctx context.Context) ([]models.DrawStore, error)                 // Получение списка
	GetDraw(ctx context.Context, drawId int) (*models.DrawStore, error)             // Получение информации по тиражу
	LoadTicketsByDrawId(ctx context.Context, drawId int) ([]*models.Ticket, error)  // Получение списка билетов по идентификатору тиража
}

// LotteryService реализует интерфейс сервиса лотереи.
type LotteryService interface {
	LotteryByName(name string) (models.Lottery, error)
	LotteryByType(name string) (models.Lottery, error)
}

type DrawOption func(*drawService) error

type drawService struct {
	repo    Repository
	lottery LotteryService

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

func WithLotteryService(lottery LotteryService) DrawOption {
	return func(r *drawService) error {
		r.lottery = lottery
		return nil
	}
}
