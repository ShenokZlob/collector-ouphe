package controllers

import (
	"net/http"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/models"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService AuthService
}

type AuthService interface {
	Register(*models.User) *models.ResponseErr
	Who(userTelegramId string) (*models.User, *models.ResponseErr)
}

type UserResponse struct {
	TelegramID       int64  `json:"telegram_id"`
	Name             string `json:"name"`
	TelegramNickname string `json:"telegram_nickname,omitempty"`
}

func NewAuthController(authService AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ac *AuthController) Register(ctx *gin.Context) {
	var user models.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	respErr := ac.authService.Register(&user)
	if respErr != nil {
		ctx.AbortWithStatusJSON(respErr.Status, respErr)
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}

func (ac *AuthController) Who(ctx *gin.Context) {
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
		TelegramID:       user.TelegramID,
		Name:             user.Name,
		TelegramNickname: user.TelegramNickname,
	})
}
