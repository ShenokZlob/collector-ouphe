package controllers

import "github.com/gin-gonic/gin"

type Controller struct {
	AuthController
	CollectionsController
	CardsController
}

type Servicer interface {
	AuthService
	CollectionsService
	CardsService
}

func NewControllers(service Servicer) *Controller {
	ctrl := Controller{}
	ctrl.authService = service
	ctrl.cardsService = service
	ctrl.collectionsService = service
	return &ctrl
}

func (c *Controller) PingPong(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}
