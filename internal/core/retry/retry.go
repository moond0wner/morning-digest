package core_retry

import (
	"context"
	"errors"
	"time"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	"go.uber.org/zap"
)

func Do(
	ctx context.Context,
	attempts int,
	delay time.Duration,
	log *core_logger.Logger,
	fn func(ctx context.Context) error,
) error {
	var lastErr error
	for i := 0; i < attempts; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err := fn(ctx)
		if err == nil {
			if i > 0 {
				log.Debug("retry succeeded",
					zap.Int("attempt", i+1),
					zap.Int("attempts", attempts),
				)
			}
			return nil
		}
		lastErr = err

		if !isRetryable(err) {
			return err
		}

		log.Debug("retry attempt failed",
			zap.Int("attempt", i+1),
			zap.Int("attempts", attempts),
			zap.Error(err),
		)

		backoff := delay * (1 << i)
		if i < attempts-1 {
			timer := time.NewTimer(backoff)
			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			}
		}

	}
	return lastErr
}

func isRetryable(err error) bool {
	var e *core_errors.Error
	if errors.As(err, &e) {
		return e.Type == core_errors.ErrExternalAPI
	}
	return errors.Is(err, context.DeadlineExceeded)
}
