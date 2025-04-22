package service

import (
	"context"
	"homework/internal/draw"
	"homework/pkg/errors"
)

func (s *drawService) ListActiveDraw(ctx context.Context) ([]draw.Draw, error) {
	list, err := s.repo.ListActiveDraw(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to list draws: %w", err)
	}

	return list, nil
}
