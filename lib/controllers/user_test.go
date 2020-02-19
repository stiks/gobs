package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/controllers"
	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/services"
	"github.com/stiks/gobs/pkg/helpers"
)

var (
	_cacheSrv = services.NewCacheService(mock.NewCacheRepository())
	_queueSrv = services.NewQueueService(mock.NewQueueRepository())
	_userSrv  = services.NewUserService(mock.NewUserRepository(), _queueSrv, _cacheSrv)
)

func TestControllers_User_NewUserController(t *testing.T) {
	assert.NotNil(t, controllers.NewUserController(_userSrv))
}

func TestControllers_User_Routes(t *testing.T) {
	t.Run("Get users", func(t *testing.T) {
		e := echo.New()
		controllers.NewUserController(_userSrv).Routes(e.Group("api"))

		c, _ := helpers.RequestTest(http.MethodGet, "/api/users", e)
		assert.Equal(t, 400, c)
	})

	t.Run("View user", func(t *testing.T) {
		e := echo.New()
		controllers.NewUserController(_userSrv).Routes(e.Group("api"))

		c, _ := helpers.RequestTest(http.MethodGet, "/api/users/123123", e)
		assert.Equal(t, 400, c)
	})

	t.Run("Create user", func(t *testing.T) {
		e := echo.New()
		controllers.NewUserController(_userSrv).Routes(e.Group("api"))

		c, _ := helpers.RequestTest(http.MethodPost, "/api/users", e)
		assert.Equal(t, 400, c)
	})

	t.Run("Update user", func(t *testing.T) {
		e := echo.New()
		controllers.NewUserController(_userSrv).Routes(e.Group("api"))

		c, _ := helpers.RequestTest(http.MethodPost, "/api/users", e)
		assert.Equal(t, 400, c)
	})

	t.Run("Delete user", func(t *testing.T) {
		e := echo.New()
		controllers.NewUserController(_userSrv).Routes(e.Group("api"))

		c, _ := helpers.RequestTest(http.MethodDelete, "/api/users/123", e)
		assert.Equal(t, 400, c)
	})
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
			assert.Contains(t, rec.Body.String(), "peter@test.com")
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

	t.Run("Existing user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		e := echo.New()
		e.SetMaxParam(2)

		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("775a5b37-1742-4e54-9439-0357e768b011")

		if assert.NoError(t, ctl.View(c)) {
			assert.Contains(t, rec.Body.String(), "peter@test.com")
		}
	})

	t.Run("Non existing user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		e := echo.New()
		e.SetMaxParam(2)

		c := e.NewContext(req, rec)
		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("5fcc94e5-c6aa-4320-8469-f5021af54b88")

		err := ctl.View(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "user not found", "error message %s", "formatted")
		}
	})
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
			Password:  "Test123456",
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
			Password:  "Test123456",
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

	t.Run("Existing user", func(t *testing.T) {
		body := models.CreateUser{
			FirstName: "John",
			LastName:  "Snow",
			Email:     "john@snow.com",
			Password:  "testpass",
			Role:      "user",
			Status:    2,
		}

		req := httptest.NewRequest(http.MethodDelete, "/", helpers.ObjectToByte(t, body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		e := echo.New()
		e.SetMaxParam(2)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("775a5b37-1742-4e54-9439-0357e768b011")

		if assert.NoError(t, ctl.Update(c)) {
			assert.Contains(t, rec.Body.String(), "peter@test.com")
		}
	})

	t.Run("Non existing user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		e := echo.New()
		e.SetMaxParam(2)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("5fcc94e5-c6aa-4320-8469-f5021af54b88")

		err := ctl.Update(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "user not found", "error message %s", "formatted")
		}
	})
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

	t.Run("Existing user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		e := echo.New()
		e.SetMaxParam(2)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("775a5b37-1742-4e54-9439-0357e768b011")

		assert.NoError(t, ctl.Delete(c))
	})

	t.Run("Non existing user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		e := echo.New()
		e.SetMaxParam(2)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath("/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("5fcc94e5-c6aa-4320-8469-f5021af54b88")

		err := ctl.Delete(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "user not found", "error message %s", "formatted")
		}
	})
}
