package todo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
	core_repository "github.com/moond0wner/morning-digest/internal/core/repository"
)

type repository struct {
	*core_repository.BaseRepository
}

func newRepository(base *core_repository.BaseRepository) *repository {
	return &repository{BaseRepository: base}
}

func (r *repository) createTable(ctx context.Context) error {
	sqlQuery := `
    CREATE TABLE IF NOT EXISTS tasks (
        id           INTEGER 	PRIMARY KEY AUTOINCREMENT,
        title        TEXT 		NOT NULL,
        description  TEXT,
        completed    INTEGER 	NOT NULL DEFAULT 0,
        created_at   TEXT 		NOT NULL,
        completed_at TEXT
    );`
	_, err := r.Exec(ctx, sqlQuery)
	return err
}

func (r *repository) createTask(ctx context.Context, task *Task) error {
	sqlQuery := `
	INSERT INTO tasks (title, description, completed, created_at) 
	VALUES (?, ?, 0, ?);
	`
	_, err := r.Exec(
		ctx,
		sqlQuery,
		task.Title,
		task.Description,
		task.CreatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"exec error",
		)
	}
	return nil
}

func (r *repository) getTasks(ctx context.Context) ([]Task, error) {
	sqlQuery := `
	SELECT id, title, description, completed, created_at, completed_at 
	FROM tasks;
	`

	var tasks []Task
	err := r.Query(ctx, sqlQuery, nil, func(rows *sql.Rows) error {
		for rows.Next() {
			var (
				completedAtNull sql.NullTime
				createdAtStr    string
				t               Task
				err             error
			)
			if err = rows.Scan(
				&t.ID,
				&t.Title,
				&t.Description,
				&t.Completed,
				&createdAtStr,
				&completedAtNull,
			); err != nil {
				return err
			}
			if t.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr); err != nil {
				return err
			}
			if completedAtNull.Valid {
				t.CompletedAt = &completedAtNull.Time
			}
			tasks = append(tasks, t)
		}
		return rows.Err()
	})

	if err != nil {
		return nil, core_errors.Wrap(core_errors.ErrInternal, err, "query tasks")
	}
	return tasks, nil
}

func (r *repository) getTask(ctx context.Context, taskID int64) (Task, error) {
	sqlQuery := `
    SELECT id, title, description, completed, created_at, completed_at
    FROM tasks
    WHERE id = ?;
    `

	var (
		t               Task
		createdAtStr    string
		completedAtNull sql.NullTime
		scanErr         error
	)

	err := r.QueryRow(ctx, sqlQuery, []any{taskID}, func(row *sql.Row) error {
		scanErr = row.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Completed,
			&createdAtStr,
			&completedAtNull,
		)
		return scanErr
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, core_errors.Wrap(core_errors.ErrNotFound, err, "task not found")
		}
		return Task{}, core_errors.Wrap(core_errors.ErrInternal, err, "scan error")
	}

	if t.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr); err != nil {
		return Task{}, core_errors.Wrap(core_errors.ErrInternal, err, "parse created_at")
	}
	if completedAtNull.Valid {
		t.CompletedAt = &completedAtNull.Time
	}

	return t, nil
}

func (r *repository) updateTask(ctx context.Context, task Task) error {
	sqlQuery := `
	UPDATE tasks
	SET 
		title = ?, 
		description = ?, 
		completed = ?, 
		completed_at = ?
	WHERE id = ?;
	`

	_, err := r.Exec(
		ctx,
		sqlQuery,
		task.Title,
		task.Description,
		task.Completed,
		task.CompletedAt,
		task.ID,
	)

	if err != nil {
		return core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"database exec",
		)
	}

	return nil
}
