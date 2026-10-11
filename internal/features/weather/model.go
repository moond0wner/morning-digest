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

func (w Weather) String() string {
	var sb strings.Builder

	sb.WriteString("Погода\n\n")
	fmt.Fprintf(&sb, "Город: %s\n", w.City)
	fmt.Fprintf(&sb, "Описание: %s\n", w.Description)
	fmt.Fprintf(&sb, "Температура: %.1f °C\n", w.Temp)
	fmt.Fprintf(&sb, "Ощущается как: %.1f °C\n", w.FeelsLike)
	fmt.Fprintf(&sb, "Влажность: %d%%\n", w.Humidity)
	fmt.Fprintf(&sb, "Скорость ветра: %.2f м/с\n", w.WindSpeed)

	return sb.String()
}
