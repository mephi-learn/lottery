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
	GetDraw(ctx context.Context, drawId int) (*models.DrawResultStore, error) // получение тиража
	SaveWinCombination(ctx context.Context, drawId int, winCombination []int) error // сохранение выигрышной комбинации
}

// // ResultService реализует интерфейс сервиса результатов лотереи.
// type ResultService interface {
// 	GetDrawResults(ctx context.Context, drawId int) (int, error) // Получение выигрышной комбинации тиража.
// 	// CheckTicketResult(ticketId int) (string, error) // Проверка результата билета.
// 	// LotteryByName(name string) (models.Lottery, error)
// 	// LotteryByType(name string) (models.Lottery, error)
// }

type DrawOption func(*resultService) error

type resultService struct {
	repo    Repository
	log     log.Logger
	lottery LotteryService
	// result  ResultService
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

func WithDrawLogger(logger log.Logger) DrawOption {
	return func(r *resultService) error {
		r.log = logger
		return nil
	}
}

func WithDrawRepository(repo Repository) DrawOption {
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