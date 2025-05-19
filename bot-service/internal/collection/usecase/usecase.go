package usecase

import (
	"context"
	"fmt"

	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/authctx"
	"github.com/ShenokZlob/collector-ouphe/pkg/collectorclient"
	"github.com/ShenokZlob/collector-ouphe/pkg/contracts/collector"
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

	token, ok := authctx.GetJWT(ctx)
	if !ok {
		u.log.Error("The user doesn't have a jwt token")
		return nil, fmt.Errorf("don't have a jwt token")
	}

	respData, err := u.collectorClient.GetUserWithCollections(&collector.GetCollectionsRequest{
		Token: token,
	})
	if err != nil {
		u.log.Error("Error when accessing the collector service", logger.Error(err))
		return nil, err
	}

	collectionsNames := make([]string, len(respData.Collections))
	for _, v := range respData.Collections {
		collectionsNames = append(collectionsNames, v.Name)
	}

	// TODO: Update cache (?)

	return collectionsNames, nil
}

func (u *collectionUsecaseImpl) CreateaCollection(ctx context.Context, name string) error {
	u.log.Info("Create collecion", logger.String("method", "CreateCollection"))

	token, ok := authctx.GetJWT(ctx)
	if !ok {
		u.log.Error("The user doesn't have a jwt token")
		return fmt.Errorf("don't have a jwt token")
	}

	_, err := u.collectorClient.CreateCollection(&collector.CreateCollectionRequest{
		Token:          token,
		CollectionName: name,
	})
	if err != nil {
		u.log.Error("Error when accessing the collector service", logger.Error(err))
		return err
	}

	// TODO: Save respData in cache

	return nil
}

func (u *collectionUsecaseImpl) RenameCollection(ctx context.Context, oldName, newName string) error {
	u.log.Info("Rename collecion", logger.String("method", "ReanameCollection"))

	token, ok := authctx.GetJWT(ctx)
	if !ok {
		u.log.Error("The user doesn't have a jwt token")
		return fmt.Errorf("don't have a jwt token")
	}

	// TODO: get collection ID by the old name (from cache or collector service)
	collectionID := getCollectionIdByName(oldName)

	err := u.collectorClient.RenameCollection(&collector.RenameCollectionRequest{
		Token:             token,
		CollectionID:      collectionID,
		NewCollectionName: newName,
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

	token, ok := authctx.GetJWT(ctx)
	if !ok {
		u.log.Error("The user doesn't have a jwt token")
		return fmt.Errorf("don't have a jwt token")
	}

	collectionID := getCollectionIdByName(name)

	err := u.collectorClient.DeleteCollection(&collector.DeleteCollectionRequest{
		Token:        token,
		CollectionID: collectionID,
	})
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
