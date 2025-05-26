package usecase

import (
	"context"
	"fmt"

	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
)

type cardSearchUsecaseImpl struct {
	log            logger.Logger
	scryfallClient *scryfall.Client
}

func NewCardSearchUsecaseImpl(log logger.Logger) *cardSearchUsecaseImpl {
	scryClient, err := scryfall.NewClient()
	if err != nil {
		log.Error("Failed to create Scryfall client", logger.Error(err))
		return nil
	}

	return &cardSearchUsecaseImpl{
		scryfallClient: scryClient,
		log:            log,
	}
}

func (c *cardSearchUsecaseImpl) SearchCard(ctx context.Context, cardName string) ([]scryfall.Card, error) {
	c.log.Info("Searching for card", logger.String("cardName", cardName))

	// Search for the card using Scryfall API
	resp, err := c.scryfallClient.SearchCards(ctx, cardName, scryfall.SearchCardsOptions{})
	if err != nil {
		c.log.Error("Error searching for card", logger.Error(err))
		return nil, err
	}

	if len(resp.Cards) == 0 {
		c.log.Info("No cards found", logger.String("cardName", cardName))
		return nil, fmt.Errorf("no cards found for name: %s", cardName)
	}

	c.log.Info("Cards found", logger.Int("count", len(resp.Cards)))
	return resp.Cards, nil

}
