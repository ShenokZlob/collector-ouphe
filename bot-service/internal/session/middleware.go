package session

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ctxKey string

const stateKey ctxKey = "userState"

func WithState(ctx context.Context, state string) context.Context {
	return context.WithValue(ctx, stateKey, state)
}

func GetState(ctx context.Context) (string, bool) {
	s, ok := ctx.Value(stateKey).(string)
	return s, ok
}

func Middleware(mgr Manager) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			if update.Message == nil || update.Message.From == nil {
				return
			}

			state, err := mgr.GetState(ctx, update.Message.From.ID)
			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "Что-то пошло не так, попробуйте позже.",
				})
				return
			}

			if state != "" {
				ctx = WithState(ctx, state)
			}

			next(ctx, b, update)
		}
	}
}

func CancelHandler(mgr Manager) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil || update.Message.From == nil {
			return
		}

		mgr.ClearState(ctx, update.Message.From.ID)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Действие отменено. Введите команду.",
		})
	}
}
