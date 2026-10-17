package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/moond0wner/morning-digest/internal/core/config"
	core_feature "github.com/moond0wner/morning-digest/internal/core/feature"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_repository "github.com/moond0wner/morning-digest/internal/core/repository"
	core_retry "github.com/moond0wner/morning-digest/internal/core/retry"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
	"github.com/moond0wner/morning-digest/internal/features/digest"
	"github.com/moond0wner/morning-digest/internal/features/echo"
	"github.com/moond0wner/morning-digest/internal/features/invest"
	"github.com/moond0wner/morning-digest/internal/features/news/habr"
	"github.com/moond0wner/morning-digest/internal/features/start"
	"github.com/moond0wner/morning-digest/internal/features/todo"
	"github.com/moond0wner/morning-digest/internal/features/weather"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.Timezone

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	loggerConfig := core_logger.NewConfigMust()

	logger, err := core_logger.NewLogger(loggerConfig)
	if err != nil {
		fmt.Printf("failed to init application logger: %v", err)
		os.Exit(1)
	}
	logger.Debug("application time zone", zap.Any("zone", time.Local))

	retrier := core_retry.NewRetrier(cfg.RetryAttempts, cfg.RetryDelay.Std(), cfg.RetryBackoff.Std(), logger)

	bot, err := core_telegram.NewBot(ctx, logger, retrier)
	if err != nil {
		logger.Error("create bot", zap.Error(err))
		os.Exit(1)
	}

	todoBaseRepo, err := core_repository.NewBaseRepository("todo")
	if err != nil {
		logger.Error("create repository", zap.Error(err))
		os.Exit(1)
	}
	defer todoBaseRepo.Close()

	features := []core_feature.Feature{
		start.NewFeature(logger),
		weather.NewFeature(logger, retrier),
		invest.NewFeature(logger, retrier),
		habr.NewFeature(logger, retrier),
		todo.NewFeature(logger, todoBaseRepo),
	}

	var briefers []digest.Briefer
	for _, f := range features {
		if b, ok := f.(digest.Briefer); ok {
			briefers = append(briefers, b)
		}
	}

	digestFeature := digest.NewFeature(logger, briefers)
	features = append(features, digestFeature)
	features = append(features, echo.NewFeature(logger)) // этот хендлер должен быть всегда последним в регистрации

	for _, f := range features {
		logger.Debug("initializing feature", zap.String("feature", f.Name()))
		if err := f.Init(ctx); err != nil {
			if err := retrier.Do(ctx, f.Init); err != nil {
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
