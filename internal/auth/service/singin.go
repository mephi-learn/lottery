package service

import (
	"context"
	"homework/internal/auth"
	"homework/pkg/errors"
	"strconv"
)

func (h *authService) SignIn(ctx context.Context, userData *auth.SignInInput) (string, error) {
	userData.Password = generatePasswordHash(userData.Username, userData.Password)
	user, err := h.repo.GetByUsernameAndPassword(ctx, userData)
	if err != nil {
		return "", err
	}

	signedToken, err := auth.GenerateJWTToken(strconv.Itoa(user.ID))
	if err != nil {
		return "", errors.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
