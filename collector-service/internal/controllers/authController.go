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
	log         logger.Logger
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

func NewAuthController(authService AuthServicer, log logger.Logger) *AuthController {
	return &AuthController{
		authService: authService,
		log:         log.With(logger.String("controller", "auth")),
	}
}

func (ac AuthController) Register(ctx *gin.Context) {
	ac.log.With(logger.String("method", "Register")).Info("registering user")

	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ac.log.Error("failed to bind json", logger.String("error", err.Error()))
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
		ac.log.Error("failed to generate token", logger.String("error", err.Error()))
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
	ac.log.With(logger.String("method", "Who")).Info("getting user by telegram ID")

	userTelegramId := ctx.Param("telegram_id")
	if userTelegramId == "" {
		ac.log.Error("telegram ID is empty", logger.String("error", "telegram ID is empty"))
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
	ac.log.With(logger.String("method", "Login")).Info("logging in user")

	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ac.log.Error("failed to bind json", logger.String("error", err.Error()))
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
		ac.log.Error("failed to generate token", logger.String("error", err.Error()))
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
