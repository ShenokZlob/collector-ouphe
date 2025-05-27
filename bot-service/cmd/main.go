package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	appbot "github.com/ShenokZlob/collector-ouphe/bot-service/internal/app/bot"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log, err := logger.NewZapLogger(false)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	log.Info("Init Redis...")
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	log.Info("Init bot...")
	app, err := appbot.NewAppBot(os.Getenv("BOT_TOKEN"), os.Getenv("COLLECTOR_URL"), log, redisClient)
	if err != nil {
		log.Error("Failed to create app bot", logger.Error(err))
	}

	log.Info("Runing app...")
	app.Run(ctx)
}
