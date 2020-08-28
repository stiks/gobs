package services

import (
	"github.com/labstack/echo/v4"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/repositories"
)

type statsService struct {
	repo repositories.StatsRepository
}

// StatsService ...
type StatsService interface {
	Process(ctx echo.Context) error
	GetStats(ctx echo.Context) (*models.Stats, error)
}

// NewStatsService ...
func NewStatsService(repo repositories.StatsRepository) StatsService {
	return &statsService{
		repo: repo,
	}
}

// Process ...
func (s *statsService) Process(ctx echo.Context) error {
	return s.repo.Process(ctx)
}

// GetStats ...
func (s *statsService) GetStats(ctx echo.Context) (*models.Stats, error) {
	return s.repo.GetStats(ctx)
}
