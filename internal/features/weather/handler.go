package weather

import (
	"context"
	"strconv"
	"strings"
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
	cities map[int64][]geoResponse
}

func newHandler(log *core_logger.Logger, srvc *service) *handler {
	return &handler{
		BaseHandler: core_telegram.NewBaseHandler(log),
		service:     srvc,
		cities:      make(map[int64][]geoResponse),
	}
}

func (h *handler) weather(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	args := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/weather"))
	if args == "" {
		h.ReplyError(ctx, b, chatID, "Укажите город: /weather Москва", core_errors.New(
			core_errors.ErrUserInput,
			"user input error",
		))
		return
	}

	candidates, err := h.service.GetCandidates(ctx, args)
	if err != nil {
		h.ReplyError(ctx, b, chatID, "Не удалось получить погоду", err)
		return
	}

	if len(candidates) == 0 {
		h.ReplyError(ctx, b, chatID, "Город не найден", core_errors.New(
			core_errors.ErrNotFound,
			"city not found",
		))
		return
	}

	if len(candidates) == 1 {
		w, err := h.service.GetByCoords(ctx, candidates[0].Lat, candidates[0].Lon)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Не удалось получить погоду",
			})
			h.Log.Warn("failed to get coords", zap.Error(err))
			return
		}
		b.SendMessage(ctx,
			&bot.SendMessageParams{
				ChatID: chatID,
				Text:   w.String(),
			})
		return
	}

	h.mtx.Lock()
	h.cities[userID] = candidates
	h.mtx.Unlock()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        "Выберите город:",
		ReplyMarkup: geoKeyboard(candidates),
	})

}

func (h *handler) handleCityCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	h.handleGetCity(ctx, b, update)

}

func (h *handler) handleGetCity(ctx context.Context, b *bot.Bot, update *models.Update) {
	data := update.CallbackQuery.Data
	userID := update.CallbackQuery.From.ID
	chatID := update.CallbackQuery.Message.Message.Chat.ID
	msgID := update.CallbackQuery.Message.Message.ID

	indexCity, err := strconv.Atoi(strings.TrimPrefix(data, "weather:city:"))
	if err != nil {
		h.ReplyError(ctx, b, chatID, "Не удалось получить погоду", err)
		return
	}

	h.mtx.RLock()
	candidates, ok := h.cities[userID]
	h.mtx.RUnlock()

	if !ok || indexCity < 0 || indexCity >= len(candidates) {
		h.ReplyError(ctx, b, chatID, "Город не найден", core_errors.New(
			core_errors.ErrNotFound,
			"city not found"))
		return
	}

	chosen := candidates[indexCity]
	w, err := h.service.GetByCoords(ctx, chosen.Lat, chosen.Lon)
	if err != nil {
		h.ReplyError(ctx, b, chatID, "Не удалось получить погоду", err)
		return
	}

	h.EditMessage(ctx, b, chatID, msgID, w.String(), core_telegram.DeleteKeyboard())

	h.mtx.Lock()
	delete(h.cities, userID)
	h.mtx.Unlock()
}

func (h *handler) handle() []core_telegram.Handle {
	return []core_telegram.Handle{
		{
			HandleType: bot.HandlerTypeMessageText,
			Pattern:    "weather",
			MatchType:  bot.MatchTypeCommand,
			Function:   h.weather,
		},
		{
			HandleType: bot.HandlerTypeCallbackQueryData,
			Pattern:    "weather:city:",
			MatchType:  bot.MatchTypePrefix,
			Function:   h.handleCityCallback,
		},
	}
}
