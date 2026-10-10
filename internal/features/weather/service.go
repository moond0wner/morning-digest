package weather

import (
	"context"
	"fmt"
)

type Service struct {
	client *client
}

func NewService(apiKey string) *Service {
	return &Service{
		client: newClient(apiKey),
	}
}

func (s *Service) Get(ctx context.Context, city string) (*Weather, error) {
	lat, lon, err := s.client.getCoordinates(ctx, city)
	if err != nil {
		return nil, fmt.Errorf("get coordinates: %w", err)
	}

	response, err := s.client.getWeather(ctx, lat, lon)
	if err != nil {
		return nil, fmt.Errorf("get weather: %w", err)
	}

	weather := &Weather{
		City:        response.Name,
		Description: response.Weather[0].Description,
		Temp:        response.Main.Temp,
		FeelsLike:   response.Main.FeelsLike,
		Humidity:    response.Main.Humidity,
		WindSpeed:   response.Wind.Speed,
	}

	return weather, nil

}
