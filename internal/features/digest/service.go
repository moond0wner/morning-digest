package digest

import (
	"context"
	"strings"

	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	"go.uber.org/zap"
)

type service struct {
	briefers []Briefer
	log      *core_logger.Logger
}

func newService(briefers []Briefer, log *core_logger.Logger) *service {
	return &service{
		briefers: briefers,
		log:      log,
	}
}

func (s *service) build(ctx context.Context) (string, error) {
	var sb strings.Builder
	sb.WriteString("Доброе утро!\n\n")

	for _, b := range s.briefers {
		text, err := b.Brief(ctx)
		if err != nil {
			s.log.Warn("brief failed", zap.String("feature", b.Name()), zap.Error(err))
			continue
		}
		if text == "" {
			continue
		}
		sb.WriteString(text)
		sb.WriteString("\n\n")
	}

	return strings.TrimSpace(sb.String()), nil
}
