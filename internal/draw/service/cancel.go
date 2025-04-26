package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (s *drawService) CancelDraw(ctx context.Context, drawId int) error {
	user, err := models.UserFromContext(ctx)
	if err != nil {
		return errors.Errorf("authentificate need: %w", err)
	}

	if !user.Admin {
		return errors.Errorf("permnission denied, admin only area")
	}

	s.log.InfoContext(ctx, "start cancel draw", "draw_id", drawId)
	err = s.repo.CancelDraw(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to cancel draw", "error", err, "draw_id", drawId)
		return errors.Errorf("failed to cancel draw: %w", err)
	}
	s.log.InfoContext(ctx, "draw canceled", "draw_id", drawId)

	return nil
}
