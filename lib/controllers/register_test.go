package controllers_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/controllers"
	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/pkg/helpers"
)

func TestControllers_Register_NewRegisterController(t *testing.T) {
	assert.NotNil(t, controllers.NewRegisterController(_userSrv))
}

func TestControllers_Register_Routes(t *testing.T) {
	t.Run("User registration", func(t *testing.T) {
		e := echo.New()
		controllers.NewRegisterController(_userSrv).Routes(e.Group("api"))

		c, _ := helpers.RequestTest(http.MethodPost, "/api/register", e)
		assert.Equal(t, 400, c)
	})
}

func TestControllers_Register_Register(t *testing.T) {
	ctl := controllers.NewRegisterController(_userSrv)

	t.Run("Non-existing user", func(t *testing.T) {
		user := models.CreateUser{
			ID:        uuid.New(),
			FirstName: "John",
			LastName:  "Snow",
			Email:     "google@test.com",
			Password:  "Test123456",
		}

		rec, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", user, echo.New())

		err := ctl.User(ctx)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

	t.Run("Existing user", func(t *testing.T) {
		user := models.CreateUser{
			ID:        uuid.New(),
			Email:     "user@test.com",
			FirstName: "User",
			LastName:  "Example",
			Role:      "user",
			Password:  "Test123456",
			Status:    models.StatusActive,
			Active:    true,
		}

		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", user, echo.New())

		err := ctl.User(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "username taken", "error message %s", "formatted")
		}
	})
}
