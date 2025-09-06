package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"spyal/cache/redis"
	"spyal/cache/valkey"
)

type Driver interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

var driver Driver

func Init() {
	drv := os.Getenv("CACHE_DRIVER")

	switch drv {
	case "redis":
		url := os.Getenv("REDIS_URL")
		if url == "" {
			panic("REDIS_URL not set for redis cache driver")
		}
		r, err := redis.NewClient("redis://" + url)
		if err != nil {
			panic(fmt.Errorf("failed to initialize redis cache: %w", err))
		}
		driver = r
	case "valkey":
		vURL := os.Getenv("VALKEY_URL")
		if vURL == "" {
			panic("VALKEY_URL not set for valkey cache driver")
		}
		addrs := []string{vURL}
		v, err := valkey.NewClient(addrs, false)
		if err != nil {
			panic(fmt.Errorf("failed to initialize valkey cache: %w", err))
		}
		driver = v
	default:
		panic("CACHE_DRIVER not set or invalid. Must be 'redis' or 'valkey'")
	}
}

func Get(ctx context.Context, key string) (string, error) {
	return driver.Get(ctx, key)
}

func Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return driver.Set(ctx, key, value, ttl)
}

func Delete(ctx context.Context, key string) error {
	return driver.Delete(ctx, key)
}

func Remember(ctx context.Context, key string, ttl time.Duration, callback func() (string, error)) (string, error) {
	val, err := Get(ctx, key)
	if err == nil && val != "" {
		return val, nil
	}

	val, err = callback()
	if err != nil {
		return "", err
	}

	if err := Set(ctx, key, val, ttl); err != nil {
		return "", err
	}

	return val, nil
}
