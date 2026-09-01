package middleware

import (
	"crypto/subtle"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// BasicAuth returns an Echo middleware for basic authentication using credentials from Config.
func BasicAuth(expectedUser, expectedPassword string) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		userMatch := subtle.ConstantTimeCompare([]byte(username), []byte(expectedUser)) == 1
		passMatch := subtle.ConstantTimeCompare([]byte(password), []byte(expectedPassword)) == 1

		return userMatch && passMatch, nil
	})
}
