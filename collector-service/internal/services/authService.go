package services

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/models"
)

type AuthService struct {
	authRepository AuthRepositorer
}

type AuthRepositorer interface {
	CreateUser(user *models.User) (*models.User, *models.ResponseErr)
	FindUserByTelegramID(telegramId int64) (*models.User, *models.ResponseErr)
}

func NewAuthService(authRepository AuthRepositorer) *AuthService {
	return &AuthService{
		authRepository: authRepository,
	}
}

func (as AuthService) Register(user *models.User) (*models.User, *models.ResponseErr) {
	respErr := validateUser(user)
	if respErr != nil {
		return nil, respErr
	}

	createdUser, respErr := as.authRepository.CreateUser(user)
	if respErr != nil {
		return nil, respErr
	}

	return createdUser, nil
}

func (as AuthService) Who(telegramId string) (*models.User, *models.ResponseErr) {
	tgIdInt64, respErr := parseTelegramID(telegramId)
	if respErr != nil {
		return nil, respErr
	}

	user, respErr := as.authRepository.FindUserByTelegramID(tgIdInt64)
	if respErr != nil {
		return nil, respErr
	}

	return user, nil
}

func validateUser(user *models.User) *models.ResponseErr {
	if user.TelegramID == 0 {
		return &models.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid user telegram ID",
		}
	}
	if user.Username == "" {
		return &models.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid user username",
		}
	}
	return nil
}

func parseTelegramID(telegramId string) (int64, *models.ResponseErr) {
	telegramIdInt64, err := strconv.ParseInt(telegramId, 10, 64)
	if err != nil {
		return 0, &models.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("Telegram_id parse error: %v", err),
		}
	}
	return telegramIdInt64, nil
}
