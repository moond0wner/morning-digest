package echo

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
	"go.uber.org/zap"
)

type handler struct {
	*core_telegram.BaseHandler
}

func newHandler(logger *core_logger.Logger) *handler {
	return &handler{
		BaseHandler: core_telegram.NewBaseHandler(logger),
	}
}

func (h *handler) Echo(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	text := "Я вас не понимаю."
	chatID := update.Message.Chat.ID

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})

	if err != nil {
		h.Log.Warn("failed to send message", zap.Int64("user_id", chatID), zap.Error(err))
	}
}

func (h *handler) Handle() []core_telegram.Handle {
	return []core_telegram.Handle{
		{
			HandleType: bot.HandlerTypeMessageText,
			Pattern:    "",
			MatchType:  bot.MatchTypePrefix,
			Function:   h.Echo,
		},
	}
}
