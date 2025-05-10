package storage

import (
	"database/sql"
	"fmt"
	"homework/pkg/log"

	"github.com/go-errors/errors"
	_ "github.com/lib/pq"
)

var _ Storage = (*storage)(nil)

type Config struct {
	Host     string `yaml:"host"     json:"host"`
	Port     int    `yaml:"port"     json:"port"`
	User     string `yaml:"user"     json:"user"`
	Password string `yaml:"password" json:"password"`
	Database string `yaml:"database" json:"database"`
	SSLMode  string `yaml:"ssl_mode" json:"ssl_mode"`
	Schema   string `yaml:"schema"   json:"schema"`
}

// Option позволяет настроить репозиторий добавлением новых функциональных опций.
type Option func(*storage) error

type storage struct {
	postgres struct {
		*Config
		*sql.DB
	}

	log log.Logger
}

// NewStorage создаёт объект репозитория, который должен удовлетворять требованиям сервисов.
func NewStorage(opts ...Option) (*storage, error) {
	var st storage

	for _, opt := range opts {
		if err := opt(&st); err != nil {
			return nil, errors.Errorf("apply option: %w", err)
		}
	}

	if st.postgres.Config == nil {
		return nil, errors.Errorf("no config")
	}

	// st.ValidateConfig()

	if st.log == nil {
		return nil, errors.Errorf("no logger provided")
	}

	dsn := connectString(st.postgres.Config)

	var err error
	st.postgres.DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Errorf("unable create storage: %w", err)
	}

	err = st.postgres.DB.Ping()
	if err != nil {
		return nil, errors.Errorf("unable create session: %w", err)
	}

	return &st, nil
}

func connectString(cfg *Config) string {
	schema := ""
	if cfg.Schema != "" {
		schema = "search_path=" + cfg.Schema + ",public"
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, schema)
}

func BuildDSN(cfg *Config) string {
	schema := "public"
	if cfg.Schema != "" {
		schema = "search_path=" + cfg.Schema + ",public"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?search_path=%s&sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, schema, cfg.SSLMode)
}

func WithLogger(logger log.Logger) Option {
	return func(r *storage) error {
		r.log = logger
		return nil
	}
}

func WithConfig(config Config) Option {
	return func(r *storage) error {
		r.postgres.Config = &config
		return nil
	}
}
