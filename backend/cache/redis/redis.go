// cache/redis/redis.go
package redis

import (
	"context"
	"time"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

func NewClient(url string) (*Client, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	r := redis.NewClient(opt)
	return &Client{client: r}, nil
}

func (r *Client) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *Client) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *Client) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
