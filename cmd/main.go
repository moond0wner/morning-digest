package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_feature "github.com/moond0wner/morning-digest/internal/core/feature"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
	"github.com/moond0wner/morning-digest/internal/features/echo"
	"github.com/moond0wner/morning-digest/internal/features/invest"
	"github.com/moond0wner/morning-digest/internal/features/news/habr"
	"github.com/moond0wner/morning-digest/internal/features/start"
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

	bot, err := core_telegram.NewBot(logger)
	if err != nil {
		logger.Error("create bot", zap.Error(err))
		os.Exit(1)
	}

	weatherService, err := weather.NewService()
	if err != nil {
		logger.Error("failed to init weather service", zap.Error(err))
		os.Exit(1)
	}

	investService, err := invest.NewService(ctx, logger)
	if err != nil {
		logger.Error("failed to init invest service", zap.Error(err))
		os.Exit(1)
	}

	defer func() {
		if err := investService.Close(); err != nil {
			logger.Warn("failed to close invest service", zap.Error(err))
		}
	}()

	habrService := habr.NewService()

	features := []core_feature.Feature{
		start.NewFeature(logger),
		weather.NewFeature(logger, weatherService),
		invest.NewFeature(logger, investService),
		habr.NewFeature(logger, habrService),
		echo.NewFeature(logger), // этот хендлер должен быть всегда последним в регистрации
	}

	for _, f := range features {
		logger.Debug("initializing feature", zap.String("feature", f.Name()))
		if err := f.Init(ctx); err != nil {
			logger.Warn("feature disabled",
				zap.String("feature", f.Name()),
				zap.Error(err),
			)
			continue
		}
		bot.Register(f.Handle()...)
		logger.Debug("feature registered", zap.String("feature", f.Name()))
	}

	bot.Start(ctx)
	logger.Debug("bot is stopped")
}
