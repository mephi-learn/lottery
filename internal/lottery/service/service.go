package service

import (
	"homework/internal/models"
	"homework/pkg/errors"
	"homework/pkg/log"
)

type LotteryOption func(*lotteryService) error

type lotteryService struct {
	list list

	log log.Logger
}

// NewLotteryService возвращает имплементацию сервиса для лотереи.
func NewLotteryService(opts ...LotteryOption) (*lotteryService, error) {
	var svc lotteryService
	svc.list.list = map[string]models.Lottery{}

	for _, opt := range opts {
		opt(&svc)
	}

	if svc.log == nil {
		return nil, errors.Errorf("no logger provided")
	}

	return &svc, nil
}

func WithLogger(logger log.Logger) LotteryOption {
	return func(r *lotteryService) error {
		r.log = logger
		return nil
	}
}
