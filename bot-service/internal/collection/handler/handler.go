package handler

import (
	"context"

	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type CollectionHandler struct {
	usecase CollectionUsecase
	log     logger.Logger
}

type CollectionUsecase interface {
}

func NewCollectionHandler(usecase CollectionUsecase, log logger.Logger) *CollectionHandler {
	return &CollectionHandler{
		usecase: usecase,
		log:     log,
	}
}

// HandleCreateCollection create a new collection for user
func (h *CollectionHandler) HandleCreateCollection(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Отправте название для новой коллекции\nДля отмены операции введите /cancel",
	})

}

// HandleGetCollectionsList show user's list of collections
// Only info about collections without cards in.
func (h *CollectionHandler) HandleGetCollectionsList(ctx context.Context, b *bot.Bot, update *models.Update) {
}

// HandleRenameCollection rename a collection
func (h *CollectionHandler) HandleRenameCollection(ctx context.Context, b *bot.Bot, update *models.Update) {
}

// HandleDeleteCollection delete a collection
func (h *CollectionHandler) HandleDeleteCollection(ctx context.Context, b *bot.Bot, update *models.Update) {
}
