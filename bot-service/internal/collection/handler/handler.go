package handler

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/session"
	"github.com/ShenokZlob/collector-ouphe/pkg/contracts/collections"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/fsm"
	"github.com/go-telegram/ui/keyboard/inline"
)

// type inputState string

// const (
// 	stateCreateCollectionCommand inputState = "waiting_collection_name"
// 	stateRenameCollectionCommand inputState = "waiting_collection_rename"
// 	stateDeleteCollectionCommand inputState = "waiting_collection_delete"
// )

type CollectionHandler struct {
	log     logger.Logger
	usecase CollectionUsecase
}

type CollectionUsecase interface {
	GetCollecionsList(ctx context.Context) ([]string, error)
	CreateaCollection(ctx context.Context, name string) (*collections.Collection, error)
	RenameCollection(ctx context.Context, oldName, newName string) error
	DeleteCollection(ctx context.Context, name string) error
}

func NewCollectionHandler(log logger.Logger, usecase CollectionUsecase) *CollectionHandler {
	return &CollectionHandler{
		log:     log,
		usecase: usecase,
	}
}

// HandleCreateCollection create a new collection for user
func (h *CollectionHandler) CreateCollectionCommand(f *fsm.FSM) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		h.log.Info("CreateCollectionCommand executing")

		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID

		currentState := f.Current(userID)
		if currentState != session.StateDefault {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Вы уже находитесь в другом процессе. Пожалуйста, завершите текущую операцию или отмените её командой /cancel.",
			})
			return
		}

		f.Transition(userID, session.StateAskCreateCollection, b, chatID)
	}
}

func (h *CollectionHandler) CallbackAskCreateCollection(f *fsm.FSM) fsm.Callback {
	return func(f *fsm.FSM, args ...any) {
		h.log.Info("CallbackAskCreateCollection executing")
		b := args[0].(*bot.Bot)
		chatID := args[1]

		b.SendMessage(context.Background(), &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Отправьте название для новой коллекции\nДля отмены операции введите /cancel",
		})
	}
}

func (h *CollectionHandler) CallbackCreateCollection(f *fsm.FSM) fsm.Callback {
	return func(f *fsm.FSM, args ...any) {
		h.log.Info("CallbackCreateCollection executing")
		ctx := args[4].(context.Context)

		b := args[0].(*bot.Bot)
		chatID := args[1]
		collectionName := args[2].(string)
		userID := args[3].(int64)

		if !validateCollectionName(collectionName) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Слишком длинное название! Используйте не более 20 символов.",
			})
			f.Transition(userID, session.StateDefault)
			return
		}

		_, err := h.usecase.CreateaCollection(ctx, collectionName)
		if err != nil {
			h.log.Error("Failed to create collection", logger.Error(err), logger.String("collection_name", collectionName))
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Не получилось создать коллекцию :(",
			})
			f.Transition(userID, session.StateDefault)
			return
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Коллекция создана!",
		})
		f.Transition(userID, session.StateDefault)
	}
}

// HandleGetCollectionsList show user's list of collections
// Only info about collections without cards in.
func (h *CollectionHandler) GetCollectionsListCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.log.Info("GetCollectionsListCommand executing")

	collections, err := h.usecase.GetCollecionsList(ctx)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Не получилось получить список коллекций :(",
		})
		return
	}

	kb := inline.New(b)
	for _, v := range collections {
		kb = kb.Row()
		kb = kb.Button(v, []byte(v), onInlineKeyboardSelect)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        formatCollectionList(collections),
		ReplyMarkup: kb,
	})
}

func formatCollectionList(collections []string) string {
	if len(collections) == 0 {
		return "У вас нет коллекций."
	}

	result := "Ваши коллекции:\n"
	for _, collection := range collections {
		result += "- " + collection + "\n"
	}
	return result
}

func onInlineKeyboardSelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	// TODO: replace this func
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Message.Chat.ID,
		Text:   "You selected: " + string(data) + "\n REPLACE ME :)",
	})
}

// HandleRenameCollection rename a collection
func (h *CollectionHandler) RenameCollectionCommand(f *fsm.FSM) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		h.log.Info("RenameCollectionCommand executing")

		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID

		currentState := f.Current(userID)
		if currentState != session.StateDefault {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Вы уже находитесь в другом процессе. Пожалуйста, завершите текущую операцию или отмените её командой /cancel.",
			})
			return
		}

		f.Transition(userID, session.StateAskRenameCollection, b, chatID)
	}
}

func (h *CollectionHandler) CallbackAskRenameCollection(f *fsm.FSM) fsm.Callback {
	return func(f *fsm.FSM, args ...any) {
		b := args[0].(*bot.Bot)
		chatID := args[1]

		b.SendMessage(context.Background(), &bot.SendMessageParams{
			ChatID: chatID,
			Text: "Введите названия старой и новой коллекции через пробел.\n" +
				"Например: `СтараяКоллекция НоваяКоллекция`\n" +
				"Для отмены операции введите /cancel",
		})
	}
}

func (h *CollectionHandler) CallbackRenameCollection(f *fsm.FSM) fsm.Callback {
	return func(f *fsm.FSM, args ...any) {
		ctx := args[5].(context.Context)

		b := args[0].(*bot.Bot)
		chatID := args[1]
		oldName := args[2].(string)
		newName := args[3].(string)
		userID := args[4].(int64)

		// TODO: change this validation
		if !validateCollectionName(oldName) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "У вас нет такой коллекции!",
			})
			f.Transition(userID, session.StateDefault)
			return
		}

		if !validateCollectionName(newName) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Слишком длинное название! Используйте не более 20 символов.",
			})
			f.Transition(userID, session.StateDefault)
			return
		}

		err := h.usecase.RenameCollection(ctx, oldName, newName)
		if err != nil {
			h.log.Error("Failed to rename collection", logger.Error(err), logger.String("old_name", oldName), logger.String("new_name", newName))
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Не получилось создать коллекцию :(",
			})
			f.Transition(userID, session.StateDefault)
			return
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Коллекция переименована!",
		})
		f.Transition(userID, session.StateDefault)
	}
}

// HandleDeleteCollection delete a collection
func (h *CollectionHandler) DeleteCollectionCommand(f *fsm.FSM) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		h.log.Info("DeleteCollectionCommand executing")

		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID

		currentState := f.Current(userID)
		if currentState != session.StateDefault {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Вы уже находитесь в другом процессе. Пожалуйста, завершите текущую операцию или отмените её командой /cancel.",
			})
			return
		}

		f.Transition(userID, session.StateAskDeleteCollection, b, chatID)
	}
}

func (h *CollectionHandler) CallbackAskDeleteCollection(f *fsm.FSM) fsm.Callback {
	return func(f *fsm.FSM, args ...any) {
		b := args[0].(*bot.Bot)
		chatID := args[1]

		b.SendMessage(context.Background(), &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Введите название коллекции, которую хотите удалить.\nДля отмены операции введите /cancel",
		})
	}
}

func (h *CollectionHandler) CallbackDeleteCollection(f *fsm.FSM) fsm.Callback {
	return func(f *fsm.FSM, args ...any) {
		ctx := args[4].(context.Context)

		b := args[0].(*bot.Bot)
		chatID := args[1]
		collectionName := args[2].(string)
		userID := args[3].(int64)

		// TODO: change this validation
		if !validateCollectionName(collectionName) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "У вас нет такой коллекции!",
			})
			f.Transition(userID, session.StateDefault)
			return
		}

		err := h.usecase.DeleteCollection(ctx, collectionName)
		if err != nil {
			h.log.Error("Failed to delete collection", logger.Error(err), logger.String("collection_name", collectionName))
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   "Не получилось удалить коллекцию :(",
			})
			f.Transition(userID, session.StateDefault)
			return
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   fmt.Sprintf("Коллекция %s удалена.", collectionName),
		})
		f.Transition(userID, session.StateDefault)
	}
}

// Correct name
// 20 or less symb (utf-8)
func validateCollectionName(name string) bool {
	return utf8.RuneCountInString(name) <= 20
}
