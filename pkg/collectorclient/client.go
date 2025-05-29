package collectorclient

import (
	"context"

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
	GetUserCollections(ctx context.Context) ([]collections.Collection, error)
	CreateCollection(ctx context.Context, req *collections.CreateCollectionRequest) (*collections.Collection, error)
	RenameCollection(ctx context.Context, collectionID string, req *collections.RenameCollectionRequest) error
	DeleteCollection(ctx context.Context, collectionID string) error
	GetUsersCollectionByName(ctx context.Context, name string) (*collections.Collection, error)
}
