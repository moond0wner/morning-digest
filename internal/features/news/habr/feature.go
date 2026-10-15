package habr

import (
	"context"
	"time"

	core_feature "github.com/moond0wner/morning-digest/internal/core/feature"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
)

type Feature struct {
	core_feature.BaseFeature
	handler *handler
}

func NewFeature(logger *core_logger.Logger) *Feature {
	return &Feature{
		BaseFeature: core_feature.NewBaseFeature("habr", logger, nil),
	}
}

func (f *Feature) Handle() []core_telegram.Handle {
	if f.handler == nil {
		return nil
	}
	return f.handler.handle()
}

func (f *Feature) Init(ctx context.Context) error {
	service := newService(f.Log)
	f.handler = newHandler(f.Log, service)

	return core_feature.WithTimeout(ctx, 30*time.Second, func(ctx context.Context) error {
		_, err := f.handler.service.baseParse(ctx, 1)
		return err
	})
}
