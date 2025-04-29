package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/auth/dto"
)

type authUsecaseImpl struct {
	collectorURL string
	localStorage *inMemoryStorage
}

func NewAuthUsecase(collectorURL string) *authUsecaseImpl {
	return &authUsecaseImpl{
		collectorURL: collectorURL,
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
		return "", err
	}

	resp, err := http.Post(a.collectorURL+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to register user, status code: %d", resp.StatusCode)
	}

	var res dto.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	// Store the token in local storage
	a.localStorage.AddUser(telegramID, res.Token)

	return res.Token, nil
}

func (a *authUsecaseImpl) IsRegistered(telegramID int64) bool {
	// Check in local struct
	if _, ok := a.localStorage.GetUser(telegramID); ok {
		return true
	}

	// TODO: Check in the database (Redis)

	// Check in the collector service
	reqData := dto.RegisterRequest{
		TelegramID: telegramID,
	}

	body, err := json.Marshal(reqData)
	if err != nil {
		return false
	}

	resp, err := http.Post(a.collectorURL+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var res dto.RegisterResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return false
		}
	}

	return false
}
