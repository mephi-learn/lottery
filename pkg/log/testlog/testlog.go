package testlog

import (
	"context"
	"log/slog"
	"testing"
)

// discardHandler - копия [slog.discardHandler] из [golang/go#62005].
// После релиза go1.24 можно просто удалить и использовать [slog.DiscardHandler] напрямую.
//
// [golang/go#62005]: https://github.com/golang/go/issues/62005
type discardHandler struct{}

func (dh discardHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (dh discardHandler) Handle(context.Context, slog.Record) error { return nil }
func (dh discardHandler) WithAttrs(attrs []slog.Attr) slog.Handler  { return dh }
func (dh discardHandler) WithGroup(name string) slog.Handler        { return dh }

type logger = *slog.Logger

// New возвращает логер, который пишет в пустоту.
// Он должен использоваться только в тестах, поэтому принимает на вход [testing.T].
func New(t *testing.T) logger {
	t.Helper()

	// см. golang/go#62005
	return slog.New(discardHandler{})
}
