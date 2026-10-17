package habr

import (
	"context"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
	"go.uber.org/zap"
)

type handler struct {
	*core_telegram.BaseHandler
	service *service

	mtx    sync.RWMutex
	states map[int64]*newsDialogState
}

func newHandler(log *core_logger.Logger, srvc *service) *handler {
	return &handler{
		BaseHandler: core_telegram.NewBaseHandler(log),
		service:     srvc,
		states:      make(map[int64]*newsDialogState),
	}
}

func (h *handler) startDialog(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	h.setState(userID, &newsDialogState{
		Step: choiceSource,
	})

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "Что парсим?",
		ReplyMarkup: sourceKeyboard(),
	})
	if err != nil {
		h.Log.Warn(
			"failed to send message",
			zap.Int64("chat_id", chatID),
			zap.Int64("user_id", userID),
			zap.Error(err),
		)
	}
}

func (h *handler) handleCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	userID := update.CallbackQuery.From.ID
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	msgID := update.CallbackQuery.Message.Message.ID

	state, ok := h.getState(userID)

	if !ok {
		h.ReplyError(ctx, b, chatID, "Ошибка, начните заново", core_errors.New(
			core_errors.ErrUserInput,
			"state error",
		))
		return
	}

	switch state.Step {
	case choiceSource:
		h.handleChoiceSource(ctx, b, update, state, chatID, msgID)
	case choiceHub:
		h.handleChoiceHub(ctx, b, update, state, chatID, msgID)
	case choiceLimit:
		h.handleChoiceLimit(ctx, b, update, state, chatID, userID)
	}
}

func (h *handler) handle() []core_telegram.Handle {
	return []core_telegram.Handle{
		{
			HandleType: bot.HandlerTypeMessageText,
			Pattern:    "news",
			MatchType:  bot.MatchTypeCommand,
			Function:   h.startDialog,
		},
		{
			HandleType: bot.HandlerTypeCallbackQueryData,
			Pattern:    "news:",
			MatchType:  bot.MatchTypePrefix,
			Function:   h.handleCallback,
		},
	}
}
