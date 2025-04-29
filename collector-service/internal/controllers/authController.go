package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/models"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	authService AuthServicer
	logger      logger.Logger
}

type AuthServicer interface {
	Register(*models.User) (*models.User, *models.ResponseErr)
	Who(userTelegramId string) (*models.User, *models.ResponseErr)
	Login(*models.User) *models.ResponseErr
}

type UserResponse struct {
	ID               string `json:"id"`
	TelegramID       int64  `json:"telegram_id"`
	Username         string `json:"username"`
	TelegramNickname string `json:"telegram_nickname,omitempty"`
}

func NewAuthController(authService AuthServicer, logger logger.Logger) *AuthController {
	return &AuthController{
		authService: authService,
		logger:      logger.With(),
	}
}

func (ac AuthController) Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	createdUser, respErr := ac.authService.Register(&user)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	token, err := generateToken(createdUser)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"user": UserResponse{
			ID:               user.ID,
			TelegramID:       createdUser.TelegramID,
			Username:         createdUser.Username,
			TelegramNickname: createdUser.TelegramNickname,
		},
		"token": token,
	})
}

func (ac AuthController) Who(ctx *gin.Context) {
	userTelegramId := ctx.Param("telegram_id")
	if userTelegramId == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, respErr := ac.authService.Who(userTelegramId)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	ctx.JSON(http.StatusOK, UserResponse{
		ID:               user.ID,
		TelegramID:       user.TelegramID,
		Username:         user.Username,
		TelegramNickname: user.TelegramNickname,
	})
}

func (ac AuthController) Login(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	respErr := ac.authService.Login(&user)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
		return
	}

	token, err := generateToken(&user)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}
