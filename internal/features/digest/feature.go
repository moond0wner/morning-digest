package digest

import (
	"context"

	core_feature "github.com/moond0wner/morning-digest/internal/core/feature"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
)

type Feature struct {
	core_feature.BaseFeature
	handler *handler
	service *service
}

func NewFeature(log *core_logger.Logger, briefers []Briefer) *Feature {
	return &Feature{
		BaseFeature: core_feature.NewBaseFeature("digest", log, nil),
		service:     newService(briefers, log),
	}
}

func (f *Feature) Init(ctx context.Context) error {
	f.handler = newHandler(f.Log, f.service)
	return nil
}

func (f *Feature) Handle() []core_telegram.Handle {
	if f.handler == nil {
		return nil
	}
	return f.handler.handle()
}
