package echo

import (
	core_feature "github.com/moond0wner/morning-digest/internal/core/feature"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
)

type Feature struct {
	*core_feature.BaseFeature
	handler *handler
}

func NewFeature(logger *core_logger.Logger) *Feature {
	h := newHandler(logger)
	return &Feature{
		BaseFeature: core_feature.NewBaseFeature("echo", logger, h.Handle()),
		handler:     h,
	}
}

// init не нужен
