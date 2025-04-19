package errors

import (
	baseErrors "errors"
	"testing"
)

var (
	ErrWrapped    = New("wrapped error")
	ErrBase       = baseErrors.New("base error")
	ErrAnyWrapped Error
	ErrAnyBase    error
)

func TestNew(t *testing.T) {
	t.Parallel()

	err := ErrWrapped
	wrapErr := New(err)
	wrap2Err := New(wrapErr)

	baseErr := ErrBase
	wrapBaseErr := New(baseErr)
	wrap2BaseErr := New(wrapBaseErr)

	tests := []struct {
		name     string
		args     error
		is       []error
		as       []any
		notIs    []error
		typeName string
	}{
		{
			// Создание новой ошибки на основе чистой ошибки это:
			// - и есть чистая ошибка
			// - может быть приведена как к обёрнутой ошибке, так и к базовой ошибке
			// - не является обернутой чистой ошибкой, дважды обёрнутой чистой ошибкой и базовой ошибкой
			name:     "new error",
			args:     err,
			is:       []error{err},
			as:       []any{&ErrAnyWrapped, &ErrAnyBase},
			notIs:    []error{wrapErr, wrap2Err, baseErr},
			typeName: "*errors.wrappedError",
		},
		{
			// Создание новой ошибки на основе обёрнутой чистой ошибки это:
			// - и чистая ошибка и её обёртывание
			// - может быть приведена как к обёрнутой ошибке, так и к базовой ошибке
			// - не является дважды обёрнутой чистой ошибкой и базовой ошибкой
			name:     "wrapped new error",
			args:     wrapErr,
			is:       []error{wrapErr, err},
			as:       []any{&ErrAnyWrapped, &ErrAnyBase},
			notIs:    []error{wrap2Err, baseErr},
			typeName: "*errors.wrappedError",
		},
		{
			// Создание новой ошибки на основе дважды обёрнутой чистой ошибки это:
			// - и чистая ошибка и её обёртывание и её двойное обёртывание
			// - может быть приведена как к обёрнутой ошибке, так и к базовой ошибке
			// - не является базовой ошибкой
			name:     "twice wrapped error",
			args:     wrap2Err,
			is:       []error{wrap2Err, wrapErr, err},
			as:       []any{&ErrAnyWrapped, &ErrAnyBase},
			notIs:    []error{baseErr},
			typeName: "*errors.wrappedError",
		},
		{
			// Создание новой ошибки на основе базовой ошибки это:
			// - и есть базовая ошибка
			// - может быть приведена как к обёрнутой ошибке, так и к базовой ошибке
			// - не является обернутой базовой ошибкой, дважды обёрнутой базовой ошибкой и чистой ошибкой
			name:     "base error",
			args:     baseErr,
			is:       []error{baseErr},
			as:       []any{&ErrAnyWrapped, &ErrAnyBase},
			notIs:    []error{wrapBaseErr, wrap2BaseErr, wrapErr},
			typeName: "*errors.errorString",
		},
		{
			// Создание новой ошибки на основе обёрнутой базовой ошибки это:
			// - и базовая ошибка и её обёртывание
			// - может быть приведена как к обёрнутой ошибке, так и к базовой ошибке
			// - не является дважды обёрнутой чистой ошибкой и чистой ошибкой
			// - тип ошибки меняется на чистую ошибку
			name:     "wrapped base error",
			args:     wrapBaseErr,
			is:       []error{wrapBaseErr, baseErr},
			as:       []any{&ErrAnyWrapped, &ErrAnyBase},
			notIs:    []error{wrap2BaseErr, wrapErr},
			typeName: "*errors.wrappedError",
		},
		{
			// Создание новой ошибки на основе дважды обёрнутой базовой ошибки это:
			// - и базовая ошибка и её обёртывание и её двойное обёртывание
			// - может быть приведена как к обёрнутой ошибке, так и к базовой ошибке
			// - не является чистой ошибкой
			// - тип ошибки остаётся чистой ошибкой
			name:     "twice wrapped base error",
			args:     wrap2BaseErr,
			is:       []error{wrap2BaseErr, wrapBaseErr, baseErr},
			as:       []any{&ErrAnyWrapped, &ErrAnyBase},
			notIs:    []error{wrapErr},
			typeName: "*errors.wrappedError",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := New(tt.args)

			// Проверяем, является ли созданная ошибка обёртышем самой себя (поведение должно совпадать со стандартным)
			for _, is := range tt.is {
				if !Is(got, is) {
					t.Errorf("Is() = %v, want %v", got, is)
				}

				if !baseErrors.Is(got, is) {
					t.Errorf("errors.Is() = %v, want %v", got, is)
				}
			}

			// Проверяем, является ли созданная ошибка своим базовым типом (поведение должно совпадать со стандартным)
			for _, as := range tt.as {
				if !As(got, &as) {
					t.Errorf("As() = %v, want %v", got, as)
				}

				if !baseErrors.As(got, &as) {
					t.Errorf("errors.As() = %v, want %v", got, as)
				}
			}

			// Проверяем, что созданная ошибка не является чужим базовым типом (поведение должно совпадать со стандартным)
			for _, notIs := range tt.notIs {
				if Is(got, notIs) {
					t.Errorf("!Is() = %v, don't want %v", got, notIs)
				}

				if baseErrors.Is(got, notIs) {
					t.Errorf("!errors.Is() = %v, don't want %v", got, notIs)
				}
			}

			// Проверяем базовый тип ошибки
			if tt.typeName != "" {
				if tn := got.(*wrappedError).TypeName(); tn != tt.typeName { //nolint:errorlint
					t.Errorf("TypeName() = %v, want %v", tn, tt.typeName)
				}
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	t.Parallel()

	baseErr := ErrBase
	wrapBaseErr := Errorf("wrapped error: %w", baseErr)
	wrap2BaseErr := Errorf("trwice wrapped error: %w", wrapBaseErr)

	err := ErrWrapped
	wrapErr := New(err)
	wrap2Err := New(wrapErr)

	tests := []struct {
		name   string
		format error
		is     []error
		notIs  []error
	}{
		{
			name:   "base error",
			format: baseErr,
			is:     []error{baseErr},
			notIs:  []error{wrapBaseErr, wrap2BaseErr, err, wrapErr, wrap2Err},
		},
		{
			name:   "wrapped base error",
			format: wrapBaseErr,
			is:     []error{wrapBaseErr, baseErr},
			notIs:  []error{wrap2BaseErr, err, wrapErr, wrap2Err},
		},
		{
			name:   "twice wrapped base error",
			format: wrap2BaseErr,
			is:     []error{wrap2BaseErr, wrapBaseErr, baseErr},
			notIs:  []error{err, wrapErr, wrap2Err},
		},
		{
			name:   "error",
			format: err,
			is:     []error{err},
			notIs:  []error{wrapErr, wrap2Err, baseErr, wrapBaseErr, wrap2BaseErr},
		},
		{
			name:   "wrapped error",
			format: wrapErr,
			is:     []error{wrapErr, err},
			notIs:  []error{wrap2Err, baseErr, wrapBaseErr, wrap2BaseErr},
		},
		{
			name:   "twice wrapped error",
			format: wrap2Err,
			is:     []error{wrap2Err, wrapErr, err},
			notIs:  []error{baseErr, wrapBaseErr, wrap2BaseErr},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.format

			// Проверяем, является ли созданная ошибка обёртышем самой себя (поведение должно совпадать со стандартным)
			for _, is := range tt.is {
				if !Is(got, is) {
					t.Errorf("Is() = %v, want %v", got, is)
				}

				if !baseErrors.Is(got, is) {
					t.Errorf("errors.Is() = %v, want %v", got, is)
				}
			}

			// Проверяем, что созданная ошибка не является чужим базовым типом (поведение должно совпадать со стандартным)
			for _, notIs := range tt.notIs {
				if Is(tt.format, notIs) {
					t.Errorf("!Is() = %v, don't want %v", got, notIs)
				}

				if baseErrors.Is(got, notIs) {
					t.Errorf("!errors.Is() = %v, don't want %v", got, notIs)
				}
			}
		})
	}
}

func Test_wrap(t *testing.T) {
	t.Parallel()

	baseErr := ErrBase
	wrapBaseErr := wrap(baseErr, 1)
	wrap2BaseErr := wrap(wrapBaseErr, 1)

	err := ErrWrapped
	wrapErr := wrap(err, 1)
	wrap2Err := wrap(wrapErr, 1)

	tests := []struct {
		name   string
		format error
		is     []error
		notIs  []error
	}{
		{
			name:   "base error",
			format: baseErr,
			is:     []error{baseErr},
			notIs:  []error{wrapBaseErr, wrap2BaseErr, err, wrapErr, wrap2Err},
		},
		{
			name:   "wrapped base error",
			format: wrapBaseErr,
			is:     []error{wrapBaseErr, baseErr},
			notIs:  []error{wrap2BaseErr, err, wrapErr, wrap2Err},
		},
		{
			name:   "twice wrapped base error",
			format: wrap2BaseErr,
			is:     []error{wrap2BaseErr, wrapBaseErr, baseErr},
			notIs:  []error{err, wrapErr, wrap2Err},
		},
		{
			name:   "error",
			format: err,
			is:     []error{err},
			notIs:  []error{wrapErr, wrap2Err, baseErr, wrapBaseErr, wrap2BaseErr},
		},
		{
			name:   "wrapped error",
			format: wrapErr,
			is:     []error{wrapErr, err},
			notIs:  []error{wrap2Err, baseErr, wrapBaseErr, wrap2BaseErr},
		},
		{
			name:   "twice wrapped error",
			format: wrap2Err,
			is:     []error{wrap2Err, wrapErr, err},
			notIs:  []error{baseErr, wrapBaseErr, wrap2BaseErr},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.format

			// Проверяем, является ли созданная ошибка обёртышем самой себя (поведение должно совпадать со стандартным)
			for _, is := range tt.is {
				if !Is(got, is) {
					t.Errorf("Is() = %v, want %v", got, is)
				}

				if !baseErrors.Is(got, is) {
					t.Errorf("errors.Is() = %v, want %v", got, is)
				}
			}

			// Проверяем, что созданная ошибка не является чужим базовым типом (поведение должно совпадать со стандартным)
			for _, notIs := range tt.notIs {
				if Is(tt.format, notIs) {
					t.Errorf("!Is() = %v, don't want %v", got, notIs)
				}

				if baseErrors.Is(got, notIs) {
					t.Errorf("!errors.Is() = %v, don't want %v", got, notIs)
				}
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	t.Parallel()

	baseErr := ErrBase
	wrapBaseErr := wrap(baseErr, 1)
	wrap2BaseErr := wrap(wrapBaseErr, 1)

	err := ErrWrapped
	wrapErr := wrap(err, 1)
	wrap2Err := wrap(wrapErr, 1)

	tests := []struct {
		name      string
		wrapped   error
		unwrapped []error
		equal     error
		notIs     []error
	}{
		{
			name:      "base error",
			wrapped:   baseErr,
			equal:     nil,
			unwrapped: []error{nil},
		},
		{
			name:      "wrapped base error",
			wrapped:   wrapBaseErr,
			equal:     baseErr,
			unwrapped: []error{baseErr},
		},
		{
			name:      "twice wrapped base error",
			wrapped:   wrap2BaseErr,
			equal:     wrapBaseErr,
			unwrapped: []error{wrapBaseErr, baseErr},
		},
		{
			name:      "error",
			wrapped:   err,
			equal:     nil, // Так как под капотом у ошибки находится базовая ошибка, а её адреса мы не знаем, то и сравнить не можем
			unwrapped: nil, // Аналогично
		},
		{
			name:      "wrapped error",
			wrapped:   wrapErr,
			equal:     err,
			unwrapped: []error{err},
		},
		{
			name:      "twice wrapped error",
			wrapped:   wrap2Err,
			equal:     wrapErr,
			unwrapped: []error{wrapErr, err},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Unwrap(tt.wrapped)

			if tt.equal != nil {
				if got != tt.equal { //nolint:errorlint
					t.Errorf("equal = %v, want %v", got, tt.equal)
				}
			}

			// Проверяем, является ли созданная ошибка обёртыванием самой себя (поведение должно совпадать со стандартным)
			for _, is := range tt.unwrapped {
				if !Is(got, is) {
					t.Errorf("Is() = %v, want %v", got, is)
				}

				if !baseErrors.Is(got, is) {
					t.Errorf("errors.Is() = %v, want %v", got, is)
				}
			}
		})
	}
}
