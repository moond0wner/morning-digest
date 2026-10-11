package weather

import (
	"errors"
	"os"
)

type config struct {
	apiKey string // openweather api key
}

func newConfig() (config, error) {
	apiKey := os.Getenv("OPEN_WEATHER")
	if apiKey == "" {
		return config{}, errors.New("empty OPEN_WEATHER")
	}
	return config{apiKey: apiKey}, nil
}
