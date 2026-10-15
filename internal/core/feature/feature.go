package core_feature

import (
	"context"
	"fmt"
	"time"

	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
)

type Feature interface {
	Name() string
	Init(ctx context.Context) error
	Handle() []core_telegram.Handle
	Close(ctx context.Context) error
}

type BaseFeature struct {
	name    string
	handles []core_telegram.Handle
	Log     *core_logger.Logger
}

func NewBaseFeature(name string, logger *core_logger.Logger, handles []core_telegram.Handle) BaseFeature {
	return BaseFeature{name: name, Log: logger, handles: handles}
}

func (b *BaseFeature) Name() string                    { return b.name }
func (b *BaseFeature) Handle() []core_telegram.Handle  { return b.handles }
func (b *BaseFeature) Init(ctx context.Context) error  { return nil } // проверяет работоспособность фичи, в каждой фиче эта ф-я переопределяется
func (b *BaseFeature) Close(ctx context.Context) error { return nil }

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
