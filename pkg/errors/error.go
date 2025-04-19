package errors

import "github.com/go-errors/errors"

var (
	_ Error = (*wrappedError)(nil)
	_ error = (*wrappedError)(nil)
)

type Error interface {
	error
	Stack() []byte
	Unwrap() error
}

type wrappedError struct {
	err *errors.Error
}

// Error возвращает сообщение ошибки.
func (e *wrappedError) Error() string {
	return e.err.Error()
}

// StackLast возвращает стек вызовов последней ошибки, отформатированный так же, как go в runtime/debug.Stack().
func (e *wrappedError) StackLast() []byte {
	return e.err.Stack()
}

// Stack возвращает стек при создании/обёртывании первой ошибки, отформатированный так же, как go в runtime/debug.Stack().
func (e *wrappedError) Stack() []byte {
	var result *wrappedError

	As(e.UnwrapMax(), &result)

	return result.StackLast()
}

// ErrorStack возвращает строку, содержащую как сообщение об ошибке, так и стек вызовов.
func (e *wrappedError) ErrorStack() string {
	return e.err.ErrorStack()
}

// TypeName возвращает тип этой ошибки. например *errors.stringError.
func (e *wrappedError) TypeName() string {
	return e.err.TypeName()
}

// Unwrap возвращает завернутую ошибку (реализует API для функции As).
func (e *wrappedError) Unwrap() error {
	return e.err.Unwrap()
}

// UnwrapMax возвращает самую первую завернутую ошибку.
func (e *wrappedError) UnwrapMax() error {
	var (
		lastUnwrap error = e
		expected   Error
	)

	for unwrap := e.Unwrap(); unwrap != nil && As(unwrap, &expected); unwrap = expected.Unwrap() {
		lastUnwrap = expected
	}

	return lastUnwrap
}

// UnwrapAll возвращает оригинальную первую.
func (e *wrappedError) UnwrapAll() error {
	var (
		lastUnwrap error = e
		unwrap     error
	)

	for unwrap = Unwrap(e.err); unwrap != nil; unwrap = Unwrap(unwrap) {
		lastUnwrap = unwrap
	}

	return lastUnwrap
}
