package collectorclient

import "github.com/ShenokZlob/collector-ouphe/pkg/contracts/collector"

type CollectorClient interface {
	CheckUser(reqData *collector.CheckUserRequest) (*collector.CheckUserResponse, error)
	RegisterUser(reqData *collector.RegisterRequest) (*collector.RegisterResponse, error)
}
