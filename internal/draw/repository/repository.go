package repository

import (
	"context"
	"database/sql"
	"homework/pkg/errors"
	"homework/pkg/log"
)

type Storage interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Conn(ctx context.Context) (*sql.Conn, error)
}

// RepoOption позволяет настроить репозиторий добавлением новых функциональных опций.
type RepoOption func(*repository) error

type repository struct {
	db Storage

	log log.Logger
}

// NewRepository создаёт объект репозитория, который должен удовлетворять требованиям сервисов.
func NewRepository(opts ...RepoOption) (*repository, error) {
	var repo repository

	for _, opt := range opts {
		if err := opt(&repo); err != nil {
			return nil, errors.Errorf("apply option: %w", err)
		}
	}

	if repo.log == nil {
		return nil, errors.Errorf("no logger provided")
	}

	if repo.db == nil {
		return nil, errors.Errorf("no database provided")
	}

	return &repo, nil
}

func WithLogger(logger log.Logger) RepoOption {
	return func(r *repository) error {
		r.log = logger
		return nil
	}
}

func WithStorage(st Storage) RepoOption {
	return func(r *repository) error {
		r.db = st
		return nil
	}
}
