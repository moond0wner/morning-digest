package weather

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
	"go.uber.org/zap"
)

type Handler struct {
	*core_telegram.BaseHandler
	service *Service
}

func NewHandler(log *core_logger.Logger, srvc *Service) *Handler {
	return &Handler{
		BaseHandler: core_telegram.NewBaseHandler(log),
		service:     srvc,
	}
}

func (h *Handler) Weather(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	weather, err := h.service.Get(ctx, "Saratov")
	if err != nil {
		h.Log.Error("get weather", zap.Error(err))
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Не удалось получить погоду",
		})

		if err != nil {
			h.Log.Warn("failed to send message", zap.Int64("user_id", chatID), zap.Error(err))
		}
		return
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   weather.String(),
	})

	if err != nil {
		h.Log.Warn("failed to send message", zap.Int64("user_id", chatID), zap.Error(err))
	}

}

func (h *Handler) Handle() core_telegram.Handle {
	return core_telegram.Handle{
		HandleType: bot.HandlerTypeMessageText,
		Pattern:    "weather",
		MatchType:  bot.MatchTypeCommand,
		Function:   h.Weather,
	}
}
