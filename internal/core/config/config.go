package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TimeZoneName string      `envconfig:"TIME_ZONE" default:"UTC"`
	Environment  Environment `envconfig:"ENVIRONMENT" default:"development"`

	TimeZone *time.Location `ignored:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("process envconfig: %w", err)
	}

	location, err := time.LoadLocation(cfg.TimeZoneName)
	if err != nil {
		return nil, fmt.Errorf(
			"load time zone %q: %w",
			cfg.TimeZoneName,
			err,
		)
	}

	cfg.TimeZone = location

	return &cfg, nil
}

func NewConfigMust() *Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get core config: %w", err)
		panic(err)
	}
	return config
}

type Environment string

const (
	Development Environment = "development"
	Testing     Environment = "testing"
	Production  Environment = "production"
)

func (e Environment) IsProduction() bool {
	return e == Production
}

func (e Environment) IsDevelopment() bool {
	return e == Development
}

func (e Environment) IsTesting() bool {
	return e == Testing
}
