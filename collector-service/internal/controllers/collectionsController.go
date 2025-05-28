// internal/app/controllers/collections_controller.go
package controllers

import (
	"net/http"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/models"
	"github.com/ShenokZlob/collector-ouphe/pkg/contracts/collections"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/gin-gonic/gin"
)

// CollectionsController отвечает за работу с коллекциями
// @Tags Collections
// @BasePath /
type CollectionsController struct {
	collectionsService CollectionsServicer
	log                logger.Logger
}

type CollectionsServicer interface {
	AllUsersCollections(userId string) ([]*models.UserCollectionRef, *models.ResponseErr)
	CreateCollection(collection *models.Collection) (*models.Collection, *models.ResponseErr)
	RenameCollection(collection *models.Collection) (*models.Collection, *models.ResponseErr)
	DeleteCollection(collection *models.Collection) *models.ResponseErr
}

// NewCollectionsController создает контроллер коллекций
func NewCollectionsController(collectionsService CollectionsServicer, log logger.Logger) *CollectionsController {
	return &CollectionsController{
		collectionsService: collectionsService,
		log:                log.With(logger.String("controller", "collections")),
	}
}

// @Summary     Get user's collections
// @Description Получить список коллекций текущего пользователя
// @Tags        Collections
// @Security    BearerAuth
// @Produce     json
// @Success     200 {array} collections.Collection
// @Failure     401 {object} collections.ErrorResponse
// @Router      /collections [get]
func (cc CollectionsController) GetCollections(ctx *gin.Context) {
	cc.log.Info("CollectionsController.GetCollections called")

	userId, respErr := getUserFromCtx(ctx)
	if respErr != nil {
		cc.log.Error("Failed to get user ID from context", logger.Error(respErr))
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	list, respErr := cc.collectionsService.AllUsersCollections(userId)
	if respErr != nil {
		cc.log.Error("Failed to get user's collections", logger.Error(respErr))
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	var out []collections.Collection
	for _, c := range list {
		out = append(out, collections.Collection{ID: c.ID, Name: c.Name})
	}
	ctx.JSON(http.StatusOK, out)
}

// @Summary     Create new collection
// @Description Создать новую коллекцию с указанным именем
// @Tags        Collections
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       input body collections.CreateCollectionRequest true "Название новой коллекции"
// @Success     201 {object} collections.Collection
// @Failure     400,401 {object} collections.ErrorResponse
// @Router      /collections [post]
func (cc CollectionsController) CreateCollection(ctx *gin.Context) {
	userId, respErr := getUserFromCtx(ctx)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	var req collections.CreateCollectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, collections.ErrorResponse{Message: err.Error()})
		return
	}

	model := &models.Collection{UserID: userId, Name: req.Name}
	created, respErr := cc.collectionsService.CreateCollection(model)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	out := collections.Collection{ID: created.ID, Name: created.Name}
	ctx.JSON(http.StatusCreated, out)
}

// @Summary     Rename collection
// @Description Переименовать коллекцию по ID
// @Tags        Collections
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       id   path string                         true "Collection ID"
// @Param       input body collections.RenameCollectionRequest true "Новое имя коллекции"
// @Success     200 {object} collections.Collection
// @Failure     400,401,404 {object} collections.ErrorResponse
// @Router      /collections/{id} [patch]
func (cc CollectionsController) RenameCollection(ctx *gin.Context) {
	userId, respErr := getUserFromCtx(ctx)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	id := ctx.Param("id")
	var req collections.RenameCollectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, collections.ErrorResponse{Message: err.Error()})
		return
	}

	model := &models.Collection{ID: id, UserID: userId, Name: req.Name}
	updated, respErr := cc.collectionsService.RenameCollection(model)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	out := collections.Collection{ID: updated.ID, Name: updated.Name}
	ctx.JSON(http.StatusOK, out)
}

// @Summary     Delete collection
// @Description Удалить коллекцию по ID
// @Tags        Collections
// @Security    BearerAuth
// @Produce     json
// @Param       id path string true "Collection ID"
// @Success     204 "No Content"
// @Failure     401,404 {object} collections.ErrorResponse
// @Router      /collections/{id} [delete]
func (cc CollectionsController) DeleteCollection(ctx *gin.Context) {
	userId, respErr := getUserFromCtx(ctx)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	id := ctx.Param("id")
	respErr = cc.collectionsService.DeleteCollection(&models.Collection{ID: id, UserID: userId})
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func getUserFromCtx(ctx *gin.Context) (string, *models.ResponseErr) {
	val, ok := ctx.Get("userID")
	if !ok {
		return "", &models.ResponseErr{Status: http.StatusUnauthorized, Message: "Invalid user ID"}
	}
	userID, ok := val.(string)
	if !ok || userID == "" {
		return "", &models.ResponseErr{Status: http.StatusUnauthorized, Message: "Invalid user ID type"}
	}
	return userID, nil
}
