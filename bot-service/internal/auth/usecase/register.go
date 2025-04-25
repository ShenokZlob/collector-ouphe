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
}

func NewAuthUsecase(collectorURL string) *authUsecaseImpl {
	return &authUsecaseImpl{
		collectorURL: collectorURL,
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

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to register user, status code: %d", resp.StatusCode)
	}

	var res dto.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.Token, nil
}
