package log

import (
	"bytes"

	"homework/pkg/log/logutil"
)

type LoggerConfig struct {
	Level   string           `yaml:"level"   json:"level"`
	Filters map[string]Level `yaml:"filters" json:"filters"`

	Stdout *DestConfig `yaml:"stdout" json:"stdout"`
	File   *FileConfig `yaml:"file"   json:"file"`
}

type levelConfig Level

func (l levelConfig) String() string {
	if Level(l) == disabled {
		return "DISABLED"
	}

	return Level(l).String()
}

func (l *levelConfig) UnmarshalText(text []byte) error {
	text = bytes.ToUpper(text)

	if bytes.Equal(text, []byte("DISABLED")) {
		*l = levelConfig(disabled)
		return nil
	}

	return (*Level)(l).UnmarshalText(text)
}

func levelValid(l Level) bool {
	switch l {
	case Debug, Info, Error, disabled:
		return true
	default:
		return false
	}
}

// FileConfig позволяет настроить дополнительный log файл.
type FileConfig struct {
	DestConfig
	Path string `yaml:"path" json:"path"`
}

// LoggerOption это тип для функциональных опций, которые будут использоваться с [New].
type LoggerOption func(*options) error

type options struct {
	LoggerConfig

	Level   Level
	hook    logutil.Hook
	writers []writer
}

func (o *options) apply(opts []LoggerOption) error {
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return err
		}
	}

	return nil
}
