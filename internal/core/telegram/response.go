package core_telegram

import (
	"context"
	"errors"

	"github.com/go-telegram/bot"
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
			userMsg = "Проверьте ввод"
		case core_errors.ErrNotFound:
			b.Log.Warn(e.Type.String(), zap.Error(err))
			userMsg = "Не найдено"
		case core_errors.ErrExternalAPI:
			b.Log.Warn(e.Type.String(), zap.Error(err))
			userMsg = "Сервис временно недоступен"
		case core_errors.ErrInternal:
			b.Log.Error(e.Type.String(), zap.Error(err))
		default:
			b.Log.Error(e.Type.String(), zap.Error(err))
		}
	} else {
		b.Log.Error("unknown error", zap.Error(err))
	}

	_, err = bt.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   userMsg,
	})
	if err != nil {
		b.Log.Warn("failed to send message", zap.Error(err))
	}
}
