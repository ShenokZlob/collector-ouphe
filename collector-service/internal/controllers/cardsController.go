package controllers

import (
	"net/http"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/models"
	"github.com/gin-gonic/gin"
)

type CardsController struct {
	cardsService CardsServicer
}

type CardsServicer interface {
	ListCardsInCollection(collectionId string) ([]*models.Card, *models.ResponseErr)
	AddCardToCollection(collectionId string, card *models.Card) *models.ResponseErr
	SetCardCountInCollection(collectionId string, card *models.Card) *models.ResponseErr
	DeleteCardFromCollection(collectionId string, card *models.Card) *models.ResponseErr
}

func NewCardsController(cardsService CardsServicer) *CardsController {
	return &CardsController{
		cardsService: cardsService,
	}
}

func (cc CardsController) ListCardsInCollection(ctx *gin.Context) {
	collectionId := ctx.Param("id")
	cards, respErr := cc.cardsService.ListCardsInCollection(collectionId)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ctx.JSON(200, cards)
}

func (cc CardsController) AddCardToCollection(ctx *gin.Context) {
	collectionId := ctx.Param("id")
	var card models.Card
	if err := ctx.ShouldBindJSON(&card); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	respErr := cc.cardsService.AddCardToCollection(collectionId, &card)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (cc CardsController) SetCardCountInCollection(ctx *gin.Context) {
	collectionId := ctx.Param("id")
	scryfallId := ctx.Param("card_id")
	var card models.Card
	if err := ctx.ShouldBindJSON(&card); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	card.ScryfallID = scryfallId
	respErr := cc.cardsService.SetCardCountInCollection(collectionId, &card)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (cc CardsController) DeleteCardFromCollection(ctx *gin.Context) {
	collectionId := ctx.Param("id")
	scryfallId := ctx.Param("card_id")

	respErr := cc.cardsService.DeleteCardFromCollection(collectionId, &models.Card{ScryfallID: scryfallId})
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}
