package core_telegram

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	core_config "github.com/moond0wner/morning-digest/internal/core/config"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
)

type Bot struct {
	bot    *bot.Bot
	logger *core_logger.Logger
	config core_config.Config
}

func NewBot(logger *core_logger.Logger, config core_config.Config) (*Bot, error) {
	b, err := bot.New(config.BotToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	return &Bot{
		bot:    b,
		logger: logger,
		config: config,
	}, nil

}

func (b *Bot) Start(ctx context.Context) {
	b.logger.Debug("Starting bot")
	b.bot.Start(ctx)
}
func (b *Bot) Register(handles ...Handle) {
	for _, handle := range handles {
		b.bot.RegisterHandler(
			handle.HandleType,
			handle.Pattern,
			handle.MatchType,
			handle.Function,
			handle.Middleware...,
		)
	}
}
