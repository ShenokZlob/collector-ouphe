package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/app"
	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/config"
	"github.com/joho/godotenv"

	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer logger.Sync()
	// logger.Info("Init logger")

	logger.Info("Init config")
	cfg := config.InitConfig()

	logger.Info("Init database")
	db := app.InitDataBase(cfg)

	logger.Info("Init app server")
	appServer := app.InitServer(cfg, logger, db)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	logger.Info("Starting app")
	go appServer.Run()
	logger.Info("App started")

	<-ctx.Done()

	logger.Info("Stopping app")
	appServer.Stop(ctx)

	os.Exit(0)
}
