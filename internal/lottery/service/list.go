package service

import (
	"homework/internal/models"
	"homework/pkg/errors"
	"sync"
)

type list struct {
	l    sync.RWMutex
	list map[string]models.Lottery
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
