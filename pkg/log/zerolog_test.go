package log

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"homework/pkg/log/filter"
	"homework/pkg/log/logutil"
	"log/slog"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNewLogHandlerDefaultWriter(t *testing.T) {
	t.Parallel()

	var opt options

	_, err := newLogHandler(opt)
	require.NoError(t, err)
}

func TestNewLogHandlerSingleWriter(t *testing.T) {
	t.Parallel()

	var opt options

	opt.writers = append(opt.writers, fileWriter{
		File:       os.Stdout,
		DestConfig: DestConfig{Format: FormatConsole},
	})

	_, err := newLogHandler(opt)
	require.NoError(t, err)
}

func TestNewLogHandlerMultiWriter(t *testing.T) {
	t.Parallel()

	var opt options

	opt.writers = append(opt.writers,
		fileWriter{
			File:       os.Stdout,
			DestConfig: DestConfig{Format: FormatConsole},
		},
		fileWriter{
			File:       os.Stderr,
			DestConfig: DestConfig{Format: FormatJSON},
		},
	)

	_, err := newLogHandler(opt)
	require.NoError(t, err)
}

func TestNewLogHandlerWriterError(t *testing.T) {
	t.Parallel()

	var opt options
	var err error

	opt.writers = append(opt.writers, fileWriter{File: os.Stdout})

	_, err = newLogHandler(opt)
	require.Error(t, err)

	opt.writers = append(opt.writers, fileWriter{File: os.Stderr})

	_, err = newLogHandler(opt)
	require.Error(t, err)
}

func TestNewLogHandlerDebug(t *testing.T) {
	t.Parallel()

	var opt options
	opt.Level = Debug

	h, err := newLogHandler(opt)
	require.NoError(t, err)

	require.Equal(t, zerolog.DebugLevel, h.zl.GetLevel())
}

func TestHandlerEnabled(t *testing.T) {
	t.Parallel()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	tests := []struct {
		name     string
		ours     zerolog.Level
		theirs   slog.Level
		expected bool
	}{
		{"global", zerolog.NoLevel, Info, false},
		{"disabled", zerolog.Disabled, Info, false},
		{"debug enabled", zerolog.DebugLevel, Debug, true},
		{"debug disabled", zerolog.InfoLevel, Debug, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			h := zlTestHandler(t, nil)
			h.zl = h.zl.Level(test.ours)

			enabled := h.Enabled(context.Background(), test.theirs)
			require.Equal(t, test.expected, enabled)
		})
	}
}

func TestHandlerEnabledFilters(t *testing.T) {
	t.Parallel()

	// Все тесты пройдут, если отладочные логи включены с помощью примененных фильтров.
	const testLevel = Debug

	tests := []struct {
		name     string
		root     Level
		filters  map[string]Level
		groups   []string
		expected bool
	}{
		{
			name:     "debug/no filters/root",
			root:     Debug,
			filters:  nil,
			groups:   nil,
			expected: true,
		},
		{
			name:     "debug/subroot info/root",
			root:     Debug,
			filters:  map[string]Level{"subroot": Info},
			groups:   nil,
			expected: true,
		},
		{
			name:     "debug/subroot info/subroot",
			root:     Debug,
			filters:  map[string]Level{"subroot": Info},
			groups:   []string{"subroot"},
			expected: false,
		},
		{
			name:     "debug/subroot info/subsubroot",
			root:     Debug,
			filters:  map[string]Level{"subroot": Info},
			groups:   []string{"subroot", "subsubroot"},
			expected: false,
		},
		{
			name:     "info/no filters/root",
			root:     Info,
			filters:  nil,
			groups:   nil,
			expected: false,
		},
		{
			name:     "info/subroot debug/subroot",
			root:     Info,
			filters:  map[string]Level{"subroot": Debug},
			groups:   []string{"subroot"},
			expected: true,
		},
		{
			name:     "info/subroot debug/subsubroot",
			root:     Info,
			filters:  map[string]Level{"subroot": Debug},
			groups:   []string{"subroot", "subsubroot"},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			h := zlTestHandler(t, nil)

			h.zl = h.zl.Level(zlevel(test.root))
			h.filter = filter.NewBackwardFilter[zlGroup](test.filters, keySeparator)

			for _, group := range test.groups {
				var ok bool
				h, ok = h.WithGroup(group).(zlHandler)
				require.True(t, ok)
			}

			enabled := h.Enabled(context.Background(), testLevel)
			require.Equal(t, test.expected, enabled)
		})
	}
}

func TestHandlerHandle(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	const (
		key      = "message"
		expected = "handle"
	)

	var tw testWriter
	h := zlTestHandler(t, &tw)

	err := h.Handle(context.Background(), zlTestRecord(t, expected))
	require.NoError(err)

	var rec map[string]interface{}
	err = json.Unmarshal(tw.Bytes(), &rec)
	require.NoError(err)

	val, ok := rec[key]
	require.True(ok)
	require.Equal(expected, val)
}

func TestHandlerContextHook(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	const (
		key      = "key"
		expected = "value"
	)

	kvhook := func(_ context.Context, stash logutil.Stash) logutil.Stash {
		return logutil.WithAttr(stash, key, expected)
	}

	var tw testWriter
	h := zlTestHandler(t, &tw)
	h.hook = append(h.hook, kvhook)

	err := h.Handle(context.Background(), zlTestRecord(t, "hook"))
	require.NoError(err)

	var rec map[string]interface{}
	err = json.Unmarshal(tw.Bytes(), &rec)
	require.NoError(err)

	val, ok := rec[key]
	require.True(ok)
	require.Equal(expected, val)
}

func TestHandlerHandleSingleGroup(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	const group = "group"
	expected := map[string]any{
		"key": "value",
	}

	var tw testWriter
	var h slog.Handler = zlTestHandler(t, &tw)

	h = h.WithGroup(group)

	err := h.Handle(context.Background(), zlTestRecord(t, "group", slog.String("key", "value")))
	require.NoError(err)

	var rec map[string]interface{}
	err = json.Unmarshal(tw.Bytes(), &rec)
	require.NoError(err)

	val, ok := rec[group]
	require.True(ok)
	require.Equal(expected, val)
}

func TestHandlerHandleNestedGroups(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	groups := []string{"top", "mid", "bottom"}
	expected := map[string]any{
		"top": map[string]any{
			"mid": map[string]any{
				"bottom": map[string]any{
					"key": "value",
				},
			},
		},
	}

	var tw testWriter
	var h slog.Handler = zlTestHandler(t, &tw)

	for _, group := range groups {
		h = h.WithGroup(group)
	}

	err := h.Handle(context.Background(), zlTestRecord(t, "groups", slog.String("key", "value")))
	require.NoError(err)

	// testWriter использует JSON, преобразуем результат и сравним с ожидаемым значением.
	var rec map[string]interface{}
	err = json.Unmarshal(tw.Bytes(), &rec)
	require.NoError(err)

	// разворачиваем и проверяем вложенность групп.
	for _, group := range groups {
		var ok bool
		rec, ok = rec[group].(map[string]any)
		require.True(ok)

		expected, ok = expected[group].(map[string]any)
		require.True(ok)

		require.Equal(expected, rec)
	}
}

func TestHandlerHandleCallerAttrs(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	const (
		key      = "key"
		expected = "value"
	)

	var tw testWriter
	h := zlTestHandler(t, &tw)

	rec := zlTestRecord(t, "attrs")
	rec.AddAttrs(slog.String(key, expected))

	err := h.Handle(context.Background(), rec)
	require.NoError(err)

	var result map[string]interface{}
	err = json.Unmarshal(tw.Bytes(), &result)
	require.NoError(err)

	val, ok := result[key]
	require.True(ok)
	require.Equal(expected, val)
}

func TestWithAttrs(t *testing.T) {
	t.Parallel()

	var h slog.Handler = zlTestHandler(t, nil)

	attrs := []slog.Attr{
		slog.String("key", "value"),
	}

	newh := h.WithAttrs(attrs)
	require.NotEqual(t, h, newh)
}

func TestWithGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		values   []string
		expected []zlGroup
	}{
		{"none", []string{}, nil},
		{"empty", []string{""}, nil},
		{"single", []string{"one"}, []zlGroup{{name: "one"}}},
		{"multi", []string{"one", "two"}, []zlGroup{{name: "two"}, {name: "one"}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var h slog.Handler = zlTestHandler(t, nil)

			for _, group := range test.values {
				h = h.WithGroup(group)
			}

			require.Equal(t, test.expected, h.(zlHandler).groups) //nolint: errcheck
		})
	}
}

func TestZlevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		slevel slog.Level
		zlevel zerolog.Level
	}{
		{"subdebug", slog.LevelDebug - 1, zerolog.DebugLevel},
		{"debug", slog.LevelDebug, zerolog.DebugLevel},
		{"supdebug", slog.LevelDebug + 1, zerolog.DebugLevel},
		{"subinfo", slog.LevelInfo - 1, zerolog.DebugLevel},
		{"info", slog.LevelInfo, zerolog.InfoLevel},
		{"supinfo", slog.LevelInfo + 1, zerolog.InfoLevel},
		{"subwarn", slog.LevelWarn - 1, zerolog.InfoLevel},
		{"warn", slog.LevelWarn, zerolog.WarnLevel},
		{"supwarn", slog.LevelWarn + 1, zerolog.WarnLevel},
		{"suberror", slog.LevelError - 1, zerolog.WarnLevel},
		{"error", slog.LevelError, zerolog.ErrorLevel},
		{"superror", slog.LevelError + 1, zerolog.ErrorLevel},
		{"subdisabled", disabled - 1, zerolog.ErrorLevel},
		{"disabled", disabled, zerolog.Disabled},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			level := zlevel(test.slevel)
			require.Equal(t, test.zlevel, level)
		})
	}
}

func TestZLAppend(t *testing.T) {
	tests := []struct {
		name  string
		value slog.Value
	}{
		{"group", slog.GroupValue(slog.Bool("true", true))},
		{"bool", slog.BoolValue(true)},
		{"float64", slog.Float64Value(1.0)},
		{"int64", slog.Int64Value(-1)},
		{"uint64", slog.Uint64Value(1)},
		{"string", slog.StringValue("value")},
		{"duration", slog.DurationValue(1 * time.Second)},
		{"time", slog.TimeValue(time.Now())},
		{"error", slog.AnyValue(errors.New("error"))},
	}

	t.Run("zerolog.Context", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				ctx := zerolog.Nop().With()

				ctx = zlAppend(ctx, slog.Attr{Key: test.name, Value: test.value})
				require.NotNil(t, ctx)
			})
		}
	})

	t.Run("zerolog.Event", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				evt := zerolog.Dict()

				evt = zlAppend(evt, slog.Attr{Key: test.name, Value: test.value})
				require.NotNil(t, evt)
			})
		}
	})
}

func TestZlDict(t *testing.T) {
	t.Parallel()

	attrs := []slog.Attr{
		slog.String("key", "value"),
	}

	dict := zlEventAttrs(zerolog.Dict(), attrs)
	require.NotNil(t, dict)
}

func TestZlWriterConsole(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	w := fileWriter{
		DestConfig: DestConfig{
			Format: FormatConsole,
		},
	}

	iow, err := zlWriter(w)
	require.NoError(err)
	require.NotNil(iow)
	require.IsType(zerolog.ConsoleWriter{}, iow)
}

func TestZlWriterJson(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	w := fileWriter{
		DestConfig: DestConfig{
			Format: FormatJSON,
		},
	}

	iow, err := zlWriter(w)
	require.NoError(err)
	require.NotNil(iow)
	require.Equal(w, iow)
}

func TestZlWriterError(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	tests := []struct {
		format string
	}{
		{""},
		{"nosj"},
	}

	for _, test := range tests {
		t.Run(test.format, func(t *testing.T) {
			t.Parallel()

			w := fileWriter{
				DestConfig: DestConfig{
					Format: Format(test.format),
				},
			}

			iow, err := zlWriter(w)
			require.Error(err)
			require.Nil(iow)
		})
	}
}

func zlTestHandler(t *testing.T, w *testWriter) zlHandler {
	t.Helper()

	var opt options
	opt.Level = Info
	opt.writers = append(opt.writers, w)

	h, err := newLogHandler(opt)
	require.NoError(t, err)

	return h
}

func zlTestRecord(t *testing.T, message string, attrs ...slog.Attr) slog.Record {
	t.Helper()

	pc, _, _, ok := runtime.Caller(1)
	require.True(t, ok)

	rec := slog.NewRecord(time.Now(), slog.LevelInfo, message, pc)
	if len(attrs) > 0 {
		rec.AddAttrs(attrs...)
	}

	return rec
}

type testWriter struct {
	bytes.Buffer
}

func (*testWriter) Config() DestConfig {
	return DestConfig{Format: FormatJSON}
}
