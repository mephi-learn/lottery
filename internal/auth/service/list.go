package service

import (
	"context"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (h *authService) List(ctx context.Context) ([]*models.User, error) {
	user, err := h.repo.List(ctx)
	if err != nil {
		return nil, errors.Errorf("failed list clients: %w", err)
	}

	return user, nil
}
