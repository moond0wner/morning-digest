package weather

import (
	"fmt"

	"github.com/go-telegram/bot/models"
)

func geoKeyboard(candidates []geoResponse) *models.InlineKeyboardMarkup {
	kb := &models.InlineKeyboardMarkup{InlineKeyboard: [][]models.InlineKeyboardButton{}}
	for i, g := range candidates {
		kb.InlineKeyboard = append(kb.InlineKeyboard, []models.InlineKeyboardButton{
			{Text: fmt.Sprintf("%s (%s)", g.Name, g.Country), CallbackData: fmt.Sprintf("weather:city:%d", i)},
		})
	}
	return kb
}
