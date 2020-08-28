package mock

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/repositories"
)

type statsRepository struct {
	stats *models.Stats
}

// NewStatsRepository ...
func NewStatsRepository() repositories.StatsRepository {
	return &statsRepository{
		stats: &models.Stats{
			Uptime:   time.Now(),
			Statuses: make(map[string]int),
		},
	}
}

// Process ...
func (r *statsRepository) Process(ctx echo.Context) error {
	status := strconv.Itoa(ctx.Response().Status)
	r.stats.Statuses[status]++

	return nil
}

// GetStats ...
func (r *statsRepository) GetStats(ctx echo.Context) (*models.Stats, error) {
	return r.stats, nil
}
