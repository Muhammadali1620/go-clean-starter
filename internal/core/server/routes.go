package server

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"new_project/internal/core/container"
	"new_project/internal/core/middleware"
)

// SetupRoutes configures all the endpoints for the application using the DI container.
func SetupRoutes(e *echo.Echo, c *container.Container) {
	// Swagger UI Route
	swaggerGroup := e.Group("/swagger", middleware.BasicAuth(c.Config.Dev.BasicAuthUser, c.Config.Dev.BasicAuthPassword))
	swaggerGroup.GET("/*", echoSwagger.EchoWrapHandler(echoSwagger.PersistAuthorization(true)))

	// API V1 Group
	v1 := e.Group("/api/v1")
	{
		// Protected routes
		protected := v1.Group("")

		// Media routes
		media := protected.Group("/media")
		{
			media.POST("/upload", c.Handlers.Media.Upload)
		}

		// Tg Webhook Route
		if c.BotService != nil && c.Config.Bot.UpdateMode == "webhook" {
			v1.POST("/tg/webhook", c.BotService.WebhookEchoHandler())
		}

	}
}
