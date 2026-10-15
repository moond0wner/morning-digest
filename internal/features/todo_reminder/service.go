package todoreminder

import (
	"context"
	"fmt"

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

func (s *service) createTask(ctx context.Context, task Task) error {
	// if err := task.validate(); err != nil {
	// 	return Task{}, fmt.Errorf("validate task domain: %w", err)
	// }

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

// func (s *service) getTask(ctx context.Context, taskID int64) (Task, error) {
// 	tasks, err := s.repository.getTask(ctx, taskID)
// 	if err != nil {
// 		return Task{}, fmt.Errorf("failed to get tasks: %w", err)
// 	}
// 	return tasks, nil
// }
