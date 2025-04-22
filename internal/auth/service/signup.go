package service

import (
	"context"
	"homework/internal/auth"
)

func (h *authService) SignUp(ctx context.Context, userData *auth.SignUpInput) (int, error) {
	userData.Password = generatePasswordHash(userData.Username, userData.Password)
	return h.repo.Create(ctx, userData)
}
