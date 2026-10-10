package weather

import "fmt"

type Weather struct {
	City        string
	Description string
	Temp        float64
	FeelsLike   float64
	Humidity    int
	WindSpeed   float64
}

func (w *Weather) String() string {
	text := fmt.Sprintf(
		"Город: %s\n"+
			"Описание: %s\n"+
			"Температура: %f\n"+
			"Ощущается как: %f\n"+
			"Влажность: %d%%\n"+
			"Скорость ветра: %f м/с\n",
		w.City, w.Description, w.Temp, w.FeelsLike, w.Humidity, w.WindSpeed,
	)
	return text
}
