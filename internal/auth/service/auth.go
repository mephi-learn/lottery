package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"homework/internal/auth"
	"homework/pkg/errors"
	"homework/pkg/log"
)

// Repository реализует интерфейс репозитория пользователей.
type Repository interface {
	Create(ctx context.Context, user *auth.SignUpInput) (userId int, err error)                   // Создание пользователя
	GetByUsernameAndPassword(ctx context.Context, userData *auth.SignInInput) (*auth.User, error) // Получение данных пользователя по его логину и паролю
	GetById(ctx context.Context, userId int) (*auth.User, error)                                  // Получение данных пользователя по его идентификатору
	// ChangePassword(ctx context.Context, drawId int, begin time.Time) error      // Смена пароля
}

type AuthOption func(*authService) error

type authService struct {
	repo Repository

	log log.Logger
}

// NewAuthService возвращает имплементацию сервиса пользователей.
func NewAuthService(opts ...AuthOption) (*authService, error) {
	var svc authService

	for _, opt := range opts {
		opt(&svc)
	}

	if svc.log == nil {
		return nil, errors.Errorf("no logger provided")
	}

	if svc.repo == nil {
		return nil, errors.Errorf("no repository provided")
	}

	return &svc, nil
}

func WithAuthLogger(logger log.Logger) AuthOption {
	return func(r *authService) error {
		r.log = logger
		return nil
	}
}

func WithAuthRepository(repo Repository) AuthOption {
	return func(r *authService) error {
		r.repo = repo
		return nil
	}
}

func generatePasswordHash(username, password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	// Будем использовать динамическую соль (имя пользователя)
	return hex.EncodeToString([]byte(username))
}
