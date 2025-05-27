package session

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	redis *redis.Client
}

const (
	ttlCache    = time.Duration(15 * time.Minute)
	prefixCache = "cache:"
)

func NewCache(redis *redis.Client) *Cache {
	return &Cache{
		redis: redis,
	}
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}) error {
	return c.redis.Set(ctx, prefixCache+key, value, ttlCache).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.redis.Get(ctx, prefixCache+key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
