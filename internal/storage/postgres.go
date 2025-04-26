package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type Storage interface {
	PingContext(ctx context.Context) error
	Ping() error
	Close() error
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	SetConnMaxLifetime(d time.Duration)
	SetConnMaxIdleTime(d time.Duration)
	Stats() sql.DBStats
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Prepare(query string) (*sql.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryRow(query string, args ...any) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Begin() (*sql.Tx, error)
	Conn(ctx context.Context) (*sql.Conn, error)
}

func (r *storage) PingContext(ctx context.Context) error {
	return r.postgres.PingContext(ctx)
}

func (r *storage) Ping() error {
	return r.postgres.Ping()
}

func (r *storage) Close() error {
	return r.postgres.Close()
}

func (r *storage) SetMaxIdleConns(n int) {
	r.postgres.SetMaxIdleConns(n)
}

func (r *storage) SetMaxOpenConns(n int) {
	r.postgres.SetMaxOpenConns(n)
}

func (r *storage) SetConnMaxLifetime(d time.Duration) {
	r.postgres.SetConnMaxLifetime(d)
}

func (r *storage) SetConnMaxIdleTime(d time.Duration) {
	r.postgres.SetConnMaxIdleTime(d)
}

func (r *storage) Stats() sql.DBStats {
	return r.postgres.Stats()
}

func (r *storage) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return r.postgres.PrepareContext(ctx, query)
}

func (r *storage) Prepare(query string) (*sql.Stmt, error) {
	return r.postgres.Prepare(query)
}

func (r *storage) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return r.postgres.ExecContext(ctx, query, args...)
}

func (r *storage) Exec(query string, args ...any) (sql.Result, error) {
	return r.postgres.Exec(query, args...)
}

func (r *storage) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return r.postgres.QueryContext(ctx, query, args...)
}

func (r *storage) Query(query string, args ...any) (*sql.Rows, error) {
	return r.postgres.Query(query, args...)
}

func (r *storage) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return r.postgres.QueryRowContext(ctx, query, args...)
}

func (r *storage) QueryRow(query string, args ...any) *sql.Row {
	return r.postgres.QueryRow(query, args...)
}

func (r *storage) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return r.postgres.BeginTx(ctx, opts)
}

func (r *storage) Begin() (*sql.Tx, error) {
	return r.postgres.Begin()
}

func (r *storage) Conn(ctx context.Context) (*sql.Conn, error) {
	return r.postgres.Conn(ctx)
}
