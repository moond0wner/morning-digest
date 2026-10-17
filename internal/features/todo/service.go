package todo

import (
	"context"
	"fmt"
	"time"

	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
)

type service struct {
	repository *repository
	log        *core_logger.Logger
}

func newService(repository *repository, logger *core_logger.Logger) *service {
	return &service{
		repository: repository,
		log:        logger,
	}
}

func (s *service) close() error {
	return s.repository.close()
}

func (s *service) createTask(ctx context.Context, task *Task) error {
	err := s.repository.createTask(ctx, task)
	if err != nil {
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

	task.Completed = true
	completedAt := time.Now()
	task.CompletedAt = &completedAt

	if err := s.repository.updateTask(ctx, task); err != nil {
		return fmt.Errorf("failed to update tasks: %w", err)
	}
	return nil
}
