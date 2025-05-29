package app

import (
	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/session"
	"github.com/redis/go-redis/v9"
)

func InitCache(client *redis.Client) *session.Cache {
	return session.NewCache(client)
}
