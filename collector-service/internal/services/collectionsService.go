package services

import (
	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/models"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
)

type CollectionsService struct {
	collectionRepository CollectionsRepositorer
	log                  logger.Logger
}

type CollectionsRepositorer interface {
	UsersCollections(userId string) ([]*models.UserCollectionRef, *models.ResponseErr)
	CreateCollection(collection *models.Collection) (*models.Collection, *models.ResponseErr)
	RenameCollection(collection *models.Collection) *models.ResponseErr
	DeleteCollection(collection *models.Collection) *models.ResponseErr
}

func NewCollectionsService(collectionRepository CollectionsRepositorer, log logger.Logger) *CollectionsService {
	return &CollectionsService{
		collectionRepository: collectionRepository,
		log:                  log.With(logger.String("service", "collections")),
	}
}

func (cs CollectionsService) AllUsersCollections(userId string) ([]*models.UserCollectionRef, *models.ResponseErr) {
	return cs.collectionRepository.UsersCollections(userId)
}

func (cs CollectionsService) CreateCollection(collection *models.Collection) (*models.Collection, *models.ResponseErr) {
	return cs.collectionRepository.CreateCollection(collection)
}

func (cs CollectionsService) RenameCollection(collecion *models.Collection) *models.ResponseErr {
	return cs.collectionRepository.RenameCollection(collecion)
}

func (cs CollectionsService) DeleteCollection(collection *models.Collection) *models.ResponseErr {
	return cs.collectionRepository.DeleteCollection(collection)
}
