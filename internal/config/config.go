package config

import (
	"homework/internal/config/parsers"
	"homework/internal/server"
	"homework/internal/storage"
	"homework/pkg/errors"
	"homework/pkg/log"
	"io"
	"os"
	"path/filepath"
)

var ErrUnknownFormat = errors.New("unknown file format")

type Config struct {
	Server struct {
		HTTP server.Config `yaml:"http" json:"http"`
	} `yaml:"server" json:"server"`

	Storage struct {
		Postgres storage.Config `yaml:"postgres" json:"postgres"`
	} `yaml:"storage" json:"storage"`

	Logger log.LoggerConfig `yaml:"logger" json:"logger"`
}

// NewConfig создаёт экземпляр конфигурации.
func NewConfig() *Config {
	return &Config{}
}

// NewConfigFromFile создаёт экземпляр конфигурации из файла конфигурации.
func NewConfigFromFile(filename string) (*Config, error) {
	var config Config
	var parser parsers.Parsers[Config]

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	file, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.Errorf("read %s: %w", f.Name(), err)
	}

	switch filepath.Ext(f.Name()) {
	case ".yaml", ".yml":
		parser = &parsers.YamlParser[Config]{}
	case ".json":
		parser = &parsers.JsonParser[Config]{}
	default:
		return nil, errors.Errorf("parse %s: %w", f.Name(), ErrUnknownFormat)
	}

	if err = parser.Parse(string(file), &config); err != nil {
		return nil, errors.Errorf("parse %s: %w", f.Name(), err)
	}

	return &config, nil
}
