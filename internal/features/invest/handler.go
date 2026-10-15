package invest

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
	service *service
}

func newHandler(log *core_logger.Logger, srvc *service) *handler {
	return &handler{
		BaseHandler: core_telegram.NewBaseHandler(log),
		service:     srvc,
	}
}

func (h *handler) invest(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	portfolio, err := h.service.getPortfolio()
	if err != nil {
		h.ReplyError(ctx, b, chatID, "Не удалось собрать информацию о Вашем счете", err)
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   portfolio.String(),
	})

	if err != nil {
		h.Log.Warn("failed to send message", zap.Int64("user_id", chatID), zap.Error(err))
	}
}

func (h *handler) handle() []core_telegram.Handle {
	return []core_telegram.Handle{
		{
			HandleType: bot.HandlerTypeMessageText,
			Pattern:    "invest",
			MatchType:  bot.MatchTypeCommand,
			Function:   h.invest,
		},
	}
}
