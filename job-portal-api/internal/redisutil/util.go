package redisutil

import (
	"context"
	"fmt"
	"job-portal-api/config"

	"github.com/go-redis/redis/v8"
)

func Redis(config config.RedisConfig) (*redis.Client, error) {
	fmt.Print("===========================", config.Host, config.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.Host, config.Port),

		Password: config.Password,
		DB:       config.DB,
	})

	// Use Ping to check if the Redis connection is successful
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return rdb, nil
}
