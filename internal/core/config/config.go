package core_config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BotToken    string `envconfig:"BOT_TOKEN" required:"true"`
	OpenWeather string `envconfig:"OPEN_WEATHER" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("API", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}
	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get core config: %w", err))
	}
	return config
}
