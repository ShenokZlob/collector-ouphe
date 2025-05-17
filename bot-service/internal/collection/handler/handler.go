package handler

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/state"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

type inputState string

const (
	stateCreateCollectionCommand inputState = "waiting_collection_name"
	stateRenameCollectionCommand inputState = "waiting_collection_rename"
	stateDeleteCollectionCommand inputState = "waiting_collection_delete"
)

type CollectionHandler struct {
	usecase CollectionUsecase
	mgr     state.Manager
	log     logger.Logger
}

type CollectionUsecase interface {
	GetCollecionsList(ctx context.Context) ([]string, error)
	CreateaCollection(ctx context.Context, name string) error
	RenameCollection(ctx context.Context, oldName, newName string) error
	DeleteCollection(ctx context.Context, name string) error
}

func NewCollectionHandler(usecase CollectionUsecase, mgr state.Manager, log logger.Logger) *CollectionHandler {
	return &CollectionHandler{
		usecase: usecase,
		mgr:     mgr,
		log:     log,
	}
}

// HandleCreateCollection create a new collection for user
func (h *CollectionHandler) CreateCollectionCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.log.Info("CreateCollectionCommand executing")

	ctx = state.WithState(ctx, string(stateCreateCollectionCommand))

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Отправьте название для новой коллекции\nДля отмены операции введите /cancel",
	})
}

func (h *CollectionHandler) CreateCollectionResponse(ctx context.Context, b *bot.Bot, update *models.Update) {
	st, _ := state.GetState(ctx)
	h.log.Info("CreateCollectionResponse executing", logger.String("state", st))
	if st != string(stateCreateCollectionCommand) {
		return
	}

	// Check collection's name
	collectionName := update.Message.Text
	if !validateCollectionName(collectionName) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Слишком длинное название! Используйте не более 20 символов.",
		})
		return
	}

	err := h.usecase.CreateaCollection(ctx, collectionName)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Не получилось создать коллекцию :(",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Коллекция создана!",
	})
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
		Text:        "Ваши коллекции карт.",
		ReplyMarkup: kb,
	})
}

// TODO: replace this
func onInlineKeyboardSelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: mes.Message.Chat.ID,
		Text:   "You selected: " + string(data) + "\n REPLACE ME :)",
	})
}

// HandleRenameCollection rename a collection
func (h *CollectionHandler) RenameCollectionCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.log.Info("RenameCollectionCommand executing")

	ctx = state.WithState(ctx, string(stateRenameCollectionCommand))

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Отправьте через пробел название старое название коллекции и новое\nДля отмены операции введите /cancel",
	})
}

func (h *CollectionHandler) RenameCollectionResponse(ctx context.Context, b *bot.Bot, update *models.Update) {
	st, _ := state.GetState(ctx)
	h.log.Info("RenameCollectionResponse executing", logger.String("state", st))
	if st != string(stateRenameCollectionCommand) {
		return
	}

	names := strings.Split(update.Message.Text, " ")

	// TODO: Check old collection's name

	// Check new collection's name
	if !validateCollectionName(names[1]) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Слишком длинное название! Используйте не более 20 символов.",
		})
		return
	}

	err := h.usecase.RenameCollection(ctx, names[0], names[1])
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Не получилось переименовать коллекцию :(",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Коллекция переименована.",
	})
}

// HandleDeleteCollection delete a collection
func (h *CollectionHandler) DeleteCollectionCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.log.Info("DeleteCollectionCommand executing")

	ctx = state.WithState(ctx, string(stateDeleteCollectionCommand))

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Введите название коллекции, которую хотите удалить.\nДля отмены операции введите /cancel",
	})
}

func (h *CollectionHandler) DeleteCollectionResponse(ctx context.Context, b *bot.Bot, update *models.Update) {
	st, _ := state.GetState(ctx)
	h.log.Info("DeleteCollectionResponse executing", logger.String("state", st))
	if st != string(stateDeleteCollectionCommand) {
		return
	}

	collectionName := update.Message.Text
	if !validateCollectionName(collectionName) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Слишком длинное название! Используйте не более 20 символов.",
		})
		return
	}

	err := h.usecase.DeleteCollection(ctx, collectionName)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Не получилось удалить коллекцию :(",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Коллекция %s удаленна.", collectionName),
	})
}

// Correct name
// 20 or less symb (utf-8)
func validateCollectionName(name string) bool {
	return utf8.RuneCountInString(name) <= 20
}
