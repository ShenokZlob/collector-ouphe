package usecase

import (
	"context"
	"fmt"

	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/auth/dto"
	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/session"
	"github.com/ShenokZlob/collector-ouphe/pkg/collectorclient"
	"github.com/ShenokZlob/collector-ouphe/pkg/contracts/auth"

	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
)

type authUsecaseImpl struct {
	log             logger.Logger
	collectorClient collectorclient.CollectorClientAuth
	cache           *session.Cache
}

func NewAuthUsecase(log logger.Logger, client collectorclient.CollectorClientAuth, cache *session.Cache) *authUsecaseImpl {
	return &authUsecaseImpl{
		log:             log,
		collectorClient: client,
		cache:           cache,
	}
}

// RegisterUser registers a user in the collector service and local storage.
// It returns the token if the registration is successful.
// If the user is already registered, it returns an empty string and nil error.
// If an error occurs during registration, it returns an empty string and the error.
func (a *authUsecaseImpl) RegisterUser(user dto.UserInfo) (string, error) {
	a.log.Info("Registering user", logger.String("method", "RegisterUser"), logger.String("telegram_id", fmt.Sprint(user.TelegramID)))

	reqData, err := a.collectorClient.RegisterUser(&auth.RegisterRequest{
		TelegramID: user.TelegramID,
		Username:   user.Username,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
	})
	if err != nil {
		a.log.Error("Failed to register user in collector service", logger.Error(err))
		return "", err
	}

	a.log.Info("Save new user in cache", logger.String("telegram_id", fmt.Sprint(user.TelegramID)))
	a.cache.Set(context.TODO(), fmt.Sprint(user.TelegramID), reqData.Token)

	return reqData.Token, nil
}

// IsRegistered return token JWT (if it exist)
func (a *authUsecaseImpl) IsRegistered(telegramID int64) (string, bool) {
	// Check in the cache (Redis)
	value, err := a.checkInCache(fmt.Sprint(telegramID))
	if err == nil {
		return value, true
	}

	// Check in the collector service
	respData, err := a.collectorClient.CheckUser(&auth.CheckUserRequest{
		TelegramID: telegramID,
	})
	if err != nil {
		a.log.Error("Failed to check registration in collector service", logger.Error(err))
		return "", false
	}

	// If the user is registered, save the token in the cache
	if respData.Success {
		a.log.Info("User found in collector service", logger.String("telegram_id", fmt.Sprint(telegramID)))
		a.cache.Set(context.TODO(), fmt.Sprint(telegramID), respData.Token)
	}

	return respData.Token, respData.Success
}

func (a *authUsecaseImpl) checkInCache(telegramID string) (string, error) {
	// Check in the cache (Redis)
	value, err := a.cache.Get(context.TODO(), telegramID)
	if err != nil {
		a.log.Error("Failed to get user from cache", logger.Error(err))
		return "", err
	}

	a.log.Info("User found in cache", logger.String("telegram_id", telegramID))
	return value, nil
}
