package todoreminder

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *handler) handleChoiceButton(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	state *taskDialogState,
	chatID int64,
	userID int64,
	msgID int,
) {
	data := update.CallbackQuery.Data
	switch data {
	case cbCreateTask:
		state.Step = inputTitle
		h.EditMessage(ctx, b, chatID, msgID, "Введите название задачи")
	case cbGetTasks:
		tasks, err := h.service.getTasks(ctx)
		if err != nil {
			h.ReplyError(ctx, b, chatID, "Не удалось получить список задач", err)
			h.deleteState(userID)
			return
		}

		var sb strings.Builder
		fmt.Fprintln(&sb, "Список задач:")
		if len(tasks) > 0 {
			for _, task := range tasks {
				fmt.Fprintln(&sb, task.String())
			}
		} else {
			fmt.Fprintln(&sb, "Задач нет")
		}

		h.EditMessage(ctx, b, chatID, msgID, sb.String())
		h.deleteState(userID)
	}
}

func (h *handler) handleTitle(
	ctx context.Context,
	b *bot.Bot,
	state *taskDialogState,
	chatID int64,
	userID int64,
	text string,
) {
	title := strings.TrimSpace(text)
	if title == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Укажите название задачи",
		})
		h.deleteState(userID)
		return
	}

	state.Step = inputDescription

	h.mtx.Lock()
	state.Task.Title = title
	h.mtx.Unlock()

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Введите описание задачи",
	})
}

func (h *handler) handleDescription(
	ctx context.Context,
	b *bot.Bot,
	state *taskDialogState,
	chatID int64,
	userID int64,
	text string,
) {
	description := strings.TrimSpace(text)
	if description == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Укажите описание задачи",
		})
		h.deleteState(userID)
		return
	}

	h.mtx.Lock()
	state.Task.Description = description
	h.mtx.Unlock()

	if err := h.service.createTask(ctx, *state.Task); err != nil {
		h.ReplyError(ctx, b, chatID, "Не удалось создать задачу", err)
		h.deleteState(userID)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Задача успешно создана!",
	})
	h.deleteState(userID)
}
