package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthControllerInterface ...
type HealthControllerInterface interface {
	Routes(g *echo.Group)
	HealthCheck(c echo.Context) error
}

type healthController struct {
}

// NewHealthController returns a controller
func NewHealthController() HealthControllerInterface {
	return &healthController{}
}

// Routes registers routes
func (ctl *healthController) Routes(g *echo.Group) {
	g.GET("/healthz", ctl.HealthCheck)
}

// HealthCheck ...
func (ctl *healthController) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"healthy": true})
}
