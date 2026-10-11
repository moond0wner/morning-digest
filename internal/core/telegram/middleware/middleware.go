package core_telegram_middleware

import "github.com/go-telegram/bot"

type Middleware func(bot.HandlerFunc) bot.HandlerFunc

func ChainMiddleware(
	b bot.HandlerFunc,
	m ...Middleware,
) bot.HandlerFunc {
	if len(m) == 0 {
		return b
	}

	for i := len(m) - 1; i >= 0; i-- {
		b = m[i](b)
	}

	return b
}
