package main

import (
	"context"
	"fmt"
	"homework/pkg/log"
	"reflect"
	"runtime"
	"strings"
)

// Manager позволяет упростить рутинную операцию завершения сервисов.
type Manager struct {
	build BuildInfo

	quit chan struct{}

	ctx *context.Context
	log log.Logger
}

func NewManager(ctx *context.Context, log log.Logger) Manager {
	b := Manager{
		build: ReadBuildInfo(),
		quit:  make(chan struct{}, 1),
		ctx:   ctx,
		log:   log,
	}

	return b
}

// Run запускает блокирующую процедуру и уведомляет после завершения.
func (b Manager) run(caller func() error) {
	defer func() { b.quit <- struct{}{} }()

	// Достаём функцию, чтобы в логе было красивое имя, а не указатель.
	fn := runtime.FuncForPC(reflect.ValueOf(caller).Pointer()).Name()
	fn = strings.TrimPrefix(fn, b.build.Name+"/")

	ctx := *b.ctx
	if err := caller(); err != nil {
		b.log.ErrorContext(ctx, "unexpected error", "func", fn, "err", err)
		return
	}

	b.log.DebugContext(ctx, "execution finished", "func", fn)
}

type Shutdowner interface {
	Shutdown(ctx context.Context) error
}

// Stop выключает сервис и пишет ошибку в лог, если завершение было грязным.
func (b Manager) stop(svc Shutdowner) {
	ctx := *b.ctx

	name := fmt.Sprintf("%T", svc)

	err := svc.Shutdown(ctx)
	if err != nil {
		b.log.ErrorContext(ctx, "dirty service shutdown", "service", name, "err", err)
		return
	}

	b.log.DebugContext(ctx, "service terminated successfully", "service", name)
}
