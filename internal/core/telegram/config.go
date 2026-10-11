package core_telegram

import (
	"errors"
	"os"
)

type config struct {
	botToken string
}

func newConfig() (config, error) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return config{}, errors.New("empty BOT_TOKEN")
	}

	return config{botToken: botToken}, nil
}
