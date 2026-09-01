package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"new_project/internal/core/server"
)

func RunHTTPServer(a *App) {
	e := echo.New()
	e.IPExtractor = echo.ExtractIPFromRealIPHeader()
	e.Validator = server.NewCustomValidator()
	e.HTTPErrorHandler = server.CustomHTTPErrorHandler

	// Middlewares
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			e.Logger.Printf("REQUEST: uri: %v, status: %v, method: %v\n", v.URI, v.Status, v.Method)
			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     a.Config.CORS.AllowOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	e.Use(middleware.BodyLimit(a.Config.Server.MaxBodySize))
	e.Use(middleware.Secure())
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			return strings.HasPrefix(path, "/api/v1/payments/") || strings.HasPrefix(path, "/api/v1/tg/webhook")
		},
		Store: middleware.NewRateLimiterMemoryStore(rate.Limit(a.Config.Server.MaxRateLimit)),
	}))

	// Health Check
	e.GET("/health", func(c echo.Context) error {
		dbStatus := "OK"
		if err := a.DB.Ping(); err != nil {
			dbStatus = "DOWN"
		}
		return c.JSON(http.StatusOK, map[string]string{
			"status":      "OK",
			"database":    dbStatus,
			"environment": a.Config.AppEnv,
		})
	})

	// Register Routes
	server.SetupRoutes(e, a.Container)

	// Telegram Webhook Setup (if configured)
	if a.Config.Bot.UpdateMode == "webhook" && a.Container.BotService != nil {
		_ = a.Container.BotService.SetupWebhook(context.Background())
	}

	// Start server with Graceful Shutdown
	go func() {
		log.Printf("🚀 HTTP Server listening on port :%s\n", a.Config.Server.Port)
		if err := e.Start(":" + a.Config.Server.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down HTTP server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to gracefully shutdown server: %v", err)
	}
	log.Println("👋 HTTP Server stopped successfully.")
}
