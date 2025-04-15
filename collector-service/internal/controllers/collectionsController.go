package controllers

import "github.com/gin-gonic/gin"

type CollectionsController struct {
	collectionsService CollectionsService
}

type CollectionsService interface {
	AllCollections()
	CreateCollection()
	RenameCollection()
	DeleteCollection()
}

func NewCollectionsController(collectionsService CollectionsService) *CollectionsController {
	return &CollectionsController{
		collectionsService: collectionsService,
	}
}

func (cc *CollectionsController) AllCollections(ctx *gin.Context) {}

func (cc *CollectionsController) CreateCollection(ctx *gin.Context) {}

func (cc *CollectionsController) RenameCollection(ctx *gin.Context) {}

func (cc *CollectionsController) DeleteCollection(ctx *gin.Context) {}
