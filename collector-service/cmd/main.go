package main

// @title           Collector Ouphe API
// @version         1.0
// @description     Сервис сбора и анализа данных Collector Ouphe
//
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
//
// @host      localhost:8080
// @BasePath  /

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ShenokZlob/collector-ouphe/collector-service/docs"
	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/app"
	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/config"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	log, err := logger.NewZapLogger(false)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

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
