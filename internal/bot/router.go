package bot

import (
	"context"
	"log"

	bot_handlers "new_project/internal/bot/handlers"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// RegisterAllRoutes registers all injected handlers.
func (s *BotService) RegisterAllRoutes(startHandler *bot_handlers.StartHandler) {
	if s == nil || s.Client == nil {
		return
	}

	// 1. /start command
	if startHandler != nil {
		s.Client.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler.Handle)
	}

	// 2. Callback buttons
	s.Client.RegisterHandler(bot.HandlerTypeCallbackQueryData, "btn_help", bot.MatchTypeExact, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.CallbackQuery == nil {
			return
		}
		_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "This is help section!",
			ShowAlert:       true,
		})
	})

	log.Println("✅ Telegram Bot routes registered")
}
