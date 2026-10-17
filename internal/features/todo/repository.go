package todo

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
	_ "modernc.org/sqlite"
)

type repository struct {
	db        *sql.DB
	opTimeout time.Duration
}

func newRepository() (*repository, error) {
	if err := os.MkdirAll("./out/data", 0755); err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"mkdir database dir",
		)
	}

	db, err := sql.Open("sqlite", "./out/data/todo.db")
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

	return &repository{
		db:        db,
		opTimeout: 10 * time.Second,
	}, nil
}

func (r *repository) close() error {
	return r.db.Close()
}

func (r *repository) createTable() error {
	sqlQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		id 				INTEGER 	 	PRIMARY KEY 	AUTOINCREMENT,
		title 			TEXT			NOT NULL,
		description 	TEXT,
		completed 		BOOLEAN			NOT NULL,
		created_at 		TIMESTAMP 		NOT NULL,
		completed_at 	TIMESTAMP
	);
	`
	if _, err := r.db.Exec(sqlQuery); err != nil {
		return core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"created table in database",
		)
	}
	return nil
}

func (r *repository) createTask(ctx context.Context, task *Task) error {
	ctx, cancel := context.WithTimeout(ctx, r.opTimeout)
	defer cancel()

	sqlQuery := `
	INSERT INTO tasks (title, description, completed, created_at)
	VALUES (?, ?, ?, ?);
	`

	result, err := r.db.Exec(
		sqlQuery,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
	)
	if err != nil {
		return core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"database exec",
		)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"get id from db",
		)
	}
	task.ID = id

	return nil
}

func (r *repository) getTasks(ctx context.Context) ([]Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.opTimeout)
	defer cancel()

	sqlQuery := `
	SELECT id, title, description, completed, created_at, completed_at
	FROM tasks
	ORDER BY id ASC;
	`

	rows, err := r.db.Query(sqlQuery)
	if err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"get tasks",
		)
	}

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Completed,
			&t.CreatedAt,
			&t.CompletedAt,
		); err != nil {
			return nil, core_errors.Wrap(
				core_errors.ErrInternal,
				err,
				"scan error",
			)
		}
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"rows error",
		)
	}

	return tasks, nil
}

func (r *repository) getTask(ctx context.Context, taskID int64) (Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.opTimeout)
	defer cancel()

	sqlQuery := `
	SELECT id, title, description, completed, created_at, completed_at
	FROM tasks
	WHERE id=$1;
	`

	row := r.db.QueryRow(sqlQuery, taskID)

	var task Task
	if err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.CompletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, core_errors.Wrap(
				core_errors.ErrNotFound,
				err,
				"task not found",
			)
		}
		return Task{}, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"scan error",
		)
	}
	return task, nil
}

func (r *repository) updateTask(ctx context.Context, task Task) error {
	ctx, cancel := context.WithTimeout(ctx, r.opTimeout)
	defer cancel()

	sqlQuery := `
	UPDATE tasks
	SET
		title=$1,
		description=$2,
		completed=$3,
		completed_at=$4
	WHERE id=$5;
	`

	_, err := r.db.Exec(
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
