package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/auth/dto"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
)

type authUsecaseImpl struct {
	collectorURL string
	log          logger.Logger
	localStorage *inMemoryStorage
}

func NewAuthUsecase(collectorURL string, log logger.Logger) *authUsecaseImpl {
	return &authUsecaseImpl{
		collectorURL: collectorURL,
		log:          log,
		// log:          log.With(logger.Field{Key: "usecase", String: "auth"}),
		localStorage: newInMemoryStorage(),
	}
}

func (a *authUsecaseImpl) RegisterUser(telegramID int64, username, telegramNickname string) (string, error) {
	reqData := dto.RegisterRequest{
		TelegramID:       telegramID,
		Username:         username,
		TelegramNickname: telegramNickname,
	}

	body, err := json.Marshal(reqData)
	if err != nil {
		a.log.Error("Failed to marshal request data", logger.Error(err))
		return "", err
	}

	resp, err := http.Post(a.collectorURL+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		a.log.Error("Failed to send request", logger.Error(err))
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		a.log.Error("Failed to register user", logger.String("status_code", fmt.Sprint(resp.StatusCode)))
		return "", fmt.Errorf("failed to register user, status code: %d", resp.StatusCode)
	}

	var res dto.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		a.log.Error("Failed to decode response", logger.Error(err))
		return "", err
	}

	// Store the token in local storage
	a.localStorage.AddUser(telegramID, res.Token)

	return res.Token, nil
}

func (a *authUsecaseImpl) IsRegistered(telegramID int64) bool {
	// Check in local struct
	if _, ok := a.localStorage.GetUser(telegramID); ok {
		a.log.Info("User found in local storage", logger.String("telegram_id", fmt.Sprint(telegramID)))
		return true
	}

	// TODO: Check in the database (Redis)

	// Check in the collector service
	a.log.Info("User not found in local storage, checking collector service", logger.String("telegram_id", fmt.Sprint(telegramID)))
	reqData := dto.RegisterRequest{
		TelegramID: telegramID,
	}

	body, err := json.Marshal(reqData)
	if err != nil {
		a.log.Error("Failed to marshal request data", logger.Error(err))
		return false
	}

	resp, err := http.Post(a.collectorURL+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		a.log.Error("Failed to send request", logger.Error(err))
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var res dto.RegisterResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			a.log.Error("Failed to decode response", logger.Error(err))
			return false
		}
	}

	return false
}
