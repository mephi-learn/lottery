package errors

import (
	"log/slog"
)

// LogValue конвертирует error в [slog.Value], который позволяет печатать трассировки стека ошибок.
func (e *wrappedError) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("message", e.Error()),
		slog.String("stack", string(e.Stack())),
	)
}
