package appbot

import (
	"context"

	authHandler "github.com/ShenokZlob/collector-ouphe/bot-service/internal/auth/handler"
	authUsecase "github.com/ShenokZlob/collector-ouphe/bot-service/internal/auth/usecase"
	"github.com/go-telegram/bot"
)

type AppBot struct {
	bot          *bot.Bot
	collectorURL string
}

func NewAppBot(token string, collectorURL string) (*AppBot, error) {
	// Auth
	authUsecase := authUsecase.NewAuthUsecase(collectorURL)
	authHandler := authHandler.NewAuthHandler(authUsecase)

	// Bot options
	opts := []bot.Option{
		bot.WithMiddlewares(authHandler.RegistrationMiddleware),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, err
	}

	// Initialize router
	b.RegisterHandler(bot.HandlerTypeMessageText, "/register", bot.MatchTypeExact, authHandler.HandleRegister)

	return &AppBot{
		bot:          b,
		collectorURL: collectorURL,
	}, nil
}

func (a *AppBot) Run(ctx context.Context) {
	a.bot.Start(ctx)
}
