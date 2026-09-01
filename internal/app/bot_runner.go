package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func RunBot(a *App) {
	if a.Container.BotService == nil {
		log.Fatal("❌ Telegram Bot is not configured. Check BOT_TOKEN in .env")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		a.Container.BotService.StartPolling(ctx)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Stopping Telegram Bot...")
	cancel()
	log.Println("👋 Telegram Bot stopped.")
}
