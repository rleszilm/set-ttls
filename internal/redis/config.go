package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kelseyhightower/envconfig"
)

// Config is the struct used to parse configuration from environment variables.
type Config struct {
	Network            string        `envconfig:"network" default:"tcp"`
	Addr               string        `envconfig:"addr" default:"localhost:6379"`
	Username           string        `envconfig:"username" default:""`
	Password           string        `envconfig:"password" default:""`
	DB                 int           `envconfig:"db" default:"0"`
	MaxRetries         int           `envconfig:"max_retries" default:"0"`
	MinRetryBackoff    time.Duration `envconfig:"min_retry_backoff" default:"8ms"`
	MaxRetryBackoff    time.Duration `envconfig:"max_retry_backoff" default:"512ms"`
	DialTimeout        time.Duration `envconfig:"dial_timeout" default:"5s"`
	ReadTimeout        time.Duration `envconfig:"read_timeout" default:"5s"`
	WriteTimeout       time.Duration `envconfig:"write_timeout" default:"5s"`
	PoolSize           int           `envconfig:"pool_size" default:"100"`
	MinIdleConns       int           `envconfig:"min_idle_conns" default:"10"`
	MaxConnAge         time.Duration `envconfig:"max_conn_age" default:""`
	PoolTimeout        time.Duration `envconfig:"pool_timeout" default:""`
	IdleTimeout        time.Duration `envconfig:"idle_timeout" default:"5m"`
	IdleCheckFrequency time.Duration `envconfig:"idle_check_frequency" default:"1m"`
}

// NewFromEnv generates a new set of configuration data.
func NewFromEnv(namespace string) (*Config, error) {
	c := &Config{}
	err := envconfig.Process(namespace, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// ToV8Options converts a Config struct to a go-redis v8 Options
func ToV8Options(config *Config) *redis.Options {
	return &redis.Options{
		Network:            config.Network,
		Addr:               config.Addr,
		Username:           config.Username,
		Password:           config.Password,
		DB:                 config.DB,
		MaxRetries:         config.MaxRetries,
		MinRetryBackoff:    config.MinRetryBackoff,
		MaxRetryBackoff:    config.MaxRetryBackoff,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxConnAge:         config.MaxConnAge,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
	}
}
