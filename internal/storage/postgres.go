package storage

import (
	"database/sql"
	"fmt"
	"github.com/go-errors/errors"
	_ "github.com/lib/pq"
	"homework/pkg/log"
)

var _ Storage = (*storage)(nil)

type Config struct {
	Host     string `yaml:"host"     json:"host"`
	Port     int    `yaml:"port"     json:"port"`
	User     string `yaml:"user"     json:"user"`
	Password string `yaml:"password" json:"password"`
	Database string `yaml:"database" json:"database"`
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

	//st.ValidateConfig()

	if st.log == nil {
		return nil, errors.Errorf("no logger provided")
	}

	var err error
	config := st.postgres.Config
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Database)
	st.postgres.DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, errors.Errorf("unable create storage: %w", err)
	}

	err = st.postgres.DB.Ping()
	if err != nil {
		return nil, errors.Errorf("unable create session: %w", err)
	}

	return &st, nil
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
