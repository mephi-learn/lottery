package service

import (
	"context"
)

func (s *resultService) GetDrawResults(ctx context.Context, drawId int) (int, error) {
	// draw, err := s.repo.GetDraw(ctx, drawId)
	// if err != nil {
	// 	return nil, errors.Errorf("failed to get draw: %w", err)
	// }

	return 1, nil
}

