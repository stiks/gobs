package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stiks/gobs/lib/services"
)

// HealthControllerInterface ...
type HealthControllerInterface interface {
	Routes(g *echo.Group)
	LiveCheck(c echo.Context) error
	HealthCheck(c echo.Context) error
}

type healthController struct {
	stats services.StatsService
}

// NewHealthController returns a controller
func NewHealthController(statsSrv services.StatsService) HealthControllerInterface {
	return &healthController{
		stats: statsSrv,
	}
}

// Routes registers routes
func (ctl *healthController) Routes(g *echo.Group) {
	g.GET("/livez", ctl.LiveCheck)
	g.GET("/healthz", ctl.HealthCheck)
}

// LiveCheck ...
func (ctl *healthController) LiveCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"started": true})
}

// HealthCheck ...
func (ctl *healthController) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"healthy": true})
}
