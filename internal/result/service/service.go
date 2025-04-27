package service

import (
	"homework/pkg/errors"
	"homework/pkg/log"
)

// Repository реализует интерфейс репозитория тиража.
// type Repository interface {
// 	Create(ctx context.Context, draw *models.DrawStore) (drawId int, err error) // Создание тиража
// 	Cancel(ctx context.Context, drawId int) error                               // Отмена тиража, все деньги возвращаются клиентам
// 	SetSaleDate(ctx context.Context, drawId int, begin time.Time) error         // Установка времени начала продажи билетов
// 	SetStartDate(ctx context.Context, drawId int, start time.Time) error        // Установка времени начала тиража
// 	ListActive(ctx context.Context) ([]models.DrawStore, error)                 // Получение списка
// 	Get(ctx context.Context, drawId int) (*models.DrawStore, error)             // Получение информации по тиражу
// }

// // ResultService реализует интерфейс сервиса результатов лотереи.
// type ResultService interface {
// 	GetDrawResults(ctx context.Context, drawId int) (int, error) // Получение выигрышной комбинации тиража.
// 	// CheckTicketResult(ticketId int) (string, error) // Проверка результата билета.
// 	// LotteryByName(name string) (models.Lottery, error)
// 	// LotteryByType(name string) (models.Lottery, error)
// }

type DrawOption func(*resultService) error

type resultService struct {
	// repo    Repository
	// result  ResultService

	log log.Logger
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

	// if svc.repo == nil {
	// 	return nil, errors.Errorf("no repository provided")
	// }

	return &svc, nil
}

func WithDrawLogger(logger log.Logger) DrawOption {
	return func(r *resultService) error {
		r.log = logger
		return nil
	}
}

// func WithDrawRepository(repo Repository) DrawOption {
// 	return func(r *resultService) error {
// 		r.repo = repo
// 		return nil
// 	}
// }

// func WithLotteryService(result ResultService) DrawOption {
// 	return func(r *resultService) error {
// 		r.result = result
// 		return nil
// 	}
// }
