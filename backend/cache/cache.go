package cache

import (
	"context"
	"fmt"
	"os"
	"time"
)

type Driver interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

const (
	maxAttempts = 5
	attemptInterval = 3
)

//nolint
var driver Driver

func Init() {
	drv := os.Getenv("CACHE_DRIVER")
	factory, err := NewDriverFactory(drv)
	if err != nil  {
		panic("CACHE_DRIVER not set or invalid. Must be 'redis' or 'valkey'")
	}

	arg := ""
	switch drv {
	case "redis":
		arg = os.Getenv("REDIS_URL")
		if arg == "" {
			panic("REDIS_URL not set for redis cache driver")
		}
	case "valkey":
		arg = os.Getenv("VALKEY_URL")
		if arg == "" {
			panic("VALKEY_URL not set for valkey cache driver")
		}
	}

	d, err := WithRetry(factory, arg, maxAttempts, attemptInterval*time.Second)
	if err != nil {
		panic(fmt.Errorf("failed to initialize %s cache: %w", drv, err))
	}
	driver = d
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
