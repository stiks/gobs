package helpers

import (
	"github.com/labstack/echo/v4"
)

// ServerHeader middleware adds a `Server` header to the response.
func DefaultHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
			c.Response().Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
			c.Response().Header().Set("Pragma", "no-cache")
			c.Response().Header().Set("Strict-Transport-Security", "max-age=15724800; includeSubDomains")
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-Frame-Options", "SAMEORIGIN")

			c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")

			return next(c)
		}
	}
}
