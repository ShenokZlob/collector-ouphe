package collectorclient

import (
	"github.com/ShenokZlob/collector-ouphe/pkg/contracts/auth"
	"github.com/ShenokZlob/collector-ouphe/pkg/contracts/collections"
)

type CollectorClient interface {
	CollectorClientAuth
	CollectorClientCollections
}

type CollectorClientAuth interface {
	RegisterUser(reqData *auth.RegisterRequest) (*auth.RegisterResponse, error)
	CheckUser(reqData *auth.CheckUserRequest) (*auth.CheckUserResponse, error)
}

type CollectorClientCollections interface {
	GetUserWithCollections(reqData *collections.GetCollectionsRequest) (*collections.GetCollectionsResponse, error)
	CreateCollection(reqData *collections.CreateCollectionRequest) (*collections.CreateCollectionResponse, error)
	RenameCollection(reqData *collections.RenameCollectionRequest) error
	DeleteCollection(reqData *collections.DeleteCollectionRequest) error
}
