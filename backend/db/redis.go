package db

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

//nolint:gochecknoglobals
var (
	redisInstance *redis.Client
	redisOnce     sync.Once
)

// GetRedis returns the singleton Redis client.
func GetRedis() (*redis.Client, error) {
	var err error
	redisOnce.Do(func() {
		addr := os.Getenv("REDIS_ADDR")
		if addr == "" {
			err = errors.New("REDIS_ADDR not set")
			return
		}

		redisInstance = redis.NewClient(&redis.Options{
			Addr: addr,
		})

		// Test the connection
		_, err = redisInstance.Ping(context.Background()).Result()
	})

	return redisInstance, err
}
