package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/config"
	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/servers"

	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer logger.Sync()
	logger.Info("Init logger")

	cfg := config.InitConfig()
	logger.Info("Init config")

	db := servers.InitDataBase(cfg)
	logger.Info("Init database")

	appServer := servers.InitServer(cfg, logger, db)
	logger.Info("Starting app server")
	go appServer.Run()

	<-ctx.Done()

	logger.Info("Stopping app")
	appServer.Stop(ctx)

	os.Exit(0)
}
