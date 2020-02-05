package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/controllers"
)

func TestNewHealthController(t *testing.T) {
	assert.NotNil(t, controllers.NewHealthController())
}

func TestHealthController_HealthCheck(t *testing.T) {
	ctl := controllers.NewHealthController()
	c := echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

	assert.NoError(t, ctl.HealthCheck(c))
}
