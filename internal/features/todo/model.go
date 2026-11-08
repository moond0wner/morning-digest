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
	CreatedAt   time.Time
	CompletedAt *time.Time
	Completed   bool
}

func newTaskUninitialized() *Task {
	return &Task{
		ID:          uninitializedID,
		Title:       "",
		Description: "",
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Completed:   false,
	}
}

func (t *Task) String() string {
	var sb strings.Builder

	status := "🔴"
	if t.Completed {
		status = "🟢"
	}

	fmt.Fprintf(&sb, "%s %s", status, t.Title)

	if t.Description != "" {
		fmt.Fprintf(&sb, "\n   └ %s", t.Description)
	}

	if t.Completed && t.CompletedAt != nil {
		fmt.Fprintf(&sb, "\n   ⏱ Выполнено: %s", t.CompletedAt.Format("02 Jan 15:04"))
	}

	return sb.String()
}
