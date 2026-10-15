package habr

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *handler) handleChoiceSource(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *newsDialogState,
	chatID int64,
	msgID int,
) {
	data := update.CallbackQuery.Data
	switch data {
	case cbSourceBase:
		state.SourceType = "base"
		state.Step = ChoiceLimit
		h.editMessage(ctx, b, chatID, msgID, "Выберите количество статей", limitKeyboard())
	case cbSourceHub:
		state.SourceType = "hub"
		state.Step = ChoiceHub
		h.editMessage(ctx, b, chatID, msgID, "Выберите хаб", hubKeyboard())
	}
}

func (h *handler) handleChoiceHub(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *newsDialogState,
	chatID int64,
	msgID int,
) {
	data := update.CallbackQuery.Data
	hub := strings.TrimPrefix(data, cbHubPrefix)
	state.Hub = hub
	state.Step = ChoiceLimit
	h.editMessage(ctx, b, chatID, msgID, "Выберите количество статей", limitKeyboard())
}

func (h *handler) handleChoiceLimit(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *newsDialogState,
	chatID int64,
	userID int64,
) {
	data := update.CallbackQuery.Data
	limitStr := strings.TrimPrefix(data, cbLimitPrefix)
	limit, _ := strconv.Atoi(limitStr)
	state.Limit = limit

	defer h.deleteState(userID)
	h.sendArticles(ctx, b, state, chatID, userID)
}
