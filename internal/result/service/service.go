package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
	"homework/pkg/log"
)

// LotteryService реализует интерфейс сервиса лотереи.
type LotteryService interface {
	LotteryByType(name string) (models.Lottery, error)
}

// Repository реализует интерфейс репозитория результатов тиража.
type Repository interface {
	GetDraw(ctx context.Context, drawId int) (*models.DrawResultStore, error)             // получение тиража
	SaveWinCombination(ctx context.Context, drawId int, winCombination []int) error       // сохранение выигрышной комбинации
	GetUserTicket(ctx context.Context, ticketId, userId int) (*models.TicketStore, error) // получение билета пользователя
	GetUserTickets(ctx context.Context, userId int) ([]models.TicketStore, error)         // получение билетов пользователя
}

type DrawService interface {
	PlannedDraw(ctx context.Context, drawId int) error
	ActiveDraw(ctx context.Context, drawId int) error
	CompletedDraw(ctx context.Context, drawId int) error
	CancelDraw(ctx context.Context, drawId int) error
	FailedDraw(ctx context.Context, drawId int) error
	Drawing(ctx context.Context, drawId int, combination []int) (*models.DrawingResult, error)
}

type DrawOption func(*resultService) error

type resultService struct {
	repo    Repository
	draw    DrawService
	log     log.Logger
	lottery LotteryService
}

// NewResultService возвращает имплементацию сервиса для получения результатов тиража.
func NewResultService(opts ...DrawOption) (*resultService, error) {
	var svc resultService

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

func WithResultLogger(logger log.Logger) DrawOption {
	return func(r *resultService) error {
		r.log = logger
		return nil
	}
}

func WithResultRepository(repo Repository) DrawOption {
	return func(r *resultService) error {
		r.repo = repo
		return nil
	}
}

func WithLotteryService(lottery LotteryService) DrawOption {
	return func(r *resultService) error {
		r.lottery = lottery
		return nil
	}
}

func WithDrawService(draw DrawService) DrawOption {
	return func(r *resultService) error {
		r.draw = draw
		return nil
	}
}
