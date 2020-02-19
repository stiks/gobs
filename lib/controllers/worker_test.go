package controllers_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/controllers"
	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/services"
	"github.com/stiks/gobs/pkg/helpers"
)

var (
	_emailSrv = services.NewEmailService(mock.NewEmailRepository())
)

func TestControllers_Worker_NewUserController(t *testing.T) {
	assert.NotNil(t, controllers.NewWorkerController(_userSrv, _queueSrv, _emailSrv))
}

func TestControllers_Worker_Routes(t *testing.T) {
	t.Run("User password reset", func(t *testing.T) {
		e := echo.New()
		controllers.NewWorkerController(_userSrv, _queueSrv, _emailSrv).Routes(e.Group("worker"))

		c, _ := helpers.RequestTest(http.MethodPost, "/worker/user-password-reset", e)
		assert.Equal(t, 400, c)
	})

	t.Run("VUser profile updated", func(t *testing.T) {
		e := echo.New()
		controllers.NewUserController(_userSrv).Routes(e.Group("worker"))

		c, _ := helpers.RequestTest(http.MethodPost, "/worker/user-profile-updated", e)
		assert.Equal(t, 400, c)
	})

	t.Run("User password changed", func(t *testing.T) {
		e := echo.New()
		controllers.NewUserController(_userSrv).Routes(e.Group("worker"))

		c, _ := helpers.RequestTest(http.MethodPost, "/worker/user-password-changed", e)
		assert.Equal(t, 400, c)
	})
}
