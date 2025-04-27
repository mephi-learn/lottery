package service

import (
	"homework/internal/models"
	"homework/pkg/errors"
	"homework/pkg/log"
	"sync"
)

type LotteryService interface {
	LotteryByName(name string) (models.Lottery, error)
	LotteryByType(name string) (models.Lottery, error)
	RegisterLottery(lottery models.Lottery) error
}

type list struct {
	l    sync.RWMutex
	list map[string]models.Lottery
}

type lotteryService struct {
	list *list
	log  log.Logger
}

type LotteryOption func(*lotteryService) error

// NewLotteryService возвращает имплементацию сервиса для лотереи.
func NewLotteryService(opts ...LotteryOption) (LotteryService, error) {
	var svc lotteryService

	for _, opt := range opts {
		if err := opt(&svc); err != nil {
			return nil, err
		}
	}

	if svc.log == nil {
		return nil, errors.Errorf("no logger provided")
	}

	if svc.list == nil {
		svc.list = &list{
			list: make(map[string]models.Lottery),
		}
	}

	// Регистрируем лотерею 5/36
	if err := svc.RegisterLottery(models.NewLottery536()); err != nil {
		return nil, err
	}

	return &svc, nil
}

func WithLogger(logger log.Logger) LotteryOption {
	return func(s *lotteryService) error {
		s.log = logger
		return nil
	}
}

func (s *lotteryService) RegisterLottery(lottery models.Lottery) error {
	s.list.l.Lock()
	defer s.list.l.Unlock()

	name := lottery.Type()

	if _, ok := s.list.list[name]; ok {
		return errors.New("lottery already register")
	}

	s.list.list[name] = lottery
	return nil
}

func (s *lotteryService) LotteryByType(lotteryType string) (models.Lottery, error) {
	s.list.l.RLock()
	defer s.list.l.RUnlock()

	lottery, ok := s.list.list[lotteryType]
	if !ok {
		return nil, errors.New("lottery not found")
	}

	return lottery, nil
}

func (s *lotteryService) LotteryByName(name string) (models.Lottery, error) {
	s.list.l.RLock()
	defer s.list.l.RUnlock()

	for _, lottery := range s.list.list {
		if lottery.Name() == name {
			return lottery, nil
		}
	}

	return nil, errors.New("lottery not found")
}
