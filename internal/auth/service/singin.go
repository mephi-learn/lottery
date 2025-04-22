package service

import (
	"context"
	"encoding/json"
	"homework/internal/auth"
	"homework/internal/models"
	"homework/pkg/errors"
)

func (h *authService) SignIn(ctx context.Context, userData *models.SignInInput) (string, error) {
	userData.Password = generatePasswordHash(userData.Username, userData.Password)
	user, err := h.repo.GetByUsernameAndPassword(ctx, userData)
	if err != nil {
		return "", err
	}

	marshaledUser, err := json.Marshal(user)
	if err != nil {
		return "", errors.Errorf("failed marshal user data: %w", err)
	}

	signedToken, err := auth.GenerateJWTToken(string(marshaledUser))
	if err != nil {
		return "", errors.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
