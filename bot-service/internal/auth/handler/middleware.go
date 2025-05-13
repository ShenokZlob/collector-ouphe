package handler

import (
	"context"

	"github.com/ShenokZlob/collector-ouphe/bot-service/internal/authctx"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// RegistrationMiddleware checks if the user is registered before allowing access to other handlers.
func (h *AuthHandler) RegistrationMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		token, isReg := h.usecase.IsRegistered(update.Message.From.ID)

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

		ctx = authctx.WithJWT(ctx, token)

		next(ctx, b, update)
	}
}
