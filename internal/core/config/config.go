package core_config

import (
	"fmt"
	"os"
)

type Config struct {
	BotToken string
}

func NewConfig() (*Config, error) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("Failed load bot token: empty")
	}

	return &Config{
		BotToken: botToken,
	}, nil
}

func NewConfigMust() *Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get core config: %w", err))
	}
	return config
}
