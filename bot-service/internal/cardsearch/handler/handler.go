package handler

import (
	"context"
	"strings"

	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type CardSearchHandler struct {
	log               logger.Logger
	cardSearchUsecase CardSearchUsecase
}

type CardSearchUsecase interface {
	SearchCard(context context.Context, cardName string) ([]scryfall.Card, error)
}

func NewCardSearchHandler(log logger.Logger, usecase CardSearchUsecase) *CardSearchHandler {
	return &CardSearchHandler{
		log:               log,
		cardSearchUsecase: usecase,
	}
}

func (h *CardSearchHandler) HandleSearchCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.log.Info("HandleSearchCommand executing")

	textSplit := strings.Split(update.Message.Text, " ")
	if len(textSplit) < 2 {
		h.log.Warn("Search command without card name", logger.String("command", update.Message.Text))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Пожалуйста, укажите название карты после команды /search.",
		})
		return
	}

	cardName := strings.Join(textSplit[1:], " ")
	if cardName == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Пожалуйста, укажите название карты после команды /search.",
		})
		return
	}

	cards, err := h.cardSearchUsecase.SearchCard(ctx, cardName)
	if err != nil {
		h.log.Error("Error searching for card", logger.Error(err))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Произошла ошибка при поиске карты. Пожалуйста, попробуйте позже.",
		})
		return
	}

	h.log.Info("Cards found", logger.Int("count", len(cards)))
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Найденные карты:\n" + formatCards(cards),
	})

}

func formatCards(cards []scryfall.Card) string {
	var sb strings.Builder
	for i, card := range cards {
		if i > 9 {
			break
		}
		sb.WriteString(card.Name + "\t")
		if card.ImageURIs != nil {
			sb.WriteString("Image: " + card.ImageURIs.Normal)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
