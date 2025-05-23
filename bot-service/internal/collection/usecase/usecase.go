package usecase

import (
	"context"

	"github.com/ShenokZlob/collector-ouphe/pkg/collectorclient"
	"github.com/ShenokZlob/collector-ouphe/pkg/contracts/collections"

	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
)

type collectionUsecaseImpl struct {
	log             logger.Logger
	collectorClient collectorclient.CollectorClientCollections
}

func NewCollectionUsecaseImpl(log logger.Logger, client collectorclient.CollectorClientCollections) *collectionUsecaseImpl {
	return &collectionUsecaseImpl{
		log:             log,
		collectorClient: client,
	}
}

func (u *collectionUsecaseImpl) GetCollecionsList(ctx context.Context) ([]string, error) {
	u.log.Info("Get collection's list", logger.String("method", "GetCollectionsList"))

	collections, err := u.collectorClient.GetUserCollections(ctx)
	if err != nil {
		u.log.Error("Error when accessing the collector service", logger.Error(err))
		return nil, err
	}

	collectionsNames := make([]string, len(collections))
	for _, v := range collections {
		collectionsNames = append(collectionsNames, v.Name)
	}

	// TODO: Update cache (?)

	return collectionsNames, nil
}

func (u *collectionUsecaseImpl) CreateaCollection(ctx context.Context, name string) (*collections.Collection, error) {
	u.log.Info("Create collecion", logger.String("method", "CreateCollection"))

	collection, err := u.collectorClient.CreateCollection(ctx, &collections.CreateCollectionRequest{
		Name: name,
	})
	if err != nil {
		u.log.Error("Error when accessing the collector service", logger.Error(err))
		return nil, err
	}

	// TODO: Save respData in cache

	return collection, nil
}

func (u *collectionUsecaseImpl) RenameCollection(ctx context.Context, oldName, newName string) error {
	u.log.Info("Rename collecion", logger.String("method", "ReanameCollection"))

	// TODO: get collection ID by the old name (from cache or collector service)
	collectionID := getCollectionIdByName(oldName)

	err := u.collectorClient.RenameCollection(ctx, collectionID, &collections.RenameCollectionRequest{
		Name: newName,
	})
	if err != nil {
		u.log.Error("Error when accessing the collector service", logger.Error(err))
		return err
	}

	// TODO: Update cache

	return nil
}

func (u *collectionUsecaseImpl) DeleteCollection(ctx context.Context, name string) error {
	u.log.Info("Delete collection", logger.String("method", "DeleteCollection"))

	collectionID := getCollectionIdByName(name)

	err := u.collectorClient.DeleteCollection(ctx, collectionID)
	if err != nil {
		u.log.Error("Error when accessing the collector service", logger.Error(err))
		return err
	}

	// TODO: Update cache

	return nil
}

func getCollectionIdByName(name string) string {
	return "Implement me!" + name
}
