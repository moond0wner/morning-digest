package todoreminder

import "github.com/go-telegram/bot/models"

const (
	cbCreateTask = "task:create"
	cbGetTasks   = "task:get:tasks"
)

func menuKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Добавить задачу", CallbackData: cbCreateTask},
				{Text: "Список задач", CallbackData: cbGetTasks},
			},
		},
	}
}
