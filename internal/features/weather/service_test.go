package weather

import (
	"context"
	"errors"
	"testing"
	"time"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
	core_logger "github.com/moond0wner/morning-digest/internal/core/logger"
	core_retry "github.com/moond0wner/morning-digest/internal/core/retry"
	core_testutils "github.com/moond0wner/morning-digest/internal/core/testutils"
	"go.uber.org/zap"
)

// 1. Создаем мок, реализующий интерфейс weatherClient
type mockWeatherClient struct {
	geo        []geoResponse
	geoErr     error
	weather    *weatherResponse
	weatherErr error
}

func (m *mockWeatherClient) getCoordinates(ctx context.Context, city string) ([]geoResponse, error) {
	return m.geo, m.geoErr
}

func (m *mockWeatherClient) getWeather(ctx context.Context, lat, lon float64) (*weatherResponse, error) {
	return m.weather, m.weatherErr
}

func TestService_GetByCoords(t *testing.T) {
	logger := &core_logger.Logger{
		Logger: zap.NewNop(),
	}
	retry := core_retry.NewRetrier(1, 1*time.Second, 1*time.Second, logger)
	t.Run("успешное получение погоды", func(t *testing.T) {
		mock := &mockWeatherClient{
			weather: &weatherResponse{
				Name: "Саратов",
				Main: struct {
					Temp      float64 `json:"temp"`
					FeelsLike float64 `json:"feels_like"`
					Humidity  int     `json:"humidity"`
				}{Temp: 25.5, FeelsLike: 24.0, Humidity: 50},
				Weather: []struct {
					Description string `json:"description"`
				}{{Description: "ясно"}},
			},
		}

		svc, err := newService(mock, logger, retry)
		if err != nil {
			t.Fatalf("newService: %v", err)
		}

		got, err := svc.getByCoords(context.Background(), 51.54, 45.99)
		if err != nil {
			t.Fatalf("getByCoords вернул ошибку: %v", err)
		}
		core_testutils.AssertEqual(t, got.City, mock.weather.Name)
		core_testutils.AssertEqual(t, got.Temp, mock.weather.Main.Temp)
		core_testutils.AssertEqual(t, got.Description, mock.weather.Weather[0].Description)
	})

	t.Run("ошибка в клиенте", func(t *testing.T) {
		mock := &mockWeatherClient{
			weatherErr: core_errors.New(core_errors.ErrExternalAPI, "api down"),
		}

		svc, err := newService(mock, logger, retry)
		if err != nil {
			t.Fatalf("newService: %v", err)
		}
		_, err = svc.getByCoords(context.Background(), 51.54, 45.99)
		if err == nil {
			t.Fatal("ожидали ошибку, получили nil")
		}

		targetErr := &core_errors.Error{Type: core_errors.ErrExternalAPI}
		if !errors.Is(err, targetErr) {
			t.Fatalf("ожидали ErrExternalAPI, получили: %v", err)
		}
	})
}
