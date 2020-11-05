package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rleszilm/set-ttls/internal/redis"
	set_ttls "github.com/rleszilm/set-ttls/internal/set-ttls"
)

// Config is a configuration struct.
type Config struct {
	set_ttls.Config
	Redis redis.Config
}

// NewFromEnv returns a new config based on env vars.
func NewFromEnv(prefix string) (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process(prefix, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
