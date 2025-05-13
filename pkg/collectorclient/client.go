package collectorclient

import "github.com/ShenokZlob/collector-ouphe/pkg/contracts/collector"

type CollectorClient interface {
	RegisterUser(reqData *collector.RegisterRequest) (*collector.RegisterResponse, error)
	CheckUser(reqData *collector.CheckUserRequest) (*collector.CheckUserResponse, error)
}
