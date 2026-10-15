package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_feature "github.com/moond0wner/morning-digest/internal/core/feature"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_retry "github.com/moond0wner/morning-digest/internal/core/retry"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
	"github.com/moond0wner/morning-digest/internal/features/echo"
	"github.com/moond0wner/morning-digest/internal/features/invest"
	"github.com/moond0wner/morning-digest/internal/features/news/habr"
	"github.com/moond0wner/morning-digest/internal/features/start"
	todoreminder "github.com/moond0wner/morning-digest/internal/features/todo_reminder"
	"github.com/moond0wner/morning-digest/internal/features/weather"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	loggerConfig := core_logger.NewConfigMust()

	logger, err := core_logger.NewLogger(loggerConfig)
	if err != nil {
		fmt.Printf("failed to init application logger: %v", err)
		os.Exit(1)
	}

	bot, err := core_telegram.NewBot(ctx, logger)
	if err != nil {
		logger.Error("create bot", zap.Error(err))
		os.Exit(1)
	}

	features := []core_feature.Feature{
		start.NewFeature(logger),
		weather.NewFeature(logger),
		invest.NewFeature(logger),
		habr.NewFeature(logger),
		todoreminder.NewFeature(logger),
		echo.NewFeature(logger), // этот хендлер должен быть всегда последним в регистрации
	}

	for _, f := range features {
		logger.Debug("initializing feature", zap.String("feature", f.Name()))
		if err := f.Init(ctx); err != nil {
			attempts := 3 // попытки
			if err := core_retry.Do(ctx, attempts, 1*time.Second, logger, f.Init); err != nil {
				logger.Warn("feature disabled",
					zap.String("feature", f.Name()),
					zap.Error(err),
				)
				continue
			}
		}
		bot.Register(f.Handle()...)
		logger.Debug("feature registered", zap.String("feature", f.Name()))
	}

	defer func() {
		for i := len(features) - 1; i >= 0; i-- {
			if err := features[i].Close(ctx); err != nil {
				logger.Warn("feature close", zap.String("feature", features[i].Name()), zap.Error(err))
			}
			logger.Debug("feature is closed", zap.String("feature", features[i].Name()))
		}
	}()

	bot.Start(ctx)
	logger.Debug("bot is stopped")
}
