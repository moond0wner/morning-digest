package core_telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_telegram_middleware "github.com/moond0wner/morning-digest/internal/core/telegram/middleware"
	"go.uber.org/zap"
)

type Bot struct {
	bot    *bot.Bot
	logger *core_logger.Logger
}

func NewBot(logger *core_logger.Logger) (*Bot, error) {
	config, err := newConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load bot config: %w", err)
	}

	safeErrors := func(err error) {
		safe := err.Error()
		if config.botToken != "" {
			safe = strings.ReplaceAll(err.Error(), config.botToken, "******")
		}
		logger.Error("some error", zap.String("error", safe))
	}

	opts := []bot.Option{
		bot.WithErrorsHandler(safeErrors),
		bot.WithWorkers(10),
	}

	b, err := bot.New(config.botToken, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	return &Bot{
		bot:    b,
		logger: logger,
	}, nil

}

func (b *Bot) Start(ctx context.Context) {
	b.logger.Debug("Starting bot")
	b.bot.Start(ctx)
}

func (b *Bot) Register(handles ...Handle) {
	for _, handle := range handles {

		// фича-middleware + хендлер
		handler := core_telegram_middleware.ChainMiddleware(
			handle.Function,
			handle.Middleware...,
		)

		// глобальные снаружи фича-middleware
		handler = core_telegram_middleware.ChainMiddleware(
			handler,
			core_telegram_middleware.TagHandler(handle.Pattern),
			core_telegram_middleware.RequestID(),
			core_telegram_middleware.Logger(b.logger),
			core_telegram_middleware.Trace(),
			core_telegram_middleware.Panic(),
		)

		b.bot.RegisterHandler(
			handle.HandleType,
			handle.Pattern,
			handle.MatchType,
			handler,
		)

		// TagHandler -> RequestID -> Logger -> Trace -> Panic -> фича-middleware -> handle.Function
	}
}
