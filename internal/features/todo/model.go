package todo

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
	// TODO: Криво, исправить
	var sb strings.Builder

	status := "❌"
	if t.Completed {
		status = "✅"
	}

	sb.WriteString(fmt.Sprintf("№%d - '%s': %s %s", t.ID, t.Title, t.Description, status))

	if t.CompletedAt != nil {
		sb.WriteString(fmt.Sprintf("\n✅ завершено %s\n", t.CompletedAt.Format("02 Jan 15:04")))
	}

	return sb.String()
}
