package usecase

import (
	"fmt"

	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/auth/dto"
	"github.com/ShenokZlob/collector-ouphe/pkg/collectorclient"
	"github.com/ShenokZlob/collector-ouphe/pkg/contracts/collector"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
)

type authUsecaseImpl struct {
	collectorClient collectorclient.CollectorClient
	log             logger.Logger
	localStorage    *inMemoryStorage
}

func NewAuthUsecase(log logger.Logger, client collectorclient.CollectorClient) *authUsecaseImpl {
	return &authUsecaseImpl{
		log:             log,
		collectorClient: client,
		localStorage:    newInMemoryStorage(),
	}
}

// RegisterUser registers a user in the collector service and local storage.
// It returns the token if the registration is successful.
// If the user is already registered, it returns an empty string and nil error.
// If an error occurs during registration, it returns an empty string and the error.
func (a *authUsecaseImpl) RegisterUser(user dto.UserInfo) (string, error) {
	a.log.Info("Registering user", logger.String("method", "RegisterUser"), logger.String("telegram_id", fmt.Sprint(user.TelegramID)))

	reqData, err := a.collectorClient.RegisterUser(&collector.RegisterRequest{
		TelegramID: user.TelegramID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Username:   user.Username,
	})
	if err != nil {
		a.log.Error("Failed to register user in collector service", logger.Error(err))
		return "", err
	}

	a.log.Info("Registering user in local storage", logger.String("telegram_id", fmt.Sprint(user.TelegramID)))
	a.localStorage.AddUser(user.TelegramID, reqData.Token)

	return reqData.Token, nil
}

func (a *authUsecaseImpl) IsRegistered(telegramID int64) bool {
	// Check in local struct
	if _, ok := a.localStorage.GetUser(telegramID); ok {
		a.log.Info("User found in local storage", logger.String("telegram_id", fmt.Sprint(telegramID)))
		return true
	}

	// TODO: Check in the database (Redis)

	// Check in the collector service
	respData, err := a.collectorClient.CheckUser(&collector.CheckUserRequest{
		TelegramID: telegramID,
	})
	if err != nil {
		a.log.Error("Failed to check registration in collector service", logger.Error(err))
		return false
	}

	return respData.Success
}
