package controllers

import (
	"net/http"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/models"
	"github.com/ShenokZlob/collector-ouphe/pkg/contracts/auth"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/gin-gonic/gin"
)

// AuthController отвечает за регистрацию, логин и проверку пользователя
// @Tags Auth
// @BasePath /
type AuthController struct {
	log         logger.Logger
	authService AuthServicer
}

type AuthServicer interface {
	Register(*models.User) (string, *models.ResponseErr)
	Who(userTelegramId string) (string, *models.ResponseErr)
	Login(*models.User) (string, *models.ResponseErr)
}

type UserResponse struct {
	ID         string `json:"id"`
	TelegramID int64  `json:"telegram_id"`
	FirstName  string `json:"first_name"`
	Username   string `json:"username,omitempty"`
}

func NewAuthController(authService AuthServicer, log logger.Logger) *AuthController {
	return &AuthController{
		authService: authService,
		log:         log.With(logger.String("controller", "auth")),
	}
}

// @Summary     Register user
// @Description Регистрация пользователя, возвращает JWT
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       input body auth.RegisterRequest true "Данные для регистрации"
// @Success     201 {object} auth.RegisterResponse
// @Failure     400 {object} models.ResponseErr
// @Router      /register [post]
func (ac AuthController) Register(ctx *gin.Context) {
	var req auth.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ac.log.Error("failed to bind json", logger.String("error", err.Error()))
		ctx.JSON(http.StatusBadRequest, models.ResponseErr{Message: err.Error()})
		return
	}

	userModel := &models.User{
		TelegramID: req.TelegramID,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Username:   req.Username,
	}
	token, respErr := ac.authService.Register(userModel)
	if respErr != nil {
		ctx.JSON(respErr.Status, respErr)
		return
	}

	ctx.JSON(http.StatusCreated, auth.RegisterResponse{Token: token})
}

// @Summary     Check user by Telegram ID
// @Description Проверяет существование пользователя и возвращает JWT
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       telegram_id path int true "Telegram ID"
// @Success     200 {object} auth.CheckUserResponse
// @Failure     400,404 {object} models.ResponseErr
// @Router      /user/telegram/{telegram_id} [get]
func (ac AuthController) Who(ctx *gin.Context) {
	telegramID := ctx.Param("telegram_id")
	token, respErr := ac.authService.Who(telegramID)
	if respErr != nil {
		ctx.JSON(respErr.Status, respErr)
		return
	}

	ctx.JSON(http.StatusOK, auth.CheckUserResponse{Token: token, Success: true})
}

// @Summary     Login user
// @Description Логин по Telegram ID, возвращает JWT
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       input body auth.CheckUserRequest true "Telegram ID для логина"
// @Success     200 {object} auth.CheckUserResponse
// @Failure     400,401 {object} models.ResponseErr
// @Router      /login [post]
func (ac AuthController) Login(ctx *gin.Context) {
	var req auth.CheckUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ac.log.Error("failed to bind json", logger.String("error", err.Error()))
		ctx.JSON(http.StatusBadRequest, models.ResponseErr{Message: err.Error()})
		return
	}

	userModel := &models.User{TelegramID: req.TelegramID}
	token, respErr := ac.authService.Login(userModel)
	if respErr != nil {
		ctx.JSON(respErr.Status, respErr)
		return
	}

	ctx.JSON(http.StatusOK, auth.CheckUserResponse{Token: token, Success: true})
}
