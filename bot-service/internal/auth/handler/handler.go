package handler

import (
	"context"

	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/auth/dto"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type AuthHandler struct {
	Usecase AuthUsecase
	log     logger.Logger
}

type AuthUsecase interface {
	RegisterUser(dto.UserInfo) (string, error)
	IsRegistered(telegramID int64) bool
}

func NewAuthHandler(usecase AuthUsecase, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		Usecase: usecase,
		// log:     log.With(logger.Field{Key: "handler", String: "auth"}),
		log: log,
	}
}

// HandleRegister handles the registration of a new user.
func (h *AuthHandler) HandleRegister(ctx context.Context, b *bot.Bot, update *models.Update) {
	user := update.Message.From
	token, err := h.Usecase.RegisterUser(dto.UserInfo{
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

// type UserInfo struct {
// 	TelegramID int64
// 	FirstName  string
// 	LastName   string
// 	Username   string
// }

// RegistrationMiddleware checks if the user is registered before allowing access to other handlers.
func (h *AuthHandler) RegistrationMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		isReg := h.Usecase.IsRegistered(update.Message.From.ID)

		// Command /register
		if update.Message.Text == "/register" {
			if !isReg {
				h.HandleRegister(ctx, b, update)
			} else {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "Вы уже зарегистрированы.",
				})
			}
			return
		}

		// All other commands
		if !isReg {
			h.log.Info("User not registered", logger.Int("user_id", int(update.Message.From.ID)))
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Вы не зарегистрированы. Пожалуйста, используйте команду /register.",
			})
			return
		}

		next(ctx, b, update)
	}
}
