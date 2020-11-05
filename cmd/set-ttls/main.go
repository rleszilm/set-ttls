package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/rleszilm/set-ttls/cmd/set-ttls/config"
	ttl_redis "github.com/rleszilm/set-ttls/internal/redis"
	set_ttls "github.com/rleszilm/set-ttls/internal/set-ttls"
)

func main() {
	ctx := context.Background()

	// get config
	cfg, err := config.NewFromEnv("SET_TTLS")
	if err != nil {
		log.Fatalln("Could not parse config:", err)
	}

	pword := cfg.Redis.Password
	cfg.Redis.Password = "fake-password-for-logging"
	log.Printf("Starting with %+v\n", cfg)
	cfg.Redis.Password = pword

	// create redis client
	rdb := redis.NewClient(ttl_redis.ToV8Options(&cfg.Redis))
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalln("Could not connect to redis:", err)
	}
	log.Println("Connectivity to redis established.")

	updater := set_ttls.NewUpdater(&cfg.Config, rdb)
	updater.SetTTLs(ctx)
}
