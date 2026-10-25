package todo

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
	core_postgres_pool "github.com/moond0wner/morning-digest/internal/core/repository/postgres/pool"
	core_pgx_pool "github.com/moond0wner/morning-digest/internal/core/repository/postgres/pool/pgx"
)

type repository struct {
	pool *core_pgx_pool.Pool
}

func newRepository(pool *core_pgx_pool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) createTask(ctx context.Context, task *Task) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
	INSERT INTO tasks (title, description, completed, created_at, completed_at) 
	VALUES ($1, $2, $3, $4, $5);
	`

	_, err := r.pool.Exec(
		ctx,
		sqlQuery,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
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
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
	SELECT id, title, description, completed, created_at, completed_at 
	FROM tasks
	ORDER BY id ASC;
	`

	rows, err := r.pool.Query(ctx, sqlQuery)
	if err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"query error",
		)
	}

	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.CompletedAt,
		); err != nil {
			return nil, core_errors.Wrap(
				core_errors.ErrInternal,
				err,
				"scan error",
			)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"next rows",
		)
	}

	return tasks, nil
}

func (r *repository) getTask(ctx context.Context, taskID int64) (Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
    SELECT id, title, description, completed, created_at, completed_at
    FROM tasks
    WHERE id=$1;
    `

	row := r.pool.QueryRow(ctx, sqlQuery, taskID)

	var task Task
	if err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.CompletedAt,
	); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return Task{}, core_errors.Wrap(
				core_errors.ErrNotFound,
				err,
				fmt.Sprintf("task with id='%d' not found", taskID),
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
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
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

	_, err := r.pool.Exec(
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
			"exec error",
		)
	}

	return nil
}
