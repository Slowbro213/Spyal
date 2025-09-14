package cache

import (
	"fmt"
	"time"

	"spyal/cache/redis"
	"spyal/cache/valkey"
)

type Factory func(string) (Driver, error)

func NewDriverFactory(name string) (Factory, error) {
	switch name {
	case "redis":
		return func(url string) (Driver, error) {
			return redis.NewClient(url)
		}, nil
	case "valkey":
		return func(addr string) (Driver, error) {
			return valkey.NewClient(addr)
		}, nil
	default:
		return nil, fmt.Errorf("unsupported cache driver: %s", name)
	}
}

func WithRetry(factory Factory, arg string, attempts int, delay time.Duration) (Driver, error) {
	var d Driver
	var err error
	for i := range attempts {
		d, err = factory(arg)
		if err == nil {
			return d, nil
		}
		if i < attempts-1 {
			time.Sleep(delay)
		}
	}
	return nil, fmt.Errorf("failed after %d attempts: %w", attempts, err)
}
