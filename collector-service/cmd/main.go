package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/app"
	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/config"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	log := logger.NewZapLogger(false)
	log.Info("Starting collector service")

	log.Info("Read config")
	cfg := config.InitConfig()

	log.Info("Init database")
	db := app.InitDataBase(cfg)

	log.Info("Init app server")
	appServer := app.InitServer(cfg, log, db)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	log.Info("Starting app")
	go appServer.Run()
	log.Info("App started")

	<-ctx.Done()

	log.Info("Stopping app")
	appServer.Stop(ctx)

	os.Exit(0)
}
