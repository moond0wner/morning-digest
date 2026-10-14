package habr

import "github.com/go-telegram/bot/models"

const (
	cbSourceBase  = "news:source:base"
	cbSourceHub   = "news:source:hub"
	cbHubPrefix   = "news:hub:"
	cbLimitPrefix = "news:limit:"
)

func sourceKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Базовое", CallbackData: cbSourceBase},
				{Text: "По хабу", CallbackData: cbSourceHub},
			},
		},
	}
}

func hubKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Go", CallbackData: "news:hub:go"},
				{Text: "C++", CallbackData: "news:hub:cpp"},
				{Text: "Rust", CallbackData: "news:hub:rust"},
			},
		},
	}
}

func limitKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "3", CallbackData: "news:limit:3"},
				{Text: "5", CallbackData: "news:limit:5"},
				{Text: "10", CallbackData: "news:limit:10"},
			},
		},
	}
}
