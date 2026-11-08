package digest

import (
	"context"
	"strings"
	"sync"
	"time"

	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	"go.uber.org/zap"
)

type service struct {
	log      *core_logger.Logger
	briefers []Briefer
}

func newService(briefers []Briefer, log *core_logger.Logger) *service {
	return &service{
		briefers: briefers,
		log:      log,
	}
}

func (s *service) build(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	type result struct {
		text string
	}
	resultCh := make(chan result, len(s.briefers))

	wg := &sync.WaitGroup{}
	for _, b := range s.briefers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			text, err := b.Brief(ctx)
			if err != nil {
				s.log.Warn("brief failed", zap.String("feature", b.Name()), zap.Error(err))
				return
			}
			if text != "" {
				resultCh <- result{text: text}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var sb strings.Builder
	sb.WriteString("Доброе утро!\n\n")

	for v := range resultCh {
		sb.WriteString(v.text)
		sb.WriteString("\n\n")
	}

	return strings.TrimSpace(sb.String()), nil
}
