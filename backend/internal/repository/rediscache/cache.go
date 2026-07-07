package rediscache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type cacheRepo struct {
	client *redis.Client
}

func NewCacheRepository(c *redis.Client) *cacheRepo {
	return &cacheRepo{client: c}
}

func (r *cacheRepo) Get(ctx context.Context, key string) ([]byte, error) {
	var value []byte

	err := r.client.Get(ctx, key).Scan(&value)
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return value, nil
}

func (r *cacheRepo) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}
