package core_repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
	_ "modernc.org/sqlite"
)

type BaseRepository struct {
	db        *sql.DB
	opTimeout time.Duration
}

func NewBaseRepository(dbName string) (*BaseRepository, error) {
	if err := os.MkdirAll("./out/data", 0755); err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"mkdir database dir",
		)
	}

	db, err := sql.Open("sqlite", fmt.Sprintf("./out/data/%s.db", dbName))
	if err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"create/open database",
		)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"ping database",
		)
	}

	db.SetMaxOpenConns(1)

	return &BaseRepository{
		db:        db,
		opTimeout: 10 * time.Second,
	}, nil
}

func (r *BaseRepository) Close() error {
	return r.db.Close()
}
func (r *BaseRepository) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(ctx, r.opTimeout)
	defer cancel()
	return r.db.ExecContext(ctx, query, args...)
}

func (r *BaseRepository) Query(
	ctx context.Context,
	query string,
	args []any,
	scan func(*sql.Rows) error,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.opTimeout)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return scan(rows)
}

func (r *BaseRepository) QueryRow(
	ctx context.Context,
	query string,
	args []any,
	scan func(*sql.Row) error,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.opTimeout)
	defer cancel()

	row := r.db.QueryRowContext(ctx, query, args...)
	return scan(row)
}
