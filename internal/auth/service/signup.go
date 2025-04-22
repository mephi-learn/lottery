package service

import (
	"context"
	"homework/internal/models"
)

func (h *authService) SignUp(ctx context.Context, userData *models.SignUpInput) (int, error) {
	userData.Password = generatePasswordHash(userData.Username, userData.Password)
	return h.repo.Create(ctx, userData)
}
