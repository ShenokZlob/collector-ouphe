package appbot

import (
	"context"
	"strings"

	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/app"
	authHandler "github.com/ShenokZlob/collector-ouphe/bot-service/internal/auth/handler"
	authUsecase "github.com/ShenokZlob/collector-ouphe/bot-service/internal/auth/usecase"
	cardsearchHandler "github.com/ShenokZlob/collector-ouphe/bot-service/internal/cardsearch/handler"
	cardsearchUsecase "github.com/ShenokZlob/collector-ouphe/bot-service/internal/cardsearch/usecase"
	collectionHandler "github.com/ShenokZlob/collector-ouphe/bot-service/internal/collection/handler"
	collectionUsecase "github.com/ShenokZlob/collector-ouphe/bot-service/internal/collection/usecase"
	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/session"
	"github.com/ShenokZlob/collector-ouphe/pkg/collectorclient"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/fsm"
	"github.com/redis/go-redis/v9"
)

type AppBot struct {
	log logger.Logger
	b   *bot.Bot
	f   *fsm.FSM
}

func NewAppBot(token string, collectorURL string, log logger.Logger, redisClient *redis.Client) (*AppBot, error) {
	appbot := &AppBot{}

	// Initialize other dependencies
	appbot.log = log
	cache := app.InitCache(redisClient)
	collectorClient := collectorclient.NewHTTPCollectorClient(collectorURL, log)

	// Auth
	authUse := authUsecase.NewAuthUsecase(log, collectorClient, cache)
	authHand := authHandler.NewAuthHandler(authUse, log)

	// Collection
	collUse := collectionUsecase.NewCollectionUsecaseImpl(log, collectorClient)
	collHand := collectionHandler.NewCollectionHandler(log, collUse)

	// Card Search
	csUse := cardsearchUsecase.NewCardSearchUsecaseImpl(log)
	csHand := cardsearchHandler.NewCardSearchHandler(log, csUse)

	// Init FSM and its callbacks
	appbot.log.Info("Initializing FSM")
	appbot.f = fsm.New(session.StateDefault, map[fsm.StateID]fsm.Callback{})
	appbot.f.AddCallbacks(map[fsm.StateID]fsm.Callback{
		session.StateAskCreateCollection: collHand.CallbackAskCreateCollection(appbot.f),
		session.StateCreateCollection:    collHand.CallbackCreateCollection(appbot.f),
		session.StateAskRenameCollection: collHand.CallbackAskRenameCollection(appbot.f),
		session.StateRenameCollection:    collHand.CallbackRenameCollection(appbot.f),
		session.StateAskDeleteCollection: collHand.CallbackAskDeleteCollection(appbot.f),
		session.StateDeleteCollection:    collHand.CallbackDeleteCollection(appbot.f),
	})

	// Bot options
	opts := []bot.Option{
		bot.WithMiddlewares(authHand.RegistrationMiddleware),
		bot.WithDefaultHandler(appbot.defaultHandler),
		bot.WithMessageTextHandler("/cancel", bot.MatchTypeExact, appbot.handlerCancel),
	}

	// Initialize bot
	var err error
	appbot.b, err = bot.New(token, opts...)
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
	_, err = appbot.b.SetMyCommands(context.TODO(), &bot.SetMyCommandsParams{
		Commands: commands,
	})
	if err != nil {
		log.Error("failed to set command", logger.Error(err))
		return nil, err
	}

	// Initialize router

	// Auth
	appbot.b.RegisterHandler(bot.HandlerTypeMessageText, "/register", bot.MatchTypeExact, authHand.HandleRegister)

	// Collection
	appbot.b.RegisterHandler(bot.HandlerTypeMessageText, "collections", bot.MatchTypeCommand, collHand.GetCollectionsListCommand)
	appbot.b.RegisterHandler(bot.HandlerTypeMessageText, "collection_new", bot.MatchTypeCommand, collHand.CreateCollectionCommand(appbot.f))
	appbot.b.RegisterHandler(bot.HandlerTypeMessageText, "collection_rename", bot.MatchTypeCommand, collHand.RenameCollectionCommand(appbot.f))
	appbot.b.RegisterHandler(bot.HandlerTypeMessageText, "collection_delete", bot.MatchTypeCommand, collHand.DeleteCollectionCommand(appbot.f))

	// Card Search
	appbot.b.RegisterHandler(bot.HandlerTypeMessageText, "search", bot.MatchTypeCommand, csHand.HandleSearchCommand)

	return appbot, nil
}

func (ab *AppBot) Run(ctx context.Context) {
	ab.b.Start(ctx)
}

func (ab *AppBot) defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	currentState := ab.f.Current(userID)

	switch currentState {
	case session.StateDefault:
		ab.b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Используйте /help для получения списка команд.",
		})
		return

	case session.StateAskCreateCollection:
		collectionName := update.Message.Text
		ab.f.Transition(userID, session.StateCreateCollection, b, chatID, collectionName, userID, ctx)
		return

	case session.StateAskRenameCollection:
		resp := update.Message.Text
		collectionsNames := strings.Split(resp, " ")
		if len(collectionsNames) != 2 {
			ab.b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Пожалуйста, введите старое и новое название коллекции через пробел.",
			})
		}
		ab.f.Transition(userID, session.StateRenameCollection, b, chatID, collectionsNames[0], collectionsNames[1], userID, ctx)
		return

	case session.StateAskDeleteCollection:
		collectionName := update.Message.Text
		ab.f.Transition(userID, session.StateDeleteCollection, b, chatID, collectionName, userID, ctx)
		return

	default:
		ab.log.Warn("unexpected state ", logger.String("state", string(currentState)))
	}

}

func (ab *AppBot) handlerCancel(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	currentState := ab.f.Current(userID)

	if currentState == session.StateDefault {
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Canceled",
	})

	ab.f.Transition(userID, session.StateDefault)
}

func (ap *AppBot) callbackFinish(f *fsm.FSM, args ...any) {
	chatID := args[0]
	userID := args[1].(int64)

	ap.b.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Успешный успех!",
	})

	f.Transition(userID, session.StateDefault)
}
