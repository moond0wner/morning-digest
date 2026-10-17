package core_telegram

import "github.com/go-telegram/bot/models"

func DeleteKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{},
	}
}
