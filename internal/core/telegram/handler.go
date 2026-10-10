package core_telegram

import (
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
)

type BaseHandler struct {
	Log *core_logger.Logger
}

func NewBaseHandler(log *core_logger.Logger) *BaseHandler {
	return &BaseHandler{
		Log: log,
	}
}
