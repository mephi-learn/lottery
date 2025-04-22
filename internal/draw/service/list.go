package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (s *drawService) ListActiveDraw(ctx context.Context) ([]models.Draw, error) {
	list, err := s.repo.ListActive(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to list draws: %w", err)
	}

	return list, nil
}
