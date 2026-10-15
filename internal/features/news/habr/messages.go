package habr

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func (h *handler) sendArticles(
	ctx context.Context,
	b *bot.Bot,
	state *newsDialogState,
	chatID int64,
	userID int64,
) {
	var articles []Article
	var err error
	switch state.SourceType {
	case "base":
		articles, err = h.service.baseParse(ctx, state.Limit)
		if err != nil {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Не удалось получить статьи",
			})
			h.Log.Warn("failed to base parse",
				zap.Int64("user_id", userID),
				zap.Int64("chat_id", chatID),
				zap.Error(err),
			)
			if err != nil {
				h.Log.Warn("failed to send message",
					zap.Int64("user_id", userID),
					zap.Int64("chat_id", chatID),
					zap.Error(err),
				)
			}
			return
		}
	case "hub":
		articles, err = h.service.hubParse(ctx, state.Hub, state.Limit)
		if err != nil {
			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Не удалось получить статьи",
			})
			h.Log.Warn("failed to hub parse",
				zap.Int64("user_id", userID),
				zap.Error(err),
			)
			if err != nil {
				h.Log.Warn("failed to send message",
					zap.Int64("user_id", userID),
					zap.Int64("chat_id", chatID),
					zap.Error(err),
				)
			}
			return
		}
	}

	if len(articles) == 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Статей не найдено",
		})
		return
	}

	var sb strings.Builder
	for _, article := range articles {
		fmt.Fprintln(&sb, article.String())
	}
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   sb.String(),
	})

	if err != nil {
		h.Log.Warn("failed to send message",
			zap.Int64("user_id", userID),
			zap.Int64("chat_id", chatID),
			zap.Error(err),
		)
	}
}

func (h *handler) editMessage(
	ctx context.Context,
	b *bot.Bot,
	chatID int64,
	msgID int,
	text string,
	kb *models.InlineKeyboardMarkup,
) {
	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   msgID,
		Text:        text,
		ReplyMarkup: kb,
	})
	if err != nil {
		h.Log.Warn("failed to edit message",
			zap.Int64("chat_id", chatID),
			zap.Int("msg_id", msgID),
			zap.Error(err),
		)
	}
}
