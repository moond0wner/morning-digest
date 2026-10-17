package weather

import (
	"context"

	core_feature "github.com/moond0wner/morning-digest/internal/core/feature"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_retry "github.com/moond0wner/morning-digest/internal/core/retry"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
)

type Feature struct {
	core_feature.BaseFeature
	handler *handler
	retrier *core_retry.Retrier
}

func NewFeature(logger *core_logger.Logger, retrier *core_retry.Retrier) *Feature {
	return &Feature{
		BaseFeature: core_feature.NewBaseFeature("weather", logger, nil),
		retrier:     retrier,
	}
}

func (f *Feature) Handle() []core_telegram.Handle {
	if f.handler == nil {
		return nil
	}
	return f.handler.handle()
}

func (f *Feature) Init(ctx context.Context) error {
	service, err := newService(f.Log, f.retrier)
	if err != nil {
		return err
	}

	f.handler = newHandler(f.Log, service)

	_, err = f.handler.service.client.getCoordinates(ctx, "Москва")
	return err
}
