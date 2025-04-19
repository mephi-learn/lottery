package server

import (
	"homework/pkg/log"
	"net"
	"net/http"
	"time"
)

type Controller interface {
	WithRouter(mux *http.ServeMux)
}

type Config struct {
	Addr            string        `yaml:"addr"             json:"addr"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" json:"shutdown_timeout"`
}

// Option позволяет настроить репозиторий добавлением новых функциональных опций.
type Option func(*options) error

type options struct {
	addr            net.Addr
	controllers     []Controller
	log             log.Logger
	shutdownTimeout time.Duration
}

func (o *options) apply(opts []Option) error {
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return err
		}
	}

	return nil
}

func WithLogger(logger log.Logger) Option {
	return func(r *options) error {
		r.log = logger

		return nil
	}
}

func WithController(controller Controller) Option {
	return func(r *options) error {
		r.controllers = append(r.controllers, controller)

		return nil
	}
}
