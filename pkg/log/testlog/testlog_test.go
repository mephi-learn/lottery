package testlog

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	logger := New(t)

	require.NotNil(t, logger)
}

func TestDiscardHandl(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	// Просто увеличиваем покрытие, потому что discardHandler - это просто заглушка.
	var handler discardHandler

	enabled := handler.Enabled(context.Background(), slog.LevelDebug)
	require.False(enabled)

	err := handler.Handle(context.Background(), slog.Record{})
	require.NoError(err)

	hattrs := handler.WithAttrs(nil)
	require.Equal(handler, hattrs)

	hgroup := handler.WithGroup("")
	require.Equal(handler, hgroup)
}
