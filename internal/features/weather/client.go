package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	baseGeoUrl     = "http://api.openweathermap.org/geo/1.0/direct?"
	baseWeatherUrl = "https://api.openweathermap.org/data/2.5/weather"
)

type client struct {
	apiKey string
	http   *http.Client
}

func newClient(apiKey string) *client {
	return &client{
		apiKey: apiKey,
		http:   &http.Client{Timeout: 10 * time.Second},
	}
}

type geoResponse struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
}

type weatherResponse struct {
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`

	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`

	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`

	Name string `json:"name"`
}

func (c *client) getCoordinates(ctx context.Context, city string) (lat, lot float64, err error) {
	url, err := url.Parse(baseGeoUrl)
	if err != nil {
		return 0, 0, fmt.Errorf("parse geocoding url: %w", err)
	}
	q := url.Query()
	q.Set("q", city)
	q.Set("limit", "1")
	q.Set("appid", c.apiKey)
	url.RawQuery = q.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return 0, 0, fmt.Errorf("create geocoding request: %w", err)
	}

	response, err := c.http.Do(request)
	if err != nil {
		return 0, 0, fmt.Errorf("do geocoding request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("geocoding status: %s", response.Status)
	}

	var geo []geoResponse
	if err := json.NewDecoder(response.Body).Decode(&geo); err != nil {
		return 0, 0, fmt.Errorf("decode geocoding response: %w", err)
	}
	if len(geo) == 0 {
		return 0, 0, fmt.Errorf("city not found: %s", city)
	}

	return geo[0].Lat, geo[0].Lon, nil
}

func (c *client) getWeather(ctx context.Context, lat, lon float64) (*weatherResponse, error) {
	url, err := url.Parse(baseWeatherUrl)
	if err != nil {
		return nil, fmt.Errorf("parse weather url: %w", err)
	}

	q := url.Query()
	q.Set("lat", fmt.Sprintf("%f", lat))
	q.Set("lon", fmt.Sprintf("%f", lon))
	q.Set("appid", c.apiKey)
	q.Set("units", "metric")
	q.Set("lang", "ru")
	url.RawQuery = q.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create weather request: %w", err)
	}

	response, err := c.http.Do(request)
	if err != nil {
		return nil, fmt.Errorf("do weather request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather status: %s", response.Status)
	}

	var w weatherResponse
	if err := json.NewDecoder(response.Body).Decode(&w); err != nil {
		return nil, fmt.Errorf("decode weather response: %w", err)
	}

	return &w, nil
}
