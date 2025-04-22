package log

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"homework/pkg/log/logutil"
)

func TestLevelConfigString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level Level
		name  string
	}{
		{disabled, "DISABLED"},
		{Info, "INFO"},
		{Error, "ERROR"},
		{Debug, "DEBUG"},
	}

	for _, test := range tests {
		lc := levelConfig(test.level)

		t.Run(lc.String(), func(t *testing.T) {
			t.Parallel()

			name := lc.String()

			require.Equal(t, test.name, name)
		})
	}
}

func TestLevelConfigUnmarshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level Level
	}{
		{"DISABLED", disabled},
		{"disabled", disabled},
		{"INFO", Info},
		{"info", Info},
		{"ERROR", Error},
		{"error", Error},
		{"DEBUG", Debug},
		{"debug", Debug},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var lc levelConfig

			err := lc.UnmarshalText([]byte(test.name))
			require.NoError(t, err)

			require.Equal(t, test.level, Level(lc))
		})
	}
}

func TestOptionsApplyNil(t *testing.T) {
	t.Parallel()

	var opt options

	err := opt.apply(nil)
	require.NoError(t, err)
}

func TestOptionsApplyError(t *testing.T) {
	t.Parallel()

	var opt options

	stub := func(*options) error { return errors.New("stub") }

	err := opt.apply([]LoggerOption{stub})
	require.Error(t, err)
}

func TestNewEmptyConfig(t *testing.T) {
	t.Parallel()

	var config LoggerConfig

	logger, err := New(config)
	require.NoError(t, err)
	require.NotNil(t, logger)
}

func TestNewInvalidLevel(t *testing.T) {
	t.Parallel()

	var config LoggerConfig

	config.Level = slog.LevelWarn.String()

	logger, err := New(config)
	require.Error(t, err)
	require.Nil(t, logger)
}

func TestNewStdout(t *testing.T) {
	t.Parallel()

	var config LoggerConfig

	config.Stdout = &DestConfig{}

	logger, err := New(config)
	require.NoError(t, err)
	require.NotNil(t, logger)
}

func TestNewFile(t *testing.T) {
	t.Parallel()

	var config LoggerConfig

	config.File = &FileConfig{Path: os.DevNull}

	logger, err := New(config)
	require.NoError(t, err)
	require.NotNil(t, logger)
}

func TestWithEventHook(t *testing.T) {
	t.Parallel()

	stub := func(_ context.Context, stash logutil.Stash) logutil.Stash {
		return stash
	}

	var opt options
	err := opt.apply([]LoggerOption{WithEventHook(stub)})
	require.NoError(t, err)
	require.Len(t, opt.hook, 1)
}
