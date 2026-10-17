package todo

import (
	"context"
	"strings"

	core_feature "github.com/moond0wner/morning-digest/internal/core/feature"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_repository "github.com/moond0wner/morning-digest/internal/core/repository"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
)

type Feature struct {
	core_feature.BaseFeature
	handler  *handler
	baseRepo *core_repository.BaseRepository
}

func NewFeature(logger *core_logger.Logger, baseRepo *core_repository.BaseRepository) *Feature {
	return &Feature{
		BaseFeature: core_feature.NewBaseFeature("todo", logger, nil),
		baseRepo:    baseRepo,
	}
}

func (f *Feature) Init(ctx context.Context) error {
	repo := newRepository(f.baseRepo)
	service := newService(repo)
	f.handler = newHandler(f.Log, service)

	err := f.handler.service.repository.createTable(ctx)
	return err
}

func (f *Feature) Handle() []core_telegram.Handle {
	if f.handler == nil {
		return nil
	}

	return f.handler.handle()
}

func (f *Feature) Close(ctx context.Context) error {
	if f.handler == nil || f.handler.service == nil {
		return nil
	}
	return f.handler.service.close()
}

func (f *Feature) Brief(ctx context.Context) (string, error) {
	tasks, err := f.handler.service.getTasks(ctx)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString("✅ Задачи:\n")
	count := 0
	for _, t := range tasks {
		if t.Completed {
			continue
		}
		sb.WriteString("• " + t.Title + "\n")
		count++
		if count >= 5 {
			break
		}
	}

	if count == 0 {
		return "✅ Задач нет", nil
	}
	return strings.TrimSpace(sb.String()), nil
}
