package appbot

import (
	"context"

	authHandler "github.com/ShenokZlob/collector-ouphe/bot-service/internal/auth/handler"
	authUsecase "github.com/ShenokZlob/collector-ouphe/bot-service/internal/auth/usecase"
	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/cache"
	cardsearchHandler "github.com/ShenokZlob/collector-ouphe/bot-service/internal/cardsearch/handler"
	cardsearchUsecase "github.com/ShenokZlob/collector-ouphe/bot-service/internal/cardsearch/usecase"
	collectionHandler "github.com/ShenokZlob/collector-ouphe/bot-service/internal/collection/handler"
	collectionUsecase "github.com/ShenokZlob/collector-ouphe/bot-service/internal/collection/usecase"
	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/state"
	"github.com/ShenokZlob/collector-ouphe/pkg/collectorclient"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type AppBot struct {
	bot          *bot.Bot
	collectorURL string
	log          logger.Logger
}

func NewAppBot(token string, collectorURL string, log logger.Logger, cache *cache.Cache) (*AppBot, error) {
	// State - save user's states
	mgr := state.NewMemoryManager()

	// HTTPCollectorClient
	collectorClient := &collectorclient.HTTPCollectorClient{
		URL: collectorURL,
		Log: log,
	}

	// Auth
	authUse := authUsecase.NewAuthUsecase(log, collectorClient, cache)
	authHand := authHandler.NewAuthHandler(authUse, log)

	// Collection
	collUse := collectionUsecase.NewCollectionUsecaseImpl(log, collectorClient)
	collHand := collectionHandler.NewCollectionHandler(collUse, mgr, log)

	// Card Search
	csUse := cardsearchUsecase.NewCardSearchUsecaseImpl(log)
	csHand := cardsearchHandler.NewCardSearchHandler(log, csUse)

	// Bot options
	opts := []bot.Option{
		bot.WithMiddlewares(authHand.RegistrationMiddleware, state.Middleware(mgr)),
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, err
	}

	// Init commands panel
	commands := []models.BotCommand{
		{Command: "search", Description: "Search for a card /command <card name>"},
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

	// Cancel command
	b.RegisterHandler(bot.HandlerTypeMessageText, "/cancel", bot.MatchTypeExact, state.CancelHandler(mgr))

	// Auth
	b.RegisterHandler(bot.HandlerTypeMessageText, "/register", bot.MatchTypeExact, authHand.HandleRegister)

	// Collection
	b.RegisterHandler(bot.HandlerTypeMessageText, "/collections", bot.MatchTypeExact, collHand.GetCollectionsListCommand)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/collection_new", bot.MatchTypeExact, collHand.CreateCollectionCommand)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/collection_rename", bot.MatchTypeExact, collHand.RenameCollectionCommand)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/collection_delete", bot.MatchTypeExact, collHand.DeleteCollectionCommand)

	// Card Search
	b.RegisterHandler(bot.HandlerTypeMessageText, "/search", bot.MatchTypeExact, csHand.HandleSearchCommand)

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
