package digest

import (
	"context"

	core_feature "github.com/moond0wner/morning-digest/internal/core/feature"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_scheduler "github.com/moond0wner/morning-digest/internal/core/scheduler"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
	"go.uber.org/zap"
)

type MessageSender interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

type Config struct {
	AccountID int64
}

type Feature struct {
	*core_feature.BaseFeature
	handler   *handler
	service   *service
	scheduler *core_scheduler.Scheduler
	cfg       Config
	sender    MessageSender
}

func NewFeature(
	log *core_logger.Logger,
	briefers []Briefer,
	scheduler *core_scheduler.Scheduler,
	cfg Config,
	sender MessageSender,
) *Feature {
	return &Feature{
		BaseFeature: core_feature.NewBaseFeature("digest", log, nil),
		service:     newService(briefers, log),
		scheduler:   scheduler,
		cfg:         cfg,
		sender:      sender,
	}
}

func (f *Feature) Init(ctx context.Context) error {
	f.handler = newHandler(f.Log, f.service)
	return nil
}

func (f *Feature) RunScheduler(ctx context.Context) {
	f.scheduler.Run(ctx, func() {
		f.sendDigest(ctx)
	})
}

func (f *Feature) sendDigest(ctx context.Context) {
	digest, err := f.service.build(ctx)
	if err != nil {
		f.Log.Error("build digest", zap.Error(err))
		return
	}

	if err := f.sender.SendMessage(ctx, f.cfg.AccountID, digest); err != nil {
		f.Log.Error("send digest to user", zap.Error(err))
	}
}

func (f *Feature) Handle() []core_telegram.Handle {
	if f.handler == nil {
		return nil
	}
	return f.handler.handle()
}
