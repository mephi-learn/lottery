package logutil

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithGroup(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	const (
		group    = "group"
		key      = "key"
		expected = "value"
	)

	var stash Stash

	// empty group
	stash = WithGroup(stash, group)
	require.Len(stash.attrs, 1)

	attr := stash.attrs[0]
	require.Equal(group, attr.Key)
	require.Equal(slog.KindGroup, attr.Value.Kind())
	require.Empty(attr.Value.Group())

	stash.attrs = nil

	// single value
	stash = WithGroup(stash, group, key, expected)
	require.Len(stash.attrs, 1)

	attr = stash.attrs[0]
	require.Equal(group, attr.Key)
	require.Equal(slog.KindGroup, attr.Value.Kind())

	g := attr.Value.Group()
	require.Len(g, 1)

	require.Equal(key, g[0].Key)
	require.Equal(slog.KindString, g[0].Value.Kind())
	require.Equal(expected, g[0].Value.String())
}

func TestWithAttrs(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	const (
		expected = "value"
		key      = "key"
	)

	var stash Stash
	stash = WithAttr(stash, key, expected)

	attrs := stash.attrs
	require.Len(attrs, 1)

	attr := attrs[0]
	require.Equal(key, attr.Key)
	require.Equal(slog.KindString, attr.Value.Kind())
	require.Equal(expected, attr.Value.String())
}

func TestHookNil(t *testing.T) {
	t.Parallel()

	var hook Hook

	attrs := hook.Attrs(context.Background())
	require.Empty(t, attrs)
}

func TestHookEmpty(t *testing.T) {
	t.Parallel()

	hook := Hook{}

	attrs := hook.Attrs(context.Background())
	require.Empty(t, attrs)
	require.NotNil(t, hook)
}

func TestHook(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	const (
		expected = "value"
		key      = "key"
	)

	pimp := func(_ context.Context, stash Stash) Stash {
		return WithAttr(stash, key, expected)
	}

	var hook Hook
	hook = append(hook, pimp)

	attrs := hook.Attrs(context.Background())
	require.Len(attrs, 1)

	attr := attrs[0]
	require.Equal(key, attr.Key)
	require.Equal(slog.KindString, attr.Value.Kind())
	require.Equal(expected, attr.Value.String())
}

func TestNilEventHook(t *testing.T) {
	t.Parallel()

	var hook Hook

	hook = append(hook, nil)

	require.Panics(t, func() {
		_ = hook.Attrs(context.Background())
	})
}
