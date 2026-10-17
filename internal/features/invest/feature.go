package invest

import (
	"context"
	"fmt"

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
		BaseFeature: core_feature.NewBaseFeature("invest", logger, nil),
		retrier:     retrier,
	}
}

func (f *Feature) Init(ctx context.Context) error {
	service, err := newService(ctx, f.Log, f.retrier)
	if err != nil {
		return err
	}

	f.handler = newHandler(f.Log, service)

	_, err = f.handler.service.getPortfolio()
	return err
}

func (f *Feature) Handle() []core_telegram.Handle {
	if f.handler == nil {
		return nil
	}
	return f.handler.handle()
}

func (f *Feature) Close(ctx context.Context) error {
	if f.handler == nil || f.handler.service == nil {
		return nil
	}
	return f.handler.service.close()
}

func (f *Feature) Brief(ctx context.Context) (string, error) {
	if f.handler == nil {
		return "", nil
	}

	p, err := f.handler.service.getPortfolio()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("💼 Портфель: %.0f ₽ (за день %+.0f ₽)", p.TotalAmount, p.DailyYield), nil
}
