package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
	default_handler "github.com/moond0wner/morning-digest/internal/features/default"
	"github.com/moond0wner/morning-digest/internal/features/invest"
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

	defaultHandler := default_handler.NewHandler(logger) // этот хендлер должен быть всегда последним в регистрации

	logger.Debug("initializing feature", zap.String("feature", "start"))
	startHandler := start.NewHandler(logger)

	logger.Debug("initializing feature", zap.String("feature", "weather"))
	weatherService, err := weather.NewService()
	if err != nil {
		logger.Error("failed to init weather service", zap.Error(err))
		os.Exit(1)
	}
	weatherHandler := weather.NewHandler(logger, weatherService)

	logger.Debug("initializing feature", zap.String("feature", "invest"))
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
	investHandler := invest.NewHandler(logger, investService)

	bot.Register(
		startHandler.Handle(),
		weatherHandler.Handle(),
		investHandler.Handle(),
		defaultHandler.Handle(),
	)
	bot.Start(ctx)
	logger.Debug("bot is stopped")
}
