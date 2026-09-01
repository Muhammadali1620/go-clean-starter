package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// BotHandler defines standard interface for telegram command/callback handlers.
type BotHandler interface {
	Handle(ctx context.Context, b *bot.Bot, update *models.Update)
}
