package core_retry

import (
	"context"
	"errors"
	"time"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	"go.uber.org/zap"
)

type Retrier struct {
	attempts   int
	delay      time.Duration
	maxBackoff time.Duration
	log        *core_logger.Logger
}

func NewRetrier(attempts int, delay, maxBackoff time.Duration, log *core_logger.Logger) *Retrier {
	return &Retrier{
		attempts:   attempts,
		delay:      delay,
		maxBackoff: maxBackoff,
		log:        log,
	}
}

func (r *Retrier) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	var lastErr error
	for i := 0; i < r.attempts; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err := fn(ctx)
		if err == nil {
			if i > 0 {
				r.log.Debug("retry succeeded",
					zap.Int("attempt", i+1),
					zap.Int("attempts", r.attempts),
				)
			}
			return nil
		}
		lastErr = err

		if !isRetryable(err) {
			return err
		}

		r.log.Debug("retry attempt failed",
			zap.Int("attempt", i+1),
			zap.Int("attempts", r.attempts),
			zap.Error(err),
		)

		backoff := r.delay * (1 << i)
		if backoff > r.maxBackoff {
			backoff = r.maxBackoff
		}
		if i < r.attempts-1 {
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
