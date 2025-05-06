package services

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/models"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
)

type AuthService struct {
	authRepository AuthRepositorer
	log            logger.Logger
}

type AuthRepositorer interface {
	CreateUser(user *models.User) (*models.User, *models.ResponseErr)
	FindUserByTelegramID(telegramId int64) (*models.User, *models.ResponseErr)
}

func NewAuthService(authRepository AuthRepositorer, log logger.Logger) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		log:            log.With(logger.String("service", "auth")),
	}
}

// Register creates a new user in the database.
func (as AuthService) Register(user *models.User) (*models.User, *models.ResponseErr) {
	as.log.With(logger.String("method", "Register")).Info("registering user")

	respErr := validateUser(user)
	if respErr != nil {
		as.log.Error("failed to validate user", logger.Error(respErr))
		return nil, respErr
	}

	createdUser, respErr := as.authRepository.CreateUser(user)
	if respErr != nil {
		as.log.Error("failed to create user", logger.Error(respErr))
		return nil, respErr
	}

	return createdUser, nil
}

// Who retrieves a user by their Telegram ID
// and returns the user object if found, or an error if not found.
func (as AuthService) Who(telegramId string) (*models.User, *models.ResponseErr) {
	as.log.With(logger.String("method", "Who")).Info("getting user by telegram ID")

	tgIdInt64, respErr := parseTelegramID(telegramId)
	if respErr != nil {
		as.log.Error("failed to parse telegram ID", logger.Error(respErr))
		return nil, respErr
	}

	user, respErr := as.authRepository.FindUserByTelegramID(tgIdInt64)
	if respErr != nil {
		as.log.Error("failed to find user by telegram ID", logger.Error(respErr))
		return nil, respErr
	}

	return user, nil
}

// Login checks if the user exists in the database
func (as AuthService) Login(user *models.User) *models.ResponseErr {
	as.log.With(logger.String("method", "Login")).Info("logging in user")

	respErr := validateUser(user)
	if respErr != nil {
		return respErr
	}

	// Check if user exists in the database
	existingUser, respErr := as.authRepository.FindUserByTelegramID(user.TelegramID)
	if respErr != nil {
		as.log.Error("failed to find user by telegram ID", logger.Error(respErr))
		return respErr
	}

	if existingUser == nil {
		as.log.Error("user not found", logger.String("error", "user not found"))
		return &models.ResponseErr{
			Status:  http.StatusUnauthorized,
			Message: "Invalid credentials",
		}
	}

	if existingUser.Username != user.Username {
		as.log.Error("invalid credentials", logger.String("error", "invalid credentials"))
		return &models.ResponseErr{
			Status:  http.StatusUnauthorized,
			Message: "Invalid credentials",
		}
	}

	return nil
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
