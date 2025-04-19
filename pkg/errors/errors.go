package errors

import (
	"fmt"

	"github.com/go-errors/errors"
)

// New создаёт новую ошибку на основе чего угодно.
//
//go:noinline // Inlining breaks error stack traces.
func New(e any) error {
	return &wrappedError{
		err: errors.Wrap(e, 1),
	}
}

// wrap обёртывает что угодно, превращая в ошибку, второй параметр нужен, чтобы ограничить вложенность стека. Не используется, оставлено для примера.
func wrap(e any, skip int) error { //nolint:unparam
	return &wrappedError{
		err: errors.Wrap(e, skip+1),
	}
}

// Errorf оборачивает ошибку, расширяя её.
//
//go:noinline // Inlining breaks error stack traces.
func Errorf(format string, a ...any) error {
	return &wrappedError{
		err: errors.Wrap(fmt.Errorf(format, a...), 1),
	}
}

// Is сравнивает ошибку со значением, при этом проверяются все обёрнутые ошибки.
func Is(e error, original error) bool {
	return errors.Is(e, original)
}

// As проверяет, относится ли ошибка к конкретному типу.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Unwrap возвращает ошибку, которая была обёрнута.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Join возвращает конкатенированную ошибку, которая была обёрнута.
func Join(errs ...error) error {
	return errors.Join(errs...)
}
