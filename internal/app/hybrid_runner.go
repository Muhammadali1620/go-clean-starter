package app

import (
	"context"
	"log"
)

func RunHybrid(a *App) {
	log.Println("⚡ Starting HYBRID Mode: (HTTP Server + Telegram Bot)...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Run Bot in background (if token provided)
	if a.Container.BotService != nil && a.Config.Bot.UpdateMode == "polling" {
		go a.Container.BotService.StartPolling(ctx)
	}

	// 2. Run HTTP Server (blocks until SIGTERM)
	RunHTTPServer(a)

	// Cancel bot context when HTTP shuts down
	cancel()
}
