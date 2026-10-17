package todo

import (
	"fmt"

	"github.com/go-telegram/bot/models"
)

const (
	cbCreateTask      = "task:create"
	cbGetTasks        = "task:get"
	cbShowTaskForDone = "task:list"
	cbCompleteTask    = "task:complete:"
)

func menuKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Добавить задачу", CallbackData: cbCreateTask},
				{Text: "Выполнить задачу", CallbackData: cbShowTaskForDone},
				{Text: "Список задач", CallbackData: cbGetTasks},
			},
		},
	}
}

func generateTaskKeyboard(tasks []Task) *models.InlineKeyboardMarkup {
	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{},
	}

	for _, task := range tasks {
		if !task.Completed {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
				[]models.InlineKeyboardButton{
					{Text: task.String(), CallbackData: fmt.Sprintf("%s%d", cbCompleteTask, task.ID)},
				},
			)
		}
	}

	return keyboard
}
