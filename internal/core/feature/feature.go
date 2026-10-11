package core_feature

import (
	"context"
	"fmt"
	"time"

	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
)

type Feature interface {
	Name() string
	Init(ctx context.Context) error
	Handle() []core_telegram.Handle
}

type BaseFeature struct {
	name    string
	handles []core_telegram.Handle
}

func NewBaseFeature(name string, handles []core_telegram.Handle) BaseFeature {
	return BaseFeature{name: name, handles: handles}
}

func (b *BaseFeature) Name() string                   { return b.name }
func (b *BaseFeature) Handle() []core_telegram.Handle { return b.handles }
func (b *BaseFeature) Init(ctx context.Context) error { return nil } // проверяет работоспособность фичи, в каждой фиче эта ф-я переопределяется

// WithTimeout запускает функцию в горутине с таймаутом.
// Возвращает ошибку если функция не успела или вернула ошибку.
func WithTimeout(ctx context.Context, d time.Duration, f func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- f(ctx)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("timeout: %w", ctx.Err())
	case err := <-done:
		return err
	}
}
