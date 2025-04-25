package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	appbot "github.com/ShenokZlob/collector-ouphe/bot-service/internal/app/bot"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Println("Initializing app...")
	app, err := appbot.NewAppBot(os.Getenv("BOT_TOKEN"), os.Getenv("COLLECTOR_URL"))
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	log.Println("Initializing router...")
	appbot.InitRouter(app)

	log.Println("Runing app...")
	app.Run(ctx)
}
