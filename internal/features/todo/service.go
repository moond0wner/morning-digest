package todo

import (
	"context"
	"fmt"
	"time"
)

type taskRepository interface {
	createTask(ctx context.Context, task *Task) error
	getTasks(ctx context.Context) ([]Task, error)
	getTask(ctx context.Context, taskID int64) (Task, error)
	updateTask(ctx context.Context, task Task) error
}

type service struct {
	repository taskRepository
}

func newService(repository taskRepository) *service {
	return &service{
		repository: repository,
	}
}

func (s *service) createTask(ctx context.Context, task *Task) error {
	if err := s.repository.createTask(ctx, task); err != nil {
		return fmt.Errorf("failed to create task :%w", err)
	}

	return nil
}

func (s *service) getTasks(ctx context.Context) ([]Task, error) {
	tasks, err := s.repository.getTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	return tasks, nil
}

func (s *service) updateTask(ctx context.Context, taskID int64) error {
	task, err := s.repository.getTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	if task.Completed {
		return nil
	}

	task.Completed = true
	completedAt := time.Now().UTC()
	task.CompletedAt = &completedAt

	if err := s.repository.updateTask(ctx, task); err != nil {
		return fmt.Errorf("failed to update tasks: %w", err)
	}
	return nil
}
