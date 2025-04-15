package controllers

import "github.com/gin-gonic/gin"

type CardsController struct {
	cardsService CardsService
}

type CardsService interface {
	AllCardsByCollection()
	AddCardToCollection()
	SetCardCount()
	DeleteCard()
}

func NewCardsController(cardsService CardsService) *CardsController {
	return &CardsController{
		cardsService: cardsService,
	}
}

func (cc *CardsController) AllCardsByCollection(ctx *gin.Context) {}

func (cc *CardsController) AddCardToCollection(ctx *gin.Context) {}

func (cc *CardsController) SetCardCount(ctx *gin.Context) {}

func (cc *CardsController) DeleteCard(ctx *gin.Context) {}
