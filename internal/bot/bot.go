package bot

import (
	"context"
	"fmt"
	"log"

	"new_project/internal/core/config"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type BotService struct {
	Client *bot.Bot
	Cfg    *config.Config
}

// NewBotService initializes telegram bot client with global middlewares.
// @inject
func NewBotService(cfg *config.Config) (*BotService, error) {
	if cfg.Bot.Token == "" {
		log.Println("⚠️ Telegram Bot Token is empty. Bot will be disabled.")
		return nil, nil
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithMiddlewares(loggingMiddleware, recoveryMiddleware),
	}

	b, err := bot.New(cfg.Bot.Token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	return &BotService{
		Client: b,
		Cfg:    cfg,
	}, nil
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Unknown command. Use /start to begin.",
		})
	}
}

func recoveryMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("🔥 [BOT PANIC RECOVERED]: %v", r)
			}
		}()
		next(ctx, b, update)
	}
}

func loggingMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil {
			log.Printf("🤖 [TG MSG] From: %d | Text: %s", update.Message.Chat.ID, update.Message.Text)
		} else if update.CallbackQuery != nil {
			log.Printf("🤖 [TG CALLBACK] From: %d | Data: %s", update.CallbackQuery.From.ID, update.CallbackQuery.Data)
		}
		next(ctx, b, update)
	}
}
