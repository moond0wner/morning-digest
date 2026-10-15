package habr

import (
	"context"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
	"go.uber.org/zap"
)

type Handler struct {
	*core_telegram.BaseHandler
	service *Service

	mtx    sync.RWMutex
	states map[int64]*NewsDialogState
}

func NewHandler(log *core_logger.Logger, srvc *Service) *Handler {
	return &Handler{
		BaseHandler: core_telegram.NewBaseHandler(log),
		service:     srvc,
		states:      make(map[int64]*NewsDialogState),
	}
}

func (h *Handler) StartDialog(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	h.mtx.Lock()
	h.states[userID] = &NewsDialogState{
		Step: ChoiceSource,
	}
	h.mtx.Unlock()

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "Что парсим?",
		ReplyMarkup: sourceKeyboard(),
	})
	if err != nil {
		h.Log.Warn("failed to send message", zap.Int64("user_id", chatID), zap.Error(err))
	}
}

func (h *Handler) HandleCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
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
	case ChoiceSource:
		h.handleChoiceSource(ctx, b, update, state, chatID, msgID)
	case ChoiceHub:
		h.handleChoiceHub(ctx, b, update, state, chatID, msgID)
	case ChoiceLimit:
		h.handleChoiceLimit(ctx, b, update, state, chatID, userID)
	}
}

func (h *Handler) Handle() []core_telegram.Handle {
	return []core_telegram.Handle{
		{
			HandleType: bot.HandlerTypeMessageText,
			Pattern:    "news",
			MatchType:  bot.MatchTypeCommand,
			Function:   h.StartDialog,
		},
		{
			HandleType: bot.HandlerTypeCallbackQueryData,
			Pattern:    "news:",
			MatchType:  bot.MatchTypePrefix,
			Function:   h.HandleCallback,
		},
	}
}
