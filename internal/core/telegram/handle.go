package core_telegram

import "github.com/go-telegram/bot"

type Handle struct {
	HandleType bot.HandlerType
	Pattern    string
	MatchType  bot.MatchType
	Function   bot.HandlerFunc
	Middleware []bot.Middleware
}
