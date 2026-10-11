package core_telegram

import (
	"github.com/go-telegram/bot"
	core_telegram_middleware "github.com/moond0wner/morning-digest/internal/core/telegram/middleware"
)

type Handle struct {
	HandleType bot.HandlerType
	Pattern    string
	MatchType  bot.MatchType
	Function   bot.HandlerFunc
	Middleware []core_telegram_middleware.Middleware
}
