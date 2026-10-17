package habr

import (
	"context"
	"strings"

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
		BaseFeature: core_feature.NewBaseFeature("habr", logger, nil),
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
	service := newService(f.Log, f.retrier)
	f.handler = newHandler(f.Log, service)

	_, err := f.handler.service.baseParse(ctx, 1)
	return err
}

func (f *Feature) Brief(ctx context.Context) (string, error) {
	if f.handler == nil {
		return "", nil
	}

	articles, err := f.handler.service.baseParse(ctx, 3)
	if err != nil {
		return "", err
	}

	if len(articles) == 0 {
		return "", nil
	}

	var sb strings.Builder
	sb.WriteString("📰 Хабр:\n")
	for _, a := range articles {
		sb.WriteString("• " + a.Title + "\n")
	}
	return strings.TrimSpace(sb.String()), nil
}
