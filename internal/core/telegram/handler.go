package core_telegram

import (
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
)

type BaseHandler struct {
	Log *core_logger.Logger
}

func NewBaseHandler(logger *core_logger.Logger) *BaseHandler {
	return &BaseHandler{
		Log: logger,
	}
}
