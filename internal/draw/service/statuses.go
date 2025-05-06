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

	s.log.InfoContext(ctx, "start change status for draw", "status", "canceled", "draw_id", drawId)
	err = s.repo.CancelDraw(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to change status draw", "status", "planned", "error", err, "draw_id", drawId)
		return errors.Errorf("failed change status draw to cancel: %w", err)
	}
	s.log.InfoContext(ctx, "draw is canceled", "draw_id", drawId)

	return nil
}

func (s *drawService) PlannedDraw(ctx context.Context, drawId int) error {
	user, err := models.UserFromContext(ctx)
	if err != nil {
		return errors.Errorf("authentificate need: %w", err)
	}

	if !user.Admin {
		return errors.Errorf("permnission denied, admin only area")
	}

	s.log.InfoContext(ctx, "start change status for draw", "status", "planned", "draw_id", drawId)
	err = s.repo.PlannedDraw(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to change status draw", "status", "planned", "error", err, "draw_id", drawId)
		return errors.Errorf("failed change status draw to planned: %w", err)
	}
	s.log.InfoContext(ctx, "draw is planned", "draw_id", drawId)

	return nil
}

func (s *drawService) ActiveDraw(ctx context.Context, drawId int) error {
	user, err := models.UserFromContext(ctx)
	if err != nil {
		return errors.Errorf("authentificate need: %w", err)
	}

	if !user.Admin {
		return errors.Errorf("permnission denied, admin only area")
	}

	s.log.InfoContext(ctx, "start change status draw", "status", "active", "draw_id", drawId)
	err = s.repo.ActiveDraw(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to change status draw", "status", "active", "error", err, "draw_id", drawId)
		return errors.Errorf("failed change status draw to active: %w", err)
	}
	s.log.InfoContext(ctx, "draw is active", "draw_id", drawId)

	return nil
}

func (s *drawService) CompletedDraw(ctx context.Context, drawId int) error {
	user, err := models.UserFromContext(ctx)
	if err != nil {
		return errors.Errorf("authentificate need: %w", err)
	}

	if !user.Admin {
		return errors.Errorf("permnission denied, admin only area")
	}

	s.log.InfoContext(ctx, "start change status draw", "status", "completed", "draw_id", drawId)
	err = s.repo.CompletedDraw(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to change status draw", "status", "completed", "error", err, "draw_id", drawId)
		return errors.Errorf("failed change status draw to completed: %w", err)
	}
	s.log.InfoContext(ctx, "draw is completed", "draw_id", drawId)

	return nil
}

func (s *drawService) FailedDraw(ctx context.Context, drawId int) error {
	user, err := models.UserFromContext(ctx)
	if err != nil {
		return errors.Errorf("authentificate need: %w", err)
	}

	if !user.Admin {
		return errors.Errorf("permnission denied, admin only area")
	}

	s.log.InfoContext(ctx, "start change status draw", "status", "failed", "draw_id", drawId)
	err = s.repo.FailedDraw(ctx, drawId)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to change status draw", "status", "failed", "error", err, "draw_id", drawId)
		return errors.Errorf("failed change status draw to failed: %w", err)
	}
	s.log.InfoContext(ctx, "draw is failed", "draw_id", drawId)

	return nil
}
