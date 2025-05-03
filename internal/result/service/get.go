package service

import (
	"context"
	"homework/pkg/errors"
)

func (s *resultService) GetDrawResults(ctx context.Context, drawId int) (int, error) {
	draw, err := s.repo.GetDraw(ctx, drawId)
	if err != nil {
		return 0, errors.Errorf("failed to get draw: %w", err)
	}

	return draw.DrawStatusId, nil
}

