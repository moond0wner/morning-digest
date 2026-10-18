package core_telegram

import (
	"context"

	"github.com/go-telegram/bot"
)

type MessageSender struct {
	bot *bot.Bot
}

func NewMessageSender(b *bot.Bot) *MessageSender {
	return &MessageSender{bot: b}
}

func (s *MessageSender) SendMessage(ctx context.Context, chatID int64, text string) error {
	_, err := s.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	return err
}
