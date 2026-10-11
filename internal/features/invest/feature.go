package invest

import (
	"context"
	"time"

	core_feature "github.com/moond0wner/morning-digest/internal/core/feature"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
)

type Feature struct {
	core_feature.BaseFeature
	handler *Handler
}

func NewFeature(logger *core_logger.Logger, srvc *Service) *Feature {
	h := NewHandler(logger, srvc)
	return &Feature{
		BaseFeature: core_feature.NewBaseFeature("invest", h.Handle()),
		handler:     h,
	}
}

func (f *Feature) Init(ctx context.Context) error {
	return core_feature.WithTimeout(ctx, 10*time.Second, func(ctx context.Context) error {
		_, err := f.handler.service.GetPortfolio(ctx)
		return err
	})
}
