package todoreminder

import (
	"context"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
	"go.uber.org/zap"
)

type handler struct {
	*core_telegram.BaseHandler
	service *service

	mtx    sync.RWMutex
	states map[int64]*taskDialogState
}

func newHandler(logger *core_logger.Logger, service *service) *handler {
	return &handler{
		BaseHandler: core_telegram.NewBaseHandler(logger),
		service:     service,
		states:      map[int64]*taskDialogState{},
	}
}

func (h *handler) startDialog(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	h.setState(userID, &taskDialogState{
		Step: choiceButton,
		Task: newTaskUninitialized(),
	})

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "Выберите что Вас интересует",
		ReplyMarkup: menuKeyboard(),
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
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Ошибка, начните заново",
		})
		if err != nil {
			h.Log.Warn("failed to send message", zap.Int64("user_id", chatID), zap.Error(err))
		}
		return
	}

	switch state.Step {
	case choiceButton:
		h.handleChoiceButton(ctx, b, update, state, chatID, userID, msgID)
	}
}

func (h *handler) handleTextMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	state, ok := h.getState(userID)
	if !ok {
		return
	}

	switch state.Step {
	case inputTitle:
		h.handleTitle(ctx, b, state, chatID, userID, text)
	case inputDescription:
		h.handleDescription(ctx, b, state, chatID, userID, text)
	}
}

func (h *handler) handle() []core_telegram.Handle {
	return []core_telegram.Handle{
		{
			HandleType: bot.HandlerTypeMessageText,
			Pattern:    "todo",
			MatchType:  bot.MatchTypeCommand,
			Function:   h.startDialog,
		},
		{
			HandleType: bot.HandlerTypeCallbackQueryData,
			Pattern:    "task:",
			MatchType:  bot.MatchTypePrefix,
			Function:   h.handleCallback,
		},
		{
			HandleType: bot.HandlerTypeMessageText,
			Pattern:    "",
			MatchType:  bot.MatchTypePrefix,
			Function:   h.handleTextMessage,
		},
	}
}
