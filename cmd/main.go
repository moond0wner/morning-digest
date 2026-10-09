package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_config "github.com/moond0wner/morning-digest/internal/core/config"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	"go.uber.org/zap"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	log *core_logger.Logger
}

func NewHandler(log *core_logger.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Printf("Failed to init application logger: %v", err)
		os.Exit(1)
	}

	handlers := NewHandler(logger)

	config := core_config.NewConfigMust()

	logger.Debug("application config")

	b, err := bot.New(config.BotToken)
	if err != nil {
		logger.Error("failed to create bot", zap.Error(err))
		os.Exit(1)
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handlers.Start)
	logger.Debug("create bot, register handlers and start")

	b.Start(ctx)
}

func (h *Handler) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := "Привет, я твой бот-помощник. Пока что я ничего не умею, но все впереди!"
	chatID := update.Message.Chat.ID
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		h.log.Warn("failed to send message", zap.Any("user_id", chatID), zap.Error(err))
	}
}
