package controller

import (
	"context"
	"homework/internal/auth"
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
	SignUp(ctx context.Context, userData *auth.SignUpInput) (userId int, err error)         // Регистрация пользователя.
	SignIn(ctx context.Context, userData *auth.SignInInput) (signedToken string, err error) // Аутентификация пользователя.
}

type RouteOption func(*handler)

func (h *handler) WithRouter(mux *http.ServeMux) {
	mux.HandleFunc("POST /sign-up", h.signUp)
	mux.HandleFunc("POST /sign-in", h.signIn)
}

func NewErrorResponse(w http.ResponseWriter, statusCode int, errors string) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(errors))
}
