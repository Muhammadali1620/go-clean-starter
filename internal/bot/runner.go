package bot

import (
	"context"
	"crypto/subtle"
	"log"
	"net/http"

	"github.com/go-telegram/bot"
	"github.com/labstack/echo/v4"
)

// StartPolling starts the bot in long polling mode (blocks or runs in goroutine).
func (s *BotService) StartPolling(ctx context.Context) {
	if s == nil || s.Client == nil {
		return
	}

	// Remove any existing webhook before starting polling
	_, _ = s.Client.DeleteWebhook(ctx, &bot.DeleteWebhookParams{
		DropPendingUpdates: false,
	})

	log.Println("🚀 Starting Telegram Bot in Long Polling mode...")
	s.Client.Start(ctx)
}

// SetupWebhook registers webhook in Telegram API.
func (s *BotService) SetupWebhook(ctx context.Context) error {
	if s == nil || s.Client == nil || s.Cfg.Bot.WebhookURL == "" {
		return nil
	}

	log.Printf("🔗 Setting up Telegram Webhook: %s", s.Cfg.Bot.WebhookURL)

	ok, err := s.Client.SetWebhook(ctx, &bot.SetWebhookParams{
		URL:         s.Cfg.Bot.WebhookURL,
		SecretToken: s.Cfg.Bot.SecretPath,
	})
	if err != nil || !ok {
		return err
	}

	log.Println("✅ Telegram Webhook registered successfully!")
	return nil
}

// WebhookEchoHandler returns an Echo HandlerFunc with Secret Token validation.
func (s *BotService) WebhookEchoHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		if s == nil || s.Client == nil {
			return c.NoContent(http.StatusOK)
		}

		// Security: verify Telegram secret token header
		if s.Cfg.Bot.SecretPath != "" {
			tokenHeader := c.Request().Header.Get("X-Telegram-Bot-Api-Secret-Token")
			if subtle.ConstantTimeCompare([]byte(tokenHeader), []byte(s.Cfg.Bot.SecretPath)) != 1 {
				return c.String(http.StatusUnauthorized, "invalid secret token")
			}
		}

		// Process the incoming update using native bot handler
		s.Client.WebhookHandler().ServeHTTP(c.Response().Writer, c.Request())
		return nil
	}
}
