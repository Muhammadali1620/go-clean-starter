package config

type BotConfig struct {
	Token      string // Telegram bot token from @BotFather
	UpdateMode string // "polling" or "webhook"
	WebhookURL string // Full public URL (e.g., https://api.domain.com/api/v1/tg/webhook)
	SecretPath string // Secret token for header X-Telegram-Bot-Api-Secret-Token
}

func loadBotConfig() BotConfig {
	return BotConfig{
		Token:      GetEnv("BOT_TOKEN", ""),
		UpdateMode: GetEnv("BOT_UPDATE_MODE", "polling"), // polling | webhook
		WebhookURL: GetEnv("BOT_WEBHOOK_URL", ""),
		SecretPath: GetEnv("BOT_SECRET_PATH", "tg-webhook-secret"),
	}
}
