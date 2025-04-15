package services

import "github.com/ShenokZlob/collector-ouphe/collector-service/internal/models"

type CollectionsService struct {
	collectionRepository CollectionsRepositorer
}

type CollectionsRepositorer interface {
}

func NewCollectionService(collectionRepository CollectionsRepositorer) *CollectionsService {
	return &CollectionsService{
		collectionRepository: collectionRepository,
	}
}

func (cs CollectionsService) AllCollections(userId string) ([]*models.Collection, *models.ResponseErr) {
	panic("not implemented") // TODO: Implement
}

func (cs CollectionsService) CreateCollection(collection *models.Collection) (*models.Collection, *models.ResponseErr) {
	panic("not implemented") // TODO: Implement
}

func (cs CollectionsService) RenameCollection(collecion *models.Collection) *models.ResponseErr {
	panic("not implemented") // TODO: Implement
}

func (cs CollectionsService) DeleteCollection(collection *models.Collection) *models.ResponseErr {
	panic("not implemented") // TODO: Implement
}
