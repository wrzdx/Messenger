package auth_jwt

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Secret []byte `envconfig:"SECRET"  required:"true"`
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
