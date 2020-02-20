package controllers_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/controllers"
	"github.com/stiks/gobs/pkg/helpers"
)

func TestControllers_Health_NewHealthController(t *testing.T) {
	assert.NotNil(t, controllers.NewHealthController())
}

func TestControllers_Health_Routes(t *testing.T) {
	e := echo.New()
	controllers.NewHealthController().Routes(e.Group("api"))

	c, _ := helpers.RequestTest(http.MethodGet, "/api/healthz", e)

	assert.Equal(t, 200, c)
}

func TestControllers_Health_HealthCheck(t *testing.T) {
	ctl := controllers.NewHealthController()
	_, ctx := helpers.RequestWithBody(http.MethodPut, "/", nil, echo.New())

	assert.NoError(t, ctl.HealthCheck(ctx))
}
