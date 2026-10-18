package digest

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
)

type handler struct {
	*core_telegram.BaseHandler
	service *service
}

func newHandler(logger *core_logger.Logger, service *service) *handler {
	return &handler{
		BaseHandler: core_telegram.NewBaseHandler(logger),
		service:     service,
	}
}

func (h *handler) handleDigest(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	text, err := h.service.build(ctx)
	if err != nil {
		h.ReplyError(ctx, b, chatID, "Не удалось собрать дайджест", err)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
}

func (h *handler) sendDigestOnUserID(ctx context.Context, b *bot.Bot, digest string, accountID int64) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: accountID,
		Text:   digest,
	})
}

func (h *handler) handle() []core_telegram.Handle {
	return []core_telegram.Handle{
		{
			HandleType: bot.HandlerTypeMessageText,
			Pattern:    "digest",
			MatchType:  bot.MatchTypeCommand,
			Function:   h.handleDigest,
		},
	}
}
