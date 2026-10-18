package core_scheduler

import (
	"context"
	"time"
)

type Scheduler struct {
	loc            *time.Location
	targetHour     int
	targetMinute   int
	tickerDuration time.Duration
}

func NewSchedueler(loc *time.Location, targetHour, targetMinute int) *Scheduler {
	return &Scheduler{
		loc:          loc,
		targetHour:   targetHour,
		targetMinute: targetMinute,
	}
}

func NewTicker(interval time.Duration) *Scheduler {
	return &Scheduler{tickerDuration: interval}
}

func (s *Scheduler) Run(ctx context.Context, onTick func()) {
	for {
		var timer *time.Timer

		// Если задан tickerDuration, значит мы в тестовом режиме
		if s.tickerDuration > 0 {
			timer = time.NewTimer(s.tickerDuration)
		} else {
			// Иначе - основной режим по расписанию
			now := time.Now()
			duration := s.calculateNextRun(now)
			timer = time.NewTimer(duration)
		}

		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			onTick()
		}
	}
}

func (s *Scheduler) calculateNextRun(now time.Time) time.Duration {
	now = now.In(s.loc)
	target := time.Date(
		now.Year(), now.Month(), now.Day(),
		s.targetHour, s.targetMinute, 0, 0,
		s.loc,
	)

	if !target.After(now) {
		target = target.AddDate(0, 0, 1)
	}

	return target.Sub(now)
}
