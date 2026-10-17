package weather

import (
	"context"
	"fmt"
	"strings"
	"time"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_retry "github.com/moond0wner/morning-digest/internal/core/retry"
)

type service struct {
	client *client
	log    *core_logger.Logger
}

func newService(logger *core_logger.Logger) (*service, error) {
	config, err := newConfig()
	if err != nil {
		return nil, fmt.Errorf("load weather config: %w", err)
	}

	return &service{
		client: newClient(config.apiKey),
		log:    logger,
	}, nil
}

func (s *service) get(ctx context.Context, city string) (Weather, error) {
	var weather Weather
	err := core_retry.Do(ctx, 3, 1*time.Second, s.log, func(ctx context.Context) error {
		lat, lon, err := s.client.getCoordinates(ctx, city)
		if err != nil {
			return err
		}

		response, err := s.client.getWeather(ctx, lat, lon)
		if err != nil {
			return err
		}

		weather, err = convertToWeather(response)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return Weather{}, fmt.Errorf("get weather for %s: %w", city, err)
	}

	return weather, nil

}

func convertToWeather(response *weatherResponse) (Weather, error) {
	if len(response.Weather) == 0 {
		return Weather{}, core_errors.New(
			core_errors.ErrNotFound,
			"weather response has no entries city",
		)
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
