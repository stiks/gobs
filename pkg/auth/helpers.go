package auth

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/stiks/gobs/pkg/parser"
)

// GetUserID ...
func GetUserID(c echo.Context) (uuid.UUID, error) {
	id, err := parser.String(c.Get("USER_ID"), nil)
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuid.Parse(id)
}
