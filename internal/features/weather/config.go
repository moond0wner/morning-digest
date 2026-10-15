package weather

import (
	"os"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
)

type config struct {
	apiKey string // openweather api key
}

func newConfig() (config, error) {
	apiKey := os.Getenv("OPEN_WEATHER")
	if apiKey == "" {
		return config{}, core_errors.New(core_errors.ErrInternal, "empty OPEN_WEATHER")
	}
	return config{apiKey: apiKey}, nil
}
