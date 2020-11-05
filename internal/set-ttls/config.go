package set_ttls

import "time"

// Config is a configuration struct.
type Config struct {
	Workers   int           `envconfig:"workers" default:"16"`
	Cursor    uint64        `envconfig:"cursor" default:"0"`
	Match     string        `envconfig:"match" default:"*"`
	BatchSize int64         `envconfig:"batch_size" default:"100"`
	RateLimit int           `envconfig:"rate_limit" default:"1000"`
	TTL       time.Duration `envconfig:"ttl" default:"2880h"`
	LogPeriod time.Duration `envconfig:"log_period" default:"5s"`
}
