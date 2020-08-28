package repositories

import (
	"github.com/labstack/echo/v4"

	"github.com/stiks/gobs/lib/models"
)

// StatsRepository ...
type StatsRepository interface {
	Process(ctx echo.Context) error
	GetStats(ctx echo.Context) (*models.Stats, error)
}
