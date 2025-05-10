package controller

import (
	"context"
	"homework/internal/auth"
	"homework/internal/models"
	"homework/pkg/errors"
	"homework/pkg/log"
	"net/http"
)

type handler struct {
	service authService
	log     log.Logger
}

type HandlerOption func(*handler)

func NewHandler(opts ...HandlerOption) (*handler, error) {
	h := &handler{}

	for _, opt := range opts {
		opt(h)
	}

	if h.log == nil {
		return nil, errors.New("logger is missing")
	}

	return h, nil
}

func WithLogger(logger log.Logger) HandlerOption {
	return func(o *handler) {
		o.log = logger
	}
}

// WithService добавляет [authService] в обработчик запросов.
func WithService(svc authService) HandlerOption {
	return func(o *handler) {
		o.service = svc
	}
}

type authService interface {
	SignUp(ctx context.Context, userData *models.SignUpInput) (userId int, err error)         // Регистрация пользователя.
	SignIn(ctx context.Context, userData *models.SignInInput) (signedToken string, err error) // Аутентификация пользователя.
	List(ctx context.Context) ([]*models.User, error)                                         // Список пользователей.
}

type RouteOption func(*handler)

func (h *handler) WithRouter(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/sign-up", h.signUp)
	mux.HandleFunc("POST /api/auth/sign-in", h.signIn)
	mux.HandleFunc("GET /api/admin/auth/list", auth.AuthenticatedAdmin(h.list))
}
