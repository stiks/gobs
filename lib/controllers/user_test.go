package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/pkg/helpers"
	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/controllers"
	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/services"
)

var (
	_cacheSrv = services.NewCacheService(mock.NewCacheRepository())
	_queueSrv = services.NewQueueService(mock.NewQueueRepository())
	_userSrv  = services.NewUserService(mock.NewUserRepository(), _queueSrv, _cacheSrv)
)

func TestControllers_User_NewUserController(t *testing.T) {
	assert.NotNil(t, controllers.NewUserController(_userSrv))
}

func TestControllers_User_List(t *testing.T) {
	ctl := controllers.NewUserController(_userSrv)

	t.Run("All users", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("USER_ID", "775a5b37-1742-4e54-9439-0357e768b011")

		err := ctl.List(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
}

func TestControllers_User_View(t *testing.T) {
	ctl := controllers.NewUserController(_userSrv)

	t.Run("Invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("USER_ID", "775a5b37-1742-4e54-9439-0357e768b011")

		err := ctl.View(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid UUID", "error message %s", "formatted")
		}
	})

	// TODO: Create view there
}

func TestControllers_User_Create(t *testing.T) {
	ctl := controllers.NewUserController(_userSrv)

	t.Run("Non-existing user", func(t *testing.T) {
		user := models.CreateUser{
			ID:        uuid.New(),
			FirstName: "John",
			LastName:  "Snow",
			Email:     "john@snow.com",
			Role:      "user",
			Password:  "Test123",
			Status:    models.StatusActive,
			Active:    true,
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, user))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("USER_ID", "775a5b37-1742-4e54-9439-0357e768b011")

		err := ctl.Create(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

	t.Run("Existing user", func(t *testing.T) {
		user := models.CreateUser{
			ID:        uuid.New(),
			Email:     "admin@test.com",
			FirstName: "Admin",
			LastName:  "Example",
			Role:      "user",
			Password:  "Test123",
			Status:    models.StatusActive,
			Active:    true,
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, user))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("USER_ID", "775a5b37-1742-4e54-9439-0357e768b011")

		err := ctl.Create(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "username taken", "error message %s", "formatted")
		}
	})
}

func TestControllers_User_Update(t *testing.T) {
	ctl := controllers.NewUserController(_userSrv)

	t.Run("Invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("USER_ID", "775a5b37-1742-4e54-9439-0357e768b011")

		err := ctl.Update(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid UUID", "error message %s", "formatted")
		}
	})

	// TODO: Update view there
}

func TestControllers_User_Delete(t *testing.T) {
	ctl := controllers.NewUserController(_userSrv)

	t.Run("Invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("USER_ID", "775a5b37-1742-4e54-9439-0357e768b011")

		err := ctl.Delete(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid UUID", "error message %s", "formatted")
		}
	})

	// TODO: Real test there
}
