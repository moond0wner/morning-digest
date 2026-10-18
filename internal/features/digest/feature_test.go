// internal/features/digest/feature_test.go

package digest

import (
	"context"
	"sync"
	"testing"
	"time"

	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_scheduler "github.com/moond0wner/morning-digest/internal/core/scheduler"
	"go.uber.org/zap"
)

// Фейковый отправитель, который записывает все сообщения в память
type fakeSender struct {
	mu       sync.Mutex
	messages []string
}

func (f *fakeSender) SendMessage(ctx context.Context, chatID int64, text string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.messages = append(f.messages, text)
	return nil
}

func (f *fakeSender) GetMessages() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string{}, f.messages...)
}

func TestFeature_SchedulerTicks(t *testing.T) {
	sched := core_scheduler.NewTicker(1 * time.Second)

	sender := &fakeSender{}

	logger := &core_logger.Logger{
		Logger: zap.NewNop(),
	}
	feature := NewFeature(logger, []Briefer{}, sched, Config{AccountID: 123}, sender)
	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	go feature.RunScheduler(ctx)

	<-ctx.Done()

	msgs := sender.GetMessages()
	if len(msgs) < 2 {
		t.Errorf("expected at least 2 messages, got %d", len(msgs))
	}
}
