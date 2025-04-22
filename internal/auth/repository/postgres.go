package repository

import (
	"context"
	"database/sql"
	"homework/internal/models"
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

func (r *repository) Create(ctx context.Context, userData *models.SignUpInput) (int, error) {
	if userData.Admin {
		var userCount int
		if err := r.db.QueryRowContext(ctx, "select count(*) cn from users where admin=$1", true).Scan(&userCount); err != nil {
			return -1, errors.Errorf("failed to check exist user: %w", err)
		}
		if userCount > 0 {
			return -1, errors.Errorf("only one user can be admin")
		}
	}

	var userId int
	if err := r.db.QueryRowContext(ctx, "insert into users(name, username, password, email, admin) values($1, $2, $3, $4, $5) returning id",
		userData.Name, userData.Username, userData.Password, userData.Email, userData.Admin).Scan(&userId); err != nil {
		return -1, errors.Errorf("failed to create user: %w", err)
	}

	return userId, nil
}

func (r *repository) GetByUsernameAndPassword(ctx context.Context, userData *models.SignInInput) (*models.User, error) {
	user := models.User{}
	if err := r.db.QueryRowContext(ctx, "select id, name, username, email, admin from users where username=$1 and password=$2", userData.Username, userData.Password).
		Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Admin); err != nil {
		return nil, errors.Errorf("user not found: %w", err)
	}

	return &user, nil
}

func (r *repository) GetById(ctx context.Context, userId int) (*models.User, error) {
	user := models.User{}
	if err := r.db.QueryRowContext(ctx, "select id, name, username, email, admin from users where id=$1", userId).
		Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Admin); err != nil {
		return nil, errors.Errorf("user not found: %w", err)
	}

	return &user, nil
}
