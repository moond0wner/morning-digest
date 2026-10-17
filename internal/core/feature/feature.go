package core_feature

import (
	"context"

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
