package valkey

import (
	"context"
	"errors"
	"time"

	"github.com/valkey-io/valkey-go"
)

type Cache struct {
	client valkey.Client
}

func NewClient(addres string) (*Cache, error) {
	addrs := []string{addres}
	c, err := valkey.NewClient(valkey.ClientOption{InitAddress: addrs, DisableCache: false})
	if err != nil {
		return nil, err
	}
	return &Cache{client: c}, nil
}

func (h *Cache) Close() {
	h.client.Close()
}

func (h *Cache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	if err := h.client.Do(ctx, h.client.B().Set().Key(key).Value(value).Build()).Error(); err != nil {
		return err
	}
	if ttl > 0 {
		if err := h.client.Do(ctx, h.client.B().Expire().Key(key).Seconds(int64(ttl.Seconds())).Build()).Error(); err != nil {
			return err
		}
	}
	return nil
}

func (h *Cache) Get(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", errors.New("key cannot be empty")
	}
	val, err := h.client.Do(ctx, h.client.B().Get().Key(key).Build()).ToString()
	return val, err
}

func (h *Cache) Delete(ctx context.Context, key string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	_, err := h.client.Do(ctx, h.client.B().Del().Key(key).Build()).AsInt64()
	return err
}
