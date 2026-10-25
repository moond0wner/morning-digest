package todo

import (
	"context"
	"testing"
)

type mockRepository struct {
	tasks           map[int64]Task
	updateWasCalled bool
}

func (m *mockRepository) createTask(ctx context.Context, task *Task) error { return nil }
func (m *mockRepository) getTasks(ctx context.Context) ([]Task, error)     { return nil, nil }
func (m *mockRepository) getTask(ctx context.Context, taskID int64) (Task, error) {
	return m.tasks[taskID], nil
}

func (m *mockRepository) updateTask(ctx context.Context, task Task) error {
	m.updateWasCalled = true
	return nil
}

func TestService_UpdateTask(t *testing.T) {
	t.Run("успешное выполнение задачи", func(t *testing.T) {
		mock := &mockRepository{
			tasks: map[int64]Task{
				1: {ID: 1, Title: "Купить хлеб", Completed: false},
			},
		}
		svc := newService(mock)

		err := svc.updateTask(context.Background(), 1)
		if err != nil {
			t.Fatalf("updateTask вернул ошибку: %v", err)
		}
		if !mock.updateWasCalled {
			t.Error("ожидали, что updateTask в репозитории вызовется, но он не вызвался")
		}
	})

	t.Run("задача уже выполнена", func(t *testing.T) {
		mock := &mockRepository{
			tasks: map[int64]Task{
				1: {ID: 1, Title: "Купить хлеб", Completed: true},
			},
		}
		svc := newService(mock)

		err := svc.updateTask(context.Background(), 1)
		if err != nil {
			t.Fatalf("updateTask вернул ошибку: %v", err)
		}
		if mock.updateWasCalled {
			t.Error("ожидали, что updateTask НЕ вызовется для выполненной задачи, но он вызвался")
		}
	})
}
