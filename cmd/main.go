package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_config "github.com/moond0wner/morning-digest/internal/core/config"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
	"github.com/moond0wner/morning-digest/internal/features/invest"
	"github.com/moond0wner/morning-digest/internal/features/start"
	"github.com/moond0wner/morning-digest/internal/features/weather"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Printf("Failed to init application logger: %v", err)
		os.Exit(1)
	}

	config := core_config.NewConfigMust()
	logger.Debug("application config")

	bot, err := core_telegram.NewBot(logger, config)
	if err != nil {
		logger.Error("create bot", zap.Error(err))
		os.Exit(1)
	}
	logger.Debug("application bot")

	logger.Debug("initializing feature", zap.String("feature", "start"))
	startHandler := start.NewHandler(logger)

	logger.Debug("initializing feature", zap.String("feature", "weather"))
	weatherService := weather.NewService(config.OpenWeather)
	weatherTransport := weather.NewHandler(logger, weatherService)

	logger.Debug("initializing feature", zap.String("feature", "invest"))
	investService := invest.NewService(ctx, logger)
	investTransport := invest.NewHandler(logger, investService)

	bot.Register(
		startHandler.Handle(),
		weatherTransport.Handle(),
		investTransport.Handle(),
	)
	bot.Start(ctx)
}
