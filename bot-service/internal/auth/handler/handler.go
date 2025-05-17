package handler

import (
	"context"

	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/auth/dto"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type AuthHandler struct {
	usecase AuthUsecase
	log     logger.Logger
}

type AuthUsecase interface {
	RegisterUser(dto.UserInfo) (string, error)
	IsRegistered(telegramID int64) (string, bool)
}

func NewAuthHandler(usecase AuthUsecase, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		usecase: usecase,
		log:     log.With(logger.String("handler", "auth")),
	}
}

// HandleRegister handles the registration of a new user.
func (h *AuthHandler) HandleRegister(ctx context.Context, b *bot.Bot, update *models.Update) {
	user := update.Message.From
	token, err := h.usecase.RegisterUser(dto.UserInfo{
		TelegramID: user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Username:   user.Username,
	})
	if err != nil {
		h.log.Error("Failed to register user", logger.Error(err))
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка регистрации. Попробуйте позже.",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Вы успешно зарегистрированы. Ваш токен: " + token,
	})
}
