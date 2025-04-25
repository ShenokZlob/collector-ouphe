package appbot

import (
	authHandler "github.com/ShenokZlob/collector-ouphe/bot-service/internal/auth/handler"
	authUsecase "github.com/ShenokZlob/collector-ouphe/bot-service/internal/auth/usecase"
	"github.com/go-telegram/bot"
)

func InitRouter(app *AppBot) {
	// Auth
	authUC := authUsecase.NewAuthUsecase(app.collectorURL)
	authHandler := authHandler.NewAuthHandler(authUC)

	// Register handlers
	app.bot.RegisterHandler(bot.HandlerTypeMessageText, "/register", bot.MatchTypeExact, authHandler.HandleRegister)
}
