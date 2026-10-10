package weather

import (
	"fmt"
	"strings"
)

type Weather struct {
	City        string
	Description string
	Temp        float64
	FeelsLike   float64
	Humidity    int
	WindSpeed   float64
}

func (w *Weather) String() string {

	var sb strings.Builder

	sb.WriteString("Погода\n\n")
	sb.WriteString(fmt.Sprintf("Город: %s\n", w.City))
	sb.WriteString(fmt.Sprintf("Описание: %s\n", w.Description))
	sb.WriteString(fmt.Sprintf("Температура: %.1f °C\n", w.Temp))
	sb.WriteString(fmt.Sprintf("Ощущается как: %.1f °C\n", w.FeelsLike))
	sb.WriteString(fmt.Sprintf("Влажность: %d%%\n", w.Humidity))
	sb.WriteString(fmt.Sprintf("Скорость ветра: %.2f м/с\n", w.WindSpeed))

	return sb.String()
}
