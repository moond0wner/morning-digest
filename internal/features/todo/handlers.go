package todo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	core_telegram "github.com/moond0wner/morning-digest/internal/core/telegram"
)

func (h *handler) handleCreateTask(
	ctx context.Context,
	b *bot.Bot,
	chatID int64,
	userID int64,
	msgID int,
) {
	state, ok := h.getState(userID)
	if !ok {
		return
	}
	state.Step = inputTitle
	h.EditMessage(ctx, b, chatID, msgID, "Введите название задачи")
}

func (h *handler) handleGetTasks(
	ctx context.Context,
	b *bot.Bot,
	chatID int64,
	userID int64,
	msgID int,
) {
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
		fmt.Fprintln(&sb, "задач нет")
	}

	h.EditMessage(ctx, b, chatID, msgID, sb.String())
	h.deleteState(userID)
}

func (h *handler) handleShowTasksForDone(
	ctx context.Context,
	b *bot.Bot,
	chatID int64,
	userID int64,
	msgID int,
) {
	tasks, err := h.service.getTasks(ctx)
	if err != nil {
		h.ReplyError(ctx, b, chatID, "Не удалось получить список задач", err)
		return
	}

	keyboard := generateTaskKeyboard(tasks)
	if len(keyboard.InlineKeyboard) == 0 {
		h.EditMessage(ctx, b, chatID, msgID, "Задач нет", core_telegram.DeleteKeyboard())
		h.deleteState(userID)
		return
	}

	h.EditMessage(ctx, b, chatID, msgID, "Выберите задачу которую хотите выполнить", keyboard)
}

func (h *handler) handleCompleteTask(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	chatID int64,
	userID int64,
	msgID int,
) {
	data := update.CallbackQuery.Data
	idStr := strings.TrimPrefix(data, "task:complete:")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.ReplyError(ctx, b, chatID, "Не удалось выполнить задачу", err)
		return
	}

	if err := h.service.updateTask(ctx, id); err != nil {
		h.ReplyError(ctx, b, chatID, "Не удалось выполнить задачу", err)
		return
	}

	tasks, err := h.service.getTasks(ctx)
	if err != nil {
		h.ReplyError(ctx, b, chatID, "Не удалось получить список задач", err)
		return
	}

	keyboard := generateTaskKeyboard(tasks)
	if len(keyboard.InlineKeyboard) == 0 {
		h.EditMessage(ctx, b, chatID, msgID, "Задач нет", core_telegram.DeleteKeyboard())
		h.deleteState(userID)
		return
	}

	h.EditMessage(ctx, b, chatID, msgID, "Выберите задачу которую хотите выполнить", keyboard)
}

func (h *handler) handleTitle(
	ctx context.Context,
	b *bot.Bot,
	chatID int64,
	userID int64,
	text string,
) {
	state, ok := h.getState(userID)
	if !ok {
		return
	}

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
	chatID int64,
	userID int64,
	text string,
) {
	state, ok := h.getState(userID)
	if !ok {
		return
	}

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

	if err := h.service.createTask(ctx, state.Task); err != nil {
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
