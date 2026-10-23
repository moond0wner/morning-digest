package start

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

func newHandler(log *core_logger.Logger) *handler {
	return &handler{
		BaseHandler: core_telegram.NewBaseHandler(log),
	}
}

func (h *handler) start(ctx context.Context, b *bot.Bot, update *models.Update) {
	// TODO: Переделать текст
	text := "Привет, я твой бот-помощник. Пока что я ничего не умею, но все впереди!"
	chatID := update.Message.Chat.ID

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})

	if err != nil {
		h.Log.Warn("failed to send message", zap.Int64("user_id", chatID), zap.Error(err))
	}
}

func (h *handler) handle() []core_telegram.Handle {
	return []core_telegram.Handle{
		{
			HandleType: bot.HandlerTypeMessageText,
			Pattern:    "start",
			MatchType:  bot.MatchTypeCommand,
			Function:   h.start,
		},
	}
}
