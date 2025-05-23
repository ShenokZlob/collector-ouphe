package collectorclient

import (
	"github.com/ShenokZlob/collector-ouphe/pkg/contracts/auth"
	"github.com/ShenokZlob/collector-ouphe/pkg/contracts/collections"
	"github.com/gin-gonic/gin"
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
	GetUserCollections(ctx *gin.Context) ([]collections.Collection, error)
	CreateCollection(ctx *gin.Context, req *collections.CreateCollectionRequest) (*collections.Collection, error)
	RenameCollection(ctx *gin.Context, collectionID string, req *collections.RenameCollectionRequest) error
	DeleteCollection(ctx *gin.Context, collectionID string) error
}
