package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (h *authService) GetById(ctx context.Context, userId int) (*models.User, error) {
	user, err := h.repo.GetById(ctx, userId)
	if err != nil {
		return nil, errors.Errorf("failed get user: %w", err)
	}

	return user, nil
}
