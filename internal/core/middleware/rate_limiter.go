package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"new_project/internal/dto"
)

func StrictRateLimit(reqPerSec float64, burst int) echo.MiddlewareFunc {
	config := echoMiddleware.RateLimiterConfig{
		Skipper: echoMiddleware.DefaultSkipper,
		Store: echoMiddleware.NewRateLimiterMemoryStoreWithConfig(
			echoMiddleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(reqPerSec),
				Burst:     burst,
				ExpiresIn: 3 * time.Minute,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return c.JSON(http.StatusTooManyRequests, dto.NewResponse("too many requests, please try again later", false))
		},
	}

	return echoMiddleware.RateLimiterWithConfig(config)
}
