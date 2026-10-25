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

	sb.WriteString("🌤 Погода\n\n")
	fmt.Fprintf(&sb, "📍 %s\n", w.City)
	fmt.Fprintf(&sb, "🌡  %.0f °C (ощущается как %.0f °C)\n", w.Temp, w.FeelsLike)
	fmt.Fprintf(&sb, "💨 Ветер: %.1f м/с\n", w.WindSpeed)
	fmt.Fprintf(&sb, "💧 Влажность: %d%%\n", w.Humidity)
	fmt.Fprintf(&sb, "☁️ %s\n", w.Description)

	return sb.String()
}
