package handler

import (
	"context"

	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type AuthHandler struct {
	Usecase AuthUsecase
	log     logger.Logger
}

type AuthUsecase interface {
	RegisterUser(telegramID int64, username, telegramNickname string) (string, error)
	IsRegistered(telegramID int64) bool
}

func NewAuthHandler(usecase AuthUsecase, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		Usecase: usecase,
		// log:     log.With(logger.Field{Key: "handler", String: "auth"}),
		log: log,
	}
}

func (h *AuthHandler) HandleRegister(ctx context.Context, b *bot.Bot, update *models.Update) {
	user := update.Message.From

	token, err := h.Usecase.RegisterUser(user.ID, user.Username, user.FirstName)
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
