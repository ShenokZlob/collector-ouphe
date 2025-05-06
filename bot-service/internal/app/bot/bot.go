package appbot

import (
	"context"

	authHandler "github.com/ShenokZlob/collector-ouphe/bot-service/internal/auth/handler"
	authUsecase "github.com/ShenokZlob/collector-ouphe/bot-service/internal/auth/usecase"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type AppBot struct {
	bot          *bot.Bot
	collectorURL string
	log          logger.Logger
}

func NewAppBot(token string, collectorURL string, log logger.Logger) (*AppBot, error) {
	// Auth
	authUsecase := authUsecase.NewAuthUsecase(collectorURL, log)
	authHandler := authHandler.NewAuthHandler(authUsecase, log)

	// Bot options
	opts := []bot.Option{
		// bot.WithMiddlewares(authHandler.RegistrationMiddleware),
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, err
	}

	// Init commands panel
	commands := []models.BotCommand{
		{Command: "collections", Description: "View your collection's list"},
		{Command: "collection_new", Description: "Create new collection /command <name>"},
		{Command: "collection_rename", Description: "Rename collection /command <old name> <new name>"},
		{Command: "collection_delete", Description: "Delete collection /command <name>"},
		{Command: "register", Description: "Register your account"},
		{Command: "help", Description: "Help"},
	}
	_, err = b.SetMyCommands(context.TODO(), &bot.SetMyCommandsParams{
		Commands: commands,
	})
	if err != nil {
		log.Error("failed to set command", logger.Error(err))
		return nil, err
	}

	// Initialize router
	// Auth
	b.RegisterHandler(bot.HandlerTypeMessageText, "/register", bot.MatchTypeExact, authHandler.HandleRegister)

	return &AppBot{
		bot:          b,
		collectorURL: collectorURL,
		log:          log,
	}, nil
}

func (a *AppBot) Run(ctx context.Context) {
	a.bot.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Choose a command",
	})
}
