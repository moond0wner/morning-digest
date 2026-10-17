package core_telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_retry "github.com/moond0wner/morning-digest/internal/core/retry"
	core_telegram_middleware "github.com/moond0wner/morning-digest/internal/core/telegram/middleware"
	"go.uber.org/zap"
)

type Bot struct {
	bot     *bot.Bot
	log     *core_logger.Logger
	retrier *core_retry.Retrier
}

func NewBot(ctx context.Context, logger *core_logger.Logger, retrier *core_retry.Retrier) (*Bot, error) {
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

	var b *bot.Bot
	if err := retrier.Do(ctx, func(ctx context.Context) error {
		var err error
		b, err = bot.New(config.botToken, opts...)
		return err
	}); err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	return &Bot{
		bot:     b,
		log:     logger,
		retrier: retrier,
	}, nil

}

func (b *Bot) Start(ctx context.Context) {
	b.log.Debug("Starting bot")
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
			core_telegram_middleware.Logger(b.log),
			core_telegram_middleware.Trace(),
			core_telegram_middleware.Panic(),
		)

		if handle.MatchFunc != nil {
			b.bot.RegisterHandlerMatchFunc(handle.MatchFunc, handler)
		} else {
			b.bot.RegisterHandler(
				handle.HandleType,
				handle.Pattern,
				handle.MatchType,
				handler,
			)
		}

		// TagHandler -> RequestID -> Logger -> Trace -> Panic -> фича-middleware -> handle.Function
	}
}
