package server

import (
	"log"

	"github.com/go-telegram/bot"

	"new_project/internal/core/container"
)

// SetupBotRoutes binds all telegram bot handlers centrally.
func SetupBotRoutes(c *container.Container) {
	if c.BotService == nil || c.BotService.Client == nil {
		return
	}

	b := c.BotService.Client
	h := c.BotHandlers

	// Register handlers from c.BotHandlers here:
	if h.Start != nil {
		b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, h.Start.Handle)
	}

	log.Println("✅ Telegram Bot routes registered")
}
