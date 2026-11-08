package core_telegram

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	core_telegram_middleware "github.com/moond0wner/morning-digest/internal/core/telegram/middleware"
)

type Handle struct {
	Middleware []core_telegram_middleware.Middleware
	MatchFunc  func(*models.Update) bool
	Pattern    string
	HandleType bot.HandlerType
	MatchType  bot.MatchType
	Function   bot.HandlerFunc
}
