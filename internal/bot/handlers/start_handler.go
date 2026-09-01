package handlers

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type StartHandler struct {
	// userService services.UserService
}

// @inject
func NewStartHandler() *StartHandler {
	return &StartHandler{}
}

func (h *StartHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	userFirstName := update.Message.From.FirstName
	chatID := update.Message.Chat.ID

	welcomeText := fmt.Sprintf("👋 Hello, <b>%s</b>!\nWelcome to our platform bot.", userFirstName)

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      welcomeText,
		ParseMode: models.ParseModeHTML,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{
					{Text: "🚀 Open Web App", URL: "https://t.me"},
					{Text: "ℹ️ Help", CallbackData: "btn_help"},
				},
			},
		},
	})
}
