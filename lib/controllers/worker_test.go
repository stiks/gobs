package controllers_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/controllers"
	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/services"
	"github.com/stiks/gobs/pkg/helpers"
)

var (
	_         = os.Setenv("PUBLIC_NAME", "something")
	_         = os.Setenv("PUBLIC_HOSTNAME", "something")
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

func TestControllers_Worker_UserPasswordReset(t *testing.T) {
	ctl := controllers.NewWorkerController(_userSrv, _queueSrv, _emailSrv)

	t.Run("Existing user", func(t *testing.T) {
		data := models.WorkerRequest{ID: helpers.UUIDFromString(t, "3ab1ba2a-6031-4e34-aae3-dcd43a987775")}
		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", data, echo.New())

		assert.NoError(t, ctl.UserPasswordReset(ctx))
	})

	t.Run("Non-existing user", func(t *testing.T) {
		data := models.WorkerRequest{ID: helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88")}
		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", data, echo.New())

		err := ctl.UserPasswordReset(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "user not found", "error message %s", "formatted")
		}
	})
}

func TestControllers_Worker_UserProfileUpdated(t *testing.T) {
	ctl := controllers.NewWorkerController(_userSrv, _queueSrv, _emailSrv)

	t.Run("Existing user", func(t *testing.T) {
		data := models.WorkerRequest{ID: helpers.UUIDFromString(t, "3ab1ba2a-6031-4e34-aae3-dcd43a987775")}
		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", data, echo.New())

		assert.NoError(t, ctl.UserProfileUpdated(ctx))
	})

	t.Run("Non-existing user", func(t *testing.T) {
		data := models.WorkerRequest{ID: helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88")}
		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", data, echo.New())

		err := ctl.UserProfileUpdated(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "user not found", "error message %s", "formatted")
		}
	})
}

func TestControllers_Worker_ConfirmEmail(t *testing.T) {
	ctl := controllers.NewWorkerController(_userSrv, _queueSrv, _emailSrv)

	t.Run("Existing user", func(t *testing.T) {
		data := models.WorkerRequest{ID: helpers.UUIDFromString(t, "3ab1ba2a-6031-4e34-aae3-dcd43a987775")}
		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", data, echo.New())

		assert.NoError(t, ctl.ConfirmEmail(ctx))
	})

	t.Run("Non-existing user", func(t *testing.T) {
		data := models.WorkerRequest{ID: helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88")}
		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", data, echo.New())

		err := ctl.ConfirmEmail(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "user not found", "error message %s", "formatted")
		}
	})
}

func TestControllers_Worker_UserPasswordChanged(t *testing.T) {
	ctl := controllers.NewWorkerController(_userSrv, _queueSrv, _emailSrv)

	t.Run("Existing user", func(t *testing.T) {
		data := models.WorkerRequest{ID: helpers.UUIDFromString(t, "3ab1ba2a-6031-4e34-aae3-dcd43a987775")}
		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", data, echo.New())

		assert.NoError(t, ctl.UserPasswordChanged(ctx))
	})

	t.Run("Non-existing user", func(t *testing.T) {
		data := models.WorkerRequest{ID: helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88")}
		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", data, echo.New())

		err := ctl.UserPasswordChanged(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "user not found", "error message %s", "formatted")
		}
	})
}
