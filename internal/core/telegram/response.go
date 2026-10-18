package core_telegram

import (
	"context"
	"errors"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
	"go.uber.org/zap"
)

func (b *BaseHandler) ReplyError(
	ctx context.Context,
	bt *bot.Bot,
	chatID int64,
	userMsg string,
	err error,
) {
	var e *core_errors.Error
	if errors.As(err, &e) {
		switch e.Type {
		case core_errors.ErrUserInput:
			b.Log.Warn(e.Type.String(), zap.Error(err))
		case core_errors.ErrNotFound:
			b.Log.Warn(e.Type.String(), zap.Error(err))
		case core_errors.ErrExternalAPI:
			b.Log.Warn(e.Type.String(), zap.Error(err))
		case core_errors.ErrInternal:
			b.Log.Error(e.Type.String(), zap.Error(err))
		default:
			b.Log.Error(e.Type.String(), zap.Error(err))
		}
	} else {
		b.Log.Error("unknown error", zap.Error(err))
	}

	if userMsg == "" {
		userMsg = "Не удалось обработать запрос"
	}
	_, err = bt.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   userMsg,
	})
	if err != nil {
		b.Log.Warn("failed to send message", zap.Error(err))
	}
}

func (b *BaseHandler) EditMessage(
	ctx context.Context,
	bt *bot.Bot,
	chatID int64,
	msgID int,
	text string,
	kb ...*models.InlineKeyboardMarkup,
) {
	params := &bot.EditMessageTextParams{
		ChatID:    chatID,
		MessageID: msgID,
		Text:      text,
	}
	if len(kb) > 0 && kb[0] != nil {
		params.ReplyMarkup = kb[0]
	}

	_, err := bt.EditMessageText(ctx, params)
	if err != nil {
		b.Log.Warn("failed to edit message",
			zap.Int64("chat_id", chatID),
			zap.Int("msg_id", msgID),
			zap.Error(err),
		)
	}
}
