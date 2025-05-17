package collectorclient

import "github.com/ShenokZlob/collector-ouphe/pkg/contracts/collector"

type CollectorClient interface {
	CollectorClientAuth
	CollectorClientCollections
}

type CollectorClientAuth interface {
	RegisterUser(reqData *collector.RegisterRequest) (*collector.RegisterResponse, error)
	CheckUser(reqData *collector.CheckUserRequest) (*collector.CheckUserResponse, error)
}

type CollectorClientCollections interface {
	GetUserWithCollections(reqData *collector.GetCollectionsRequest) (*collector.GetCollectionsResponse, error)
	CreateCollection(reqData *collector.CreateCollectionRequest) (*collector.CreateCollectionResponse, error)
	RenameCollection(reqData *collector.RenameCollectionRequest) error
	DeleteCollection(reqData *collector.DeleteCollectionRequest) error
}
