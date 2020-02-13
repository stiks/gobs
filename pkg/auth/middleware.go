package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RequiredAuth ...
func RequiredAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Get("AUTHORISED") == nil || !c.Get("AUTHORISED").(bool) {
				return echo.NewHTTPError(http.StatusUnauthorized, "authorisation required")
			}

			return next(c)
		}
	}
}
