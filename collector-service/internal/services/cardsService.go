package services

import (
	"net/http"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/models"
)

type CardsService struct {
	cardsRepository CardsRepositorer
}

type CardsRepositorer interface {
	GetCollection(collectionId string) (*models.Collection, *models.ResponseErr)
	AddCardToCollection(collectionId string, card *models.Card) *models.ResponseErr
	SetCardCountInCollection(collectionId string, card *models.Card) *models.ResponseErr
	DeleteCardFromCollection(collectionId string, card *models.Card) *models.ResponseErr
}

func NewCardsService(cardsRepository CardsRepositorer) *CardsService {
	return &CardsService{
		cardsRepository: cardsRepository,
	}
}

// ListCardsInCollection retrieves all cards in a collection by its ID.
func (cs CardsService) ListCardsInCollection(collectionId string) ([]*models.Card, *models.ResponseErr) {
	collection, err := cs.cardsRepository.GetCollection(collectionId)
	if err != nil {
		return nil, err
	}

	if collection == nil {
		return nil, &models.ResponseErr{
			Status:  http.StatusNotFound,
			Message: "Collection not found",
		}
	}

	if collection.Cards == nil {
		collection.Cards = []*models.Card{}
	}

	return collection.Cards, nil
}

// AddCardToCollection adds a card to a collection by its ID.
func (cs CardsService) AddCardToCollection(collectionId string, card *models.Card) *models.ResponseErr {
	return cs.cardsRepository.AddCardToCollection(collectionId, card)
}

// SetCardCountInCollection updates the count of a card in a collection by its ID.
func (cs CardsService) SetCardCountInCollection(collectionId string, card *models.Card) *models.ResponseErr {
	return cs.cardsRepository.SetCardCountInCollection(collectionId, card)
}

// DeleteCardFromCollection removes a card from a collection by its ID.
func (cs CardsService) DeleteCardFromCollection(collectionId string, card *models.Card) *models.ResponseErr {
	return cs.cardsRepository.DeleteCardFromCollection(collectionId, card)
}
