//nolint:errorlint // Пакет errors специально проверяет наличие ошибок.
package errors

import (
	baseErrors "errors"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:noinline
func fn(msg string) error {
	if err := fn1("fn " + msg); err != nil {
		return Errorf("fn wrapped %w", err)
	}

	return nil
}

//go:noinline
func fn1(msg string) error {
	if err := fn2("fn1 " + msg); err != nil {
		return Errorf("fn1 wrapped %w", err)
	}

	return nil
}

//go:noinline
func fn2(msg string) error {
	return New("fn2 error: " + msg)
}

func panicFn1() error {
	panicFn2(5)

	return nil
}

func panicFn2(_ int) {
	panicFn3()
}

func panicFn3() {
	panic("panicFn")
}

func TestErr_Error(t *testing.T) {
	t.Parallel()

	errMain := New("main error")
	errBased := baseErrors.New("base error")

	tests := []struct {
		name     string
		fields   error
		want     string
		typeName string
	}{
		{
			name:   "main",
			fields: errMain,
			want:   "main error",
		},
		{
			name:   "base",
			fields: errBased,
			want:   "base error",
		},
		{
			name:   "wrapped",
			fields: fn("wrapped test"),
			want:   "fn wrapped fn1 wrapped fn2 error: fn1 fn wrapped test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := tt.fields
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErr_StackLast(t *testing.T) {
	t.Parallel()

	result := fn("hello")
	stack := string(result.(*wrappedError).StackLast())
	stacks := strings.Split(stack, "\n")
	require.Contains(t, stacks[1], "fn: return Errorf(\"fn wrapped %w\", err)",
		"Stack trace does not contain source line: 'fn: return Errorf(\"fn wrapped %%w\", err)'", stack)

	defer func() {
		err := recover()
		if err != "panicFn" {
			t.Fatal(err)
		}

		e := Errorf("hi")

		stack = string(e.(*wrappedError).StackLast())
		stacks = strings.Split(stack, "\n")
		require.Contains(t, stacks[11], "panicFn1: panicFn2(5)", "Stack trace does not contain source line: 'panicFn1: panicFn2(5))'", stack)
		require.Contains(t, stacks[0], "error_test.go:", "Stack trace does not contain file name: 'error_test.go:'", stack)
	}()

	_ = panicFn1()
}

func TestErr_Stack(t *testing.T) {
	t.Parallel()

	result := fn("hello").(Error).Stack()
	stack := string(result)
	stacks := strings.Split(stack, "\n")
	require.Contains(t, stacks[1], "fn2: return New(\"fn2 error: \" + msg)",
		"Stack trace does not contain source line: 'fn2: return New(\"fn2 error: \" + msg)", stack)

	wrapped := New(baseErrors.New("base error"))
	require.Equal(t, wrapped.(*wrappedError).StackLast(), wrapped.(Error).Stack())
}

func TestErr_ErrorStack(t *testing.T) {
	t.Parallel()

	result := fn("hello")

	require.Equal(t, result.(*wrappedError).ErrorStack(), result.(*wrappedError).TypeName()+" "+result.Error()+"\n"+string(result.(*wrappedError).StackLast()),
		"ErrorStack is in the wrong format")
}

func TestErr_TypeName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fields error
		want   string
	}{
		{
			name:   "main",
			fields: New(ErrWrapped),
			want:   "*errors.wrappedError",
		},
		{
			name:   "base",
			fields: New(ErrBase),
			want:   "*errors.errorString",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := tt.fields

			if got := e.(*wrappedError).TypeName(); got != tt.want {
				t.Errorf("TypeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErr_Unwrap(t *testing.T) {
	t.Parallel()

	err := ErrWrapped
	wrapErr := New(err)
	wrap2Err := New(wrapErr)

	tests := []struct {
		name     string
		wrapped  error
		equal    error
		notEqual error
	}{
		{
			name:     "wrap level 1",
			wrapped:  wrapErr,
			equal:    err,
			notEqual: wrapErr,
		},
		{
			name:     "wrap level 2",
			wrapped:  wrap2Err,
			equal:    wrapErr,
			notEqual: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.wrapped.(Error).Unwrap()
			if err == nil || err != tt.equal {
				t.Errorf("Unwrap() error = %v, wantErr %v", err, tt.equal)
			}

			if err == nil || err == tt.notEqual {
				t.Errorf("!Unwrap() error = %v, wantErr %v", err, tt.notEqual)
			}
		})
	}
}

func TestWrappedError_UnwrapMax(t *testing.T) {
	t.Parallel()

	err := baseErrors.New("base error")
	wrapErr := New(err)
	wrap2Err := New(wrapErr)
	wrap3Err := New(wrap2Err)

	tests := []struct {
		name     string
		wrapped  error
		equal    error
		notEqual error
	}{
		{
			name:     "wrap level 1",
			wrapped:  wrapErr,
			equal:    wrapErr,
			notEqual: nil,
		},
		{
			name:     "wrap level 2",
			wrapped:  wrap2Err,
			equal:    wrapErr,
			notEqual: err,
		},
		{
			name:     "wrap level 3",
			wrapped:  wrap3Err,
			equal:    wrapErr,
			notEqual: err,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.wrapped.(*wrappedError).UnwrapMax()
			if err != nil {
				if tt.equal != nil && *err.(*wrappedError) != *tt.equal.(*wrappedError) {
					t.Errorf("UnwrapMax() error = %v, wantErr %v", err, tt.equal)
				}

				if tt.notEqual != nil && err == tt.notEqual {
					t.Errorf("!UnwrapMax() error = %v, wantErr %v", err, tt.notEqual)
				}
			}
		})
	}
}

func TestWrappedError_UnwrapAll(t *testing.T) {
	t.Parallel()

	err := baseErrors.New("base error")
	wrapErr := New(err)
	wrap2Err := New(wrapErr)
	wrap3Err := New(wrap2Err)

	tests := []struct {
		name     string
		wrapped  error
		equal    error
		notEqual error
	}{
		{
			name:     "wrap level 1",
			wrapped:  wrapErr,
			equal:    err,
			notEqual: wrapErr,
		},
		{
			name:     "wrap level 2",
			wrapped:  wrap2Err,
			equal:    err,
			notEqual: wrapErr,
		},
		{
			name:     "wrap level 2",
			wrapped:  wrap3Err,
			equal:    err,
			notEqual: wrapErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.wrapped.(*wrappedError).UnwrapAll()
			if err == nil || err != tt.equal {
				t.Errorf("UnwrapMax() error = %v, wantErr %v", err, tt.equal)
			}

			if err == nil || err == tt.notEqual {
				t.Errorf("!UnwrapMax() error = %v, wantErr %v", err, tt.notEqual)
			}
		})
	}
}

func TestJoinReturnsNil(t *testing.T) {
	t.Parallel()

	if err := Join(); err != nil {
		t.Errorf("Join() = %v, want nil", err)
	}

	if err := Join(nil); err != nil {
		t.Errorf("Join(nil) = %v, want nil", err)
	}

	if err := Join(nil, nil); err != nil {
		t.Errorf("Join(nil, nil) = %v, want nil", err)
	}
}

func TestJoin(t *testing.T) {
	t.Parallel()

	err1 := New("err1")
	err2 := New("err2")

	for _, test := range []struct {
		errs []error
		want []error
	}{
		{errs: []error{err1}, want: []error{err1}},
		{errs: []error{err1, err2}, want: []error{err1, err2}},
		{errs: []error{err1, nil, err2}, want: []error{err1, err2}},
	} {
		got := Join(test.errs...).(interface{ Unwrap() []error }).Unwrap()
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Join(%v) = %v; want %v", test.errs, got, test.want)
		}

		if len(got) != cap(got) {
			t.Errorf("Join(%v) returns errors with len=%v, cap=%v; want len==cap", test.errs, len(got), cap(got))
		}
	}
}

func TestJoinErrorMethod(t *testing.T) {
	t.Parallel()

	err1 := New("err1")
	err2 := New("err2")

	for _, test := range []struct {
		errs []error
		want string
	}{
		{errs: []error{err1}, want: "err1"},
		{errs: []error{err1, err2}, want: "err1\nerr2"},
		{errs: []error{err1, nil, err2}, want: "err1\nerr2"},
	} {
		got := Join(test.errs...).Error()
		if got != test.want {
			t.Errorf("Join(%v).Error() = %q; want %q", test.errs, got, test.want)
		}
	}
}
