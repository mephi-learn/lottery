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
	"strconv"
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

// EnvEnrichment обогащает конфигурацию данными из переменных окружения.
func EnvEnrichment(config *Config) *Config {
	config.Storage.Postgres.Host = getEnv("DB_HOST", config.Storage.Postgres.Host)
	if port, err := strconv.Atoi(getEnv("DB_PORT", strconv.Itoa(config.Storage.Postgres.Port))); err == nil {
		config.Storage.Postgres.Port = port
	}
	config.Storage.Postgres.User = getEnv("DB_USER", config.Storage.Postgres.User)
	config.Storage.Postgres.Password = getEnv("DB_PASSWORD", config.Storage.Postgres.Password)
	config.Storage.Postgres.Database = getEnv("DB_NAME", config.Storage.Postgres.Database)
	config.Storage.Postgres.Schema = getEnv("DB_SCHEMA", config.Storage.Postgres.Schema)
	config.Storage.Postgres.SSLMode = getEnv("DB_SSL_MODE", config.Storage.Postgres.SSLMode)

	return config
}

// getEnv получает значение переменной окружения.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
