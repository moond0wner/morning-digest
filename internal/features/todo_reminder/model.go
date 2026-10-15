package todoreminder

import (
	"fmt"
	"strings"
	"time"
)

var (
	uninitializedID int64 = -1
)

type Task struct {
	ID          int64
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

func newTaskUninitialized() *Task {
	return &Task{
		uninitializedID,
		"",
		"",
		false,
		time.Now(),
		nil,
	}
}

func (t *Task) String() string {
	var sb strings.Builder

	status := "⬜"
	if t.Completed {
		status = "✅"
	}

	sb.WriteString(fmt.Sprintf("%s #%d — %s", status, t.ID, t.Title))

	if t.Description != "" {
		sb.WriteString("\n " + t.Description)
	}

	if t.CompletedAt != nil {
		sb.WriteString(fmt.Sprintf("\n✅ завершено %s", t.CompletedAt.Format("02 Jan 15:04")))
	}

	return sb.String()
}
