package controllers

import (
	"net/http"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/models"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/gin-gonic/gin"
)

type CollectionsController struct {
	collectionsService CollectionsServicer
	log                logger.Logger
}

type CollectionsServicer interface {
	AllUsersCollections(userId string) ([]*models.UserCollectionRef, *models.ResponseErr)
	CreateCollection(collection *models.Collection) (*models.Collection, *models.ResponseErr)
	RenameCollection(collecion *models.Collection) *models.ResponseErr
	DeleteCollection(collection *models.Collection) *models.ResponseErr
}

func NewCollectionsController(collectionsService CollectionsServicer, log logger.Logger) *CollectionsController {
	return &CollectionsController{
		collectionsService: collectionsService,
		log:                log.With(logger.String("controller", "collections")),
	}
}

func (cc CollectionsController) AllUsersCollections(ctx *gin.Context) {
	userId, respErr := getUserIDFromKeys(ctx)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	collections, respErr := cc.collectionsService.AllUsersCollections(userId)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ctx.JSON(http.StatusOK, collections)
}

func (cc CollectionsController) CreateCollection(ctx *gin.Context) {
	userId, respErr := getUserIDFromKeys(ctx)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	var collection models.Collection
	if err := ctx.ShouldBindJSON(&collection); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// вероятно, так делать не стоит...
	collection.UserID = userId

	createdColl, respErr := cc.collectionsService.CreateCollection(&collection)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ctx.JSON(http.StatusCreated, createdColl)
}

// хз как правильно реализовать
func (cc CollectionsController) RenameCollection(ctx *gin.Context) {
	userId, respErr := getUserIDFromKeys(ctx)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	var collection models.Collection
	if err := ctx.ShouldBindJSON(&collection); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// вероятно, так делать не стоит...
	collection.UserID = userId

	respErr = cc.collectionsService.RenameCollection(&collection)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (cc CollectionsController) DeleteCollection(ctx *gin.Context) {
	userId, respErr := getUserIDFromKeys(ctx)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}
	collectionId := ctx.Param("id")

	respErr = cc.collectionsService.DeleteCollection(&models.Collection{
		ID:     collectionId,
		UserID: userId,
	})
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}
	ctx.Status(http.StatusNoContent)

}

func getUserIDFromKeys(ctx *gin.Context) (string, *models.ResponseErr) {
	userId, ok := ctx.Keys["user_id"].(string)
	if !ok {
		return "", &models.ResponseErr{
			Status:  http.StatusUnauthorized,
			Message: "Invalid user ID",
		}
	}
	return userId, nil
}
