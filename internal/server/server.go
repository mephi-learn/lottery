package server

import (
	"homework/pkg/errors"
	"net"
	"net/http"
	"sync/atomic"
)

type Server struct {
	options

	mux *http.ServeMux

	done atomic.Bool // для сигнализаци сессиям о начале shutdown-а
}

func New(config Config, opts ...Option) (*Server, error) {
	var opt options

	opt.shutdownTimeout = config.ShutdownTimeout

	if config.Addr == "" {
		return nil, errors.New("empty listen address")
	}

	var err error
	opt.addr, err = net.ResolveTCPAddr("tcp", config.Addr)
	if err != nil {
		return nil, errors.Errorf("resolve TCP address '%s': %w", config.Addr, err)
	}

	if err := opt.apply(opts); err != nil {
		return nil, errors.Errorf("applying server options: %w", err)
	}

	var server Server

	server.options = opt
	server.done.Store(false)

	if server.log == nil {
		return nil, errors.New("logger is missing")
	}

	return &server, nil
}

// ListenAndServe запускает сервер и принимает входящие запросы.
func (s *Server) ListenAndServe() error {
	mux := &http.ServeMux{}
	for _, controller := range s.options.controllers {
		controller.WithRouter(mux)
	}
	s.mux = mux

	return http.ListenAndServe(s.options.addr.String(), s.mux)
}
