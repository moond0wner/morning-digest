package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
)

const (
	baseGeoURL     = "http://api.openweathermap.org/geo/1.0/direct?"
	baseWeatherURL = "https://api.openweathermap.org/data/2.5/weather"
)

type client struct {
	apiKey         string
	http           *http.Client
	geoBaseURL     string
	weatherBaseURL string
}

func newClient(apiKey, geoURL, weatherURL string) *client {
	return &client{
		apiKey:         apiKey,
		http:           &http.Client{Timeout: 10 * time.Second},
		geoBaseURL:     geoURL,
		weatherBaseURL: weatherURL,
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

func (c *client) getCoordinates(ctx context.Context, city string) ([]geoResponse, error) {
	params := map[string]string{
		"q":     city,
		"limit": "5",
		"appid": c.apiKey,
	}

	url, err := buildURL(c.geoBaseURL, params)
	if err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"parse geocoding url",
		)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrExternalAPI,
			err,
			"create geocoding request",
		)
	}

	response, err := c.http.Do(request)
	if err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrExternalAPI,
			err,
			"do geocoding request",
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, core_errors.Wrap(
			core_errors.CheckStatusCode(response.StatusCode),
			fmt.Errorf("status %s", response.Status),
			"geocoding status",
		)
	}

	var geo []geoResponse
	if err := json.NewDecoder(response.Body).Decode(&geo); err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"decode geocoding response",
		)
	}
	if len(geo) == 0 {
		return nil, core_errors.Wrap(
			core_errors.ErrNotFound,
			fmt.Errorf("city %q", city),
			fmt.Sprintf("city not found: %s", city),
		)
	}
	return geo, nil
}

func (c *client) getWeather(ctx context.Context, lat, lon float64) (*weatherResponse, error) {
	params := map[string]string{
		"lat":   fmt.Sprintf("%f", lat),
		"lon":   fmt.Sprintf("%f", lon),
		"appid": c.apiKey,
		"units": "metric",
		"lang":  "ru",
	}

	url, err := buildURL(c.weatherBaseURL, params)
	if err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"parse weather url",
		)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"create weather request",
		)
	}

	response, err := c.http.Do(request)
	if err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrExternalAPI,
			err,
			"do weather request",
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, core_errors.Wrap(
			core_errors.CheckStatusCode(response.StatusCode),
			fmt.Errorf("status %s", response.Status),
			"weather status",
		)
	}

	var w weatherResponse
	if err := json.NewDecoder(response.Body).Decode(&w); err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"decode weather response",
		)
	}

	return &w, nil
}

func buildURL(base string, params map[string]string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
