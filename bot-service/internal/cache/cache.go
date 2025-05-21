package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	redis *redis.Client
}

func NewCache(redis *redis.Client) *Cache {
	return &Cache{
		redis: redis,
	}
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}) error {
	return c.redis.Set(ctx, key, value, 0).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
