package core_auth_jwt

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Secret  string `envconfig:"SECRET"  required:"true"`
	AccessTokenTTL time.Duration `envconfig:"ACCESS_TTL" default:"15m"`
	RefreshTokenTTL time.Duration `envconfig:"REFRESH_TTL" default:"1h"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("JWT", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get JWT config: %w", err)
		panic(err)
	}

	return config
}
