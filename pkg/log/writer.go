package log

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"homework/pkg/errors"
)

// Format тип для формата вывода.
type Format string

const (
	formatNone    Format = ""
	FormatConsole Format = "console"
	FormatJSON    Format = "json"
)

// DestConfig определяет общую конфигурацию вывода.
type DestConfig struct {
	Format Format `yaml:"format" json:"format"` // "console" или "json"
}

// appendFileWriter конфигурирует [writer] для [os.Stdout].
// Если формат не установлен, по умолчанию будет [FormatConsole].
func appendFileWriter(writers []writer, f *os.File, config DestConfig) ([]writer, error) {
	if config.Format == formatNone {
		config.Format = FormatConsole
	}

	exist := slices.ContainsFunc(writers, func(w writer) bool {
		fw, ok := w.(fileWriter)
		return ok && fw.equals(f)
	})

	if exist {
		return nil, errors.Errorf("duplicate file writer for '%s'", f.Name())
	}

	fw := fileWriter{File: f, DestConfig: config}

	return append(writers, fw), nil
}

// appendPathWriter конфигурирует [writer] для файла, используя его путь.
// Файл создастся, если не существует.
// Если формат не установлен, по умолчанию будет [FormatJSON] для файла с расширением json.
func appendPathWriter(writers []writer, path string, config DestConfig) ([]writer, error) {
	if config.Format == formatNone {
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".json" {
			config.Format = FormatJSON
		}
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	return appendFileWriter(writers, f, config)
}

// writer это [io.Writer] с оберткой для вложенной конфигурации.
type writer interface {
	io.Writer

	Config() DestConfig
}

type fileWriter struct {
	*os.File
	DestConfig
}

func (f fileWriter) Config() DestConfig {
	return f.DestConfig
}

func (f fileWriter) equals(other *os.File) bool {
	fi, ef := f.Stat()
	oi, eo := other.Stat()

	return ef == nil && eo == nil && os.SameFile(fi, oi)
}
