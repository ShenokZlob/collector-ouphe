package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Will check if the user is registered
func (h *AuthHandler) RegistrationMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		user := update.Message.From
		if ok := IsRegistered(user.ID); !ok {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Вы не зарегистрированы. Пожалуйста, зарегистрируйтесь.",
			})
			return
		}

		next(ctx, b, update)
	}
}
