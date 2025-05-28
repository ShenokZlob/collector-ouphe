package session

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Manager interface {
	SetState(ctx context.Context, telegramID int64, state string) error
	GetState(ctx context.Context, telegramID int64) (string, error)
	ClearState(ctx context.Context, telegramID int64) error
}

type stateManager struct {
	redis *redis.Client
}

const (
	ttlState    = time.Duration(3 * time.Minute)
	prefixState = "state:"
)

func NewStateRedis(redis *redis.Client) Manager {
	return &stateManager{
		redis: redis,
	}
}

func (m *stateManager) SetState(ctx context.Context, telegramID int64, state string) error {
	tgIDstring := strconv.FormatInt(telegramID, 10)
	return m.redis.Set(ctx, prefixState+tgIDstring, state, ttlState).Err()
}

func (m *stateManager) GetState(ctx context.Context, telegramID int64) (string, error) {
	tgIDstring := strconv.FormatInt(telegramID, 10)
	val, err := m.redis.Get(ctx, prefixState+tgIDstring).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // State not found
		}
		return "", err // Other error
	}
	return val, nil
}

func (m *stateManager) ClearState(ctx context.Context, telegramID int64) error {
	tgIDstring := strconv.FormatInt(telegramID, 10)
	return m.redis.Del(ctx, prefixState+tgIDstring).Err()
}
