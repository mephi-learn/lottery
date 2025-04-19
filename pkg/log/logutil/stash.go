package logutil

import (
	"context"
	"log/slog"
	"sync"
)

type Stash struct {
	attrs []slog.Attr
}

func (s Stash) append(attr slog.Attr) Stash {
	s.attrs = append(s.attrs, attr)
	return s
}

func (s *Stash) reset() {
	s.attrs = s.attrs[:0]
}

// WithGroup возвращает копию [Stash] с добавленной группой значений.
// Список аргументов должен быть в виде чередующихся пар ключ-значение.
func WithGroup(stash Stash, group string, args ...any) Stash {
	return stash.append(slog.Group(group, args...))
}

// WithAttr возвращает копию [Stash] с добавленным аттрибутом.
// Для определения допустимых типов значений val см. [slog.Value].
func WithAttr[T any](stash Stash, key string, val T) Stash {
	return stash.append(slog.Any(key, val))
}

// EventHook позволяет подключать кастомный обработчик к логеру.
type EventHook = func(context.Context, Stash) Stash

// Hook содержит набор [EventHook].
type Hook []EventHook

//nolint:gochecknoglobals
var stashPool = sync.Pool{
	New: func() any {
		return Stash{attrs: make([]slog.Attr, 0, 8)} //nolint:mnd // Просто произвольное значение.
	},
}

// Attrs определяет атрибуты на основе [context.Context] из установленных hook функции.
func (h Hook) Attrs(ctx context.Context) []slog.Attr {
	if len(h) == 0 {
		return nil
	}

	// сокращаем выделение ресурсов внутри вызовов hooks, повторно используя attr slice.
	stash := stashPool.Get().(Stash) //nolint:errcheck // В stashPool не может попасть что-то еще.
	defer func() {
		stash.reset()
		stashPool.Put(stash)
	}()

	for _, hook := range h {
		stash = hook(ctx, stash)
	}

	var attrs []slog.Attr
	if len(stash.attrs) > 0 {
		// это будет единственная аллокация, содержащая только те атрибуты, которые мы фактически создали.
		attrs = make([]slog.Attr, len(stash.attrs))
		copy(attrs, stash.attrs)
	}

	return attrs
}
