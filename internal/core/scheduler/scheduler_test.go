package core_scheduler

import (
	"testing"
	"time"
)

func TestCalculateNextRun(t *testing.T) {
	moscow, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		t.Fatalf("failed to load location: %v", err)
	}

	s := &Scheduler{
		loc:          moscow,
		targetHour:   8,
		targetMinute: 0,
	}

	tests := []struct {
		name     string
		now      time.Time
		expected time.Duration
	}{
		{
			name:     "before target time today",
			now:      time.Date(2026, 10, 18, 7, 0, 0, 0, moscow),
			expected: 1 * time.Hour, // до 08:00 ровно час
		},
		{
			name:     "after target time today",
			now:      time.Date(2026, 10, 18, 9, 0, 0, 0, moscow),
			expected: 23 * time.Hour, // до завтрашних 08:00 = 23 часа
		},
		{
			name:     "exactly at target time",
			now:      time.Date(2026, 10, 18, 8, 0, 0, 0, moscow),
			expected: 24 * time.Hour, // ровно в цель — следующий запуск завтра
		},
		{
			name:     "one second before target",
			now:      time.Date(2026, 10, 18, 7, 59, 59, 0, moscow),
			expected: 1 * time.Second,
		},
		{
			name:     "one second after target",
			now:      time.Date(2026, 10, 18, 8, 0, 1, 0, moscow),
			expected: 23*time.Hour + 59*time.Minute + 59*time.Second,
		},
		{
			name:     "just after midnight",
			now:      time.Date(2026, 10, 18, 0, 30, 0, 0, moscow),
			expected: 7*time.Hour + 30*time.Minute, // до 08:00
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.calculateNextRun(tt.now)
			if got != tt.expected {
				t.Errorf("calculateNextRun(%v) = %v, want %v",
					tt.now, got, tt.expected)
			}
		})
	}
}

func TestCalculateNextRun_EdgeCases(t *testing.T) {
	moscow, _ := time.LoadLocation("Europe/Moscow")

	s := &Scheduler{
		loc:          moscow,
		targetHour:   8,
		targetMinute: 0,
	}

	tests := []struct {
		name     string
		now      time.Time
		expected time.Duration
	}{
		{
			name:     "last day of month",
			now:      time.Date(2026, 10, 31, 9, 0, 0, 0, moscow), // после 8:00
			expected: 23 * time.Hour,                              // 1 ноября 08:00
		},
		{
			name:     "last day of year",
			now:      time.Date(2026, 12, 31, 9, 0, 0, 0, moscow),
			expected: 23 * time.Hour, // 1 января 08:00
		},
		{
			name:     "leap day",
			now:      time.Date(2028, 2, 29, 9, 0, 0, 0, moscow), // 29 февраля
			expected: 23 * time.Hour,                             // 1 марта 08:00
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.calculateNextRun(tt.now)
			if got != tt.expected {
				t.Errorf("calculateNextRun(%v) = %v, want %v",
					tt.now, got, tt.expected)
			}
		})
	}
}

func TestCalculateNextRun_DST(t *testing.T) {
	berlin, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatalf("failed to load location: %v", err)
	}

	s := &Scheduler{
		loc:          berlin,
		targetHour:   8,
		targetMinute: 0,
	}

	// 25 октября 2025, 09:00 — после цели. Следующий запуск — 26 октября, но там переход.
	now := time.Date(2025, 10, 25, 9, 0, 0, 0, berlin)
	got := s.calculateNextRun(now)

	// Несмотря на то, что часов в сутках будет 25, "календарный" день длиннее на час,
	// и если сейчас 9:00, до 8:00 следующего дня = 24 часа (потому что мы спим "до 8 утра").
	// ВАЖНО: этот тест проверяет, что результат != 23 часа (как если бы часы не переводили)
	if got == 23*time.Hour {
		t.Errorf("DST not accounted for: got %v, expected something else (not 23h)", got)
	}

	// Для точности: target = 26 октября 08:00 (тоже в Berlin).
	// now = 25 октября 09:00.
	// Разница = target.Sub(now). Учитывая переход часов ночью, разница будет 24 часа.
	t.Logf("got duration: %v", got)
}
