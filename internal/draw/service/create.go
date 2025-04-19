package service

import (
	"context"
	"time"
)

// CreateDraw создание тиража
func (s *drawService) CreateDraw(ctx context.Context, begin time.Time, start time.Time) (int, error) {

	// Создаём тираж в хранилище, получая идентификатор
	drawId, err := s.repo.Create(ctx, begin, start, "536")
	if err != nil {
		return 0, err
	}

	return drawId, nil
}
