package core_config

import (
	"os"
	"strconv"
	"time"

	core_errors "github.com/moond0wner/morning-digest/internal/core/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	RetryDelay    Duration `yaml:"retry_delay"`
	RetryBackoff  Duration `yaml:"retry_backoff"`
	RetryAttempts int      `yaml:"retry_attempts"`
	InitTimeout   Duration `yaml:"init_timeout"`

	TimezoneStr string `yaml:"timezone"`

	Timezone *time.Location `yaml:"-"`

	AccountID int64
}

func newConfig() (*Config, error) {
	yamlConfig, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"read yaml config",
		)
	}

	var config *Config
	if err := yaml.Unmarshal(yamlConfig, &config); err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"yaml unmarshal",
		)
	}

	loc, err := time.LoadLocation(config.TimezoneStr)
	if err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"load timezone",
		)
	}
	config.Timezone = loc

	tgIDStr := os.Getenv("TELEGRAM_ACCOUNT_ID")
	if tgIDStr == "" {
		return nil, core_errors.New(
			core_errors.ErrInternal,
			"empty TELEGRAM_ACCOUNT_ID",
		)
	}

	tgID, err := strconv.ParseInt(tgIDStr, 10, 64)
	if err != nil {
		return nil, core_errors.Wrap(
			core_errors.ErrInternal,
			err,
			"error convertation tgID",
		)
	}

	config.AccountID = tgID

	return config, nil
}

func NewConfigMust() *Config {
	cfg, err := newConfig()
	if err != nil {
		panic(err)
	}

	return cfg
}
