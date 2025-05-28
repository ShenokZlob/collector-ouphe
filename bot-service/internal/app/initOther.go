package app

import (
	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/session"
	"github.com/redis/go-redis/v9"
)

func InitSessions(client *redis.Client) (*session.Cache, session.Manager) {
	cache := initCache(client)
	state := initState(client)

	return cache, state
}

func initCache(client *redis.Client) *session.Cache {
	return session.NewCache(client)
}

func initState(client *redis.Client) session.Manager {
	return session.NewStateRedis(client)
}
