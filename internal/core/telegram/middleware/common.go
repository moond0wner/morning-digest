package core_telegram_middleware

import (
	"context"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	"go.uber.org/zap"
)

type ctxKey string

const (
	handlerTagKey ctxKey = "handler_tag"
	RequestIDKey  ctxKey = "request_id"
	KeyLength            = 8
)

func TagHandler(name string) Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return bot.HandlerFunc(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			ctx = context.WithValue(ctx, handlerTagKey, name)
			next(ctx, b, update)
		})
	}
}

func GetHandlerTag(ctx context.Context) string {
	if name, ok := ctx.Value(handlerTagKey).(string); ok {
		return name
	}
	return "unknown"
}

func RequestID() Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return bot.HandlerFunc(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			requestID := bot.RandomString(KeyLength)
			ctx = context.WithValue(ctx, RequestIDKey, requestID)
			next(ctx, b, update)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return bot.HandlerFunc(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			requestID := ctx.Value(RequestIDKey)

			fields := []zap.Field{
				zap.Any("request_id", requestID),
			}

			if update.Message != nil {
				fields = append(fields,
					zap.Int64("user_id", update.Message.Chat.ID),
					zap.String("text", update.Message.Text),
				)
			}

			l := log.With(fields...)
			ctx = core_logger.ToContext(ctx, l)
			next(ctx, b, update)
		})
	}
}

func Trace() Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return bot.HandlerFunc(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			log := core_logger.FromContext(ctx)

			before := time.Now()
			log.Debug(
				">>> incoming handler request",
				zap.String("handler", GetHandlerTag(ctx)),
			)

			next(ctx, b, update)

			log.Debug(
				"<<< done handler request",
				zap.Duration("latency", time.Since(before)),
			)
		})
	}
}

func Panic() Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return bot.HandlerFunc(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			log := core_logger.FromContext(ctx)
			defer func() {
				if p := recover(); p != nil {
					log.Error("during handler request got unexpected panic",
						zap.Any("panic", p),
						zap.Stack("stack"),
					)

					var chatID int64
					switch {
					case update.Message != nil:
						chatID = update.Message.Chat.ID
					case update.CallbackQuery != nil && update.CallbackQuery.Message.Message != nil:
						chatID = update.CallbackQuery.Message.Message.Chat.ID
					default:
						return
					}

					if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: chatID,
						Text:   "При обработке запроса возникла непредвиденная ошибка",
					}); err != nil {
						log.Error("failed to send panic message", zap.Error(err))
					}
				}
			}()

			next(ctx, b, update)
		})
	}
}
