package repository

import (
	"context"
	"database/sql"
	"homework/internal/models"
	"homework/pkg/errors"
	"homework/pkg/log"
	"time"
)

type Storage interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
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

func (r *repository) Create(ctx context.Context, begin time.Time, start time.Time, lotteryType string) (int, error) {
	return 0, nil
}

func (r *repository) Cancel(ctx context.Context, drawId int) error {
	return nil
}

func (r *repository) SetBeginTime(ctx context.Context, drawId int, begin time.Time) error {
	return nil
}

func (r *repository) SetStartTime(ctx context.Context, drawId int, start time.Time) error {
	return nil
}

func (r *repository) ListActive(ctx context.Context) ([]models.Draw, error) {
	return nil, nil
}
