package log

import (
	"log/slog"
	"math"
	"os"

	"homework/pkg/errors"
	"homework/pkg/log/logutil"
)

const keySeparator = '/'

// Level это псевдоним для [slog.Level].
type Level = slog.Level

const (
	Debug = slog.LevelDebug
	Info  = slog.LevelInfo
	Error = slog.LevelError

	disabled Level = math.MaxInt
)

// Logger это псевдоним для *[slog.Logger].
type Logger = *slog.Logger

// New создает новый экземпляр логера.
//
// Если не настроены места вывода, то будет установлен вывод в консоль.
//
// Если не установлен параметр логирования [LoggerConfig.Level], то [Info] будет уровнем по умолчанию.
func New(config LoggerConfig, logopts ...LoggerOption) (Logger, error) {
	var hopts options
	hopts.LoggerConfig = config

	if level := config.Level; level == "" {
		config.Level = Info.String()
	}

	err := hopts.Level.UnmarshalText([]byte(config.Level))
	if err != nil {
		return nil, errors.Errorf("invalid level '%s'", config.Level)
	}

	if !levelValid(hopts.Level) {
		return nil, errors.Errorf("unsupported level '%s'", hopts.Level)
	}

	for key, level := range hopts.Filters {
		if !levelValid(level) {
			return nil, errors.Errorf("unsupported filter level '%s': '%s'", level, key)
		}
	}

	if cfg := config.Stdout; cfg != nil {
		hopts.writers, err = appendFileWriter(hopts.writers, os.Stdout, *cfg)
		if err != nil {
			return nil, errors.Errorf("creating stdout writer: %w", err)
		}
	}

	if cfg := config.File; cfg != nil {
		hopts.writers, err = appendPathWriter(hopts.writers, cfg.Path, cfg.DestConfig)
		if err != nil {
			return nil, errors.Errorf("creating file writer: %w", err)
		}
	}

	if err = hopts.apply(logopts); err != nil {
		return nil, errors.Errorf("applying logger options: %w", err)
	}

	// ПРИМЕЧАНИЕ: не должно быть никакой логики, зависящей от конкретной реализации,
	// поэтому мы можем переключить реализацию newLogHandler с помощью простых тегов сборки(build tags).
	h, err := newLogHandler(hopts)
	if err != nil {
		return nil, err
	}

	return slog.New(h), nil
}

// WithEventHook конфигурирует [EventHook] для событий логера.
func WithEventHook(hook logutil.EventHook) LoggerOption {
	return func(o *options) error {
		o.hook = append(o.hook, hook)
		return nil
	}
}
