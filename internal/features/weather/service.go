package weather

import (
	"context"
	"fmt"
	"strings"
)

type Service struct {
	client *client
}

func NewService() (*Service, error) {
	config, err := newConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load weather config: %w", err)
	}

	return &Service{
		client: newClient(config.apiKey),
	}, nil
}

func (s *Service) Get(ctx context.Context, city string) (Weather, error) {
	lat, lon, err := s.client.getCoordinates(ctx, city)
	if err != nil {
		return Weather{}, fmt.Errorf("get coordinates for %s: %w", city, err)
	}

	response, err := s.client.getWeather(ctx, lat, lon)
	if err != nil {
		return Weather{}, fmt.Errorf("get weather for %s: %w", city, err)
	}

	if len(response.Weather) == 0 {
		return Weather{}, fmt.Errorf("weather response has no entries for city %q", city)
	}
	descriptions := make([]string, len(response.Weather))
	for i, w := range response.Weather {
		descriptions[i] = w.Description
	}
	description := strings.Join(descriptions, ", ")

	weather := Weather{
		City:        response.Name,
		Description: description,
		Temp:        response.Main.Temp,
		FeelsLike:   response.Main.FeelsLike,
		Humidity:    response.Main.Humidity,
		WindSpeed:   response.Wind.Speed,
	}

	return weather, nil

}
