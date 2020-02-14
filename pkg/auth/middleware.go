package auth

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/stiks/gobs/pkg/env"
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

// EnableAuthorisation ...
func EnableAuthorisation() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		// CredentialsOptional: true,
		SigningKey:    []byte(env.MustGetString("AUTH_SECRET_KEY")),
		ContextKey:    "user",
		SigningMethod: middleware.AlgorithmHS256,
		BeforeFunc: func(c echo.Context) {
			c.Set("AUTHORISED", false)
		},
		SuccessHandler: func(c echo.Context) {
			/* we only authorise users when we have users details in Context */
			if c.Get("user") != nil && c.Get("user") != "" {
				claims := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)
				if claims == nil {
					return
				}

				uid, ok := claims["uid"]
				if !ok {
					return
				}

				c.Set("AUTHORISED", true)
				c.Set("USER_ID", uid)

				// check role
				if val, ok := claims["auth"]; ok {
					c.Set("ROLE", val)
				}
			}
		},
	})
}
