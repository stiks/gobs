package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/controllers"
	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/pkg/helpers"
)

func TestControllers_Account_NewAccountController(t *testing.T) {
	assert.NotNil(t, controllers.NewAccountController(_userSrv))
}

func TestControllers_Account_GetProfile(t *testing.T) {
	ctl := controllers.NewAccountController(_userSrv)

	t.Run("Invalid UUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := ctl.GetProfile(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid UUID", "error message %s", "formatted")
		}
	})

	t.Run("Non-existing user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("USER_ID", "921c3683-e8e6-41fd-8adb-cdb54429ad51")

		err := ctl.GetProfile(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "user not found", "error message %s", "formatted")
		}
	})

	t.Run("Existing user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("USER_ID", "775a5b37-1742-4e54-9439-0357e768b011")

		err := ctl.GetProfile(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "peter@test.com")
		}
	})
}

func TestControllers_Account_ResetRequest(t *testing.T) {
	ctl := controllers.NewAccountController(_userSrv)

	t.Run("Non-existing user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, models.PasswordResetRequest{Email: "test@google.com"}))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := ctl.ResetRequest(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "ok")
		}
	})

	t.Run("Existing user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, models.PasswordResetRequest{Email: "peter@test.com"}))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := ctl.ResetRequest(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "ok")
		}
	})
}

func TestControllers_Account_PasswordConfirm(t *testing.T) {
	ctl := controllers.NewAccountController(_userSrv)

	t.Run("Blank code", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, models.PasswordResetRequest{Email: "test@google.com"}))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := ctl.PasswordConfirm(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "email confirmation code cannot be blank")
		}
	})

	t.Run("Non-existing user", func(t *testing.T) {
		data := models.EmailConfirmationCode{
			Code:     "zzZzzZzz",
			Password: "ASdkjnw3rwdf234sdf",
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := ctl.PasswordConfirm(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid email confirmation code supplied")
		}
	})

	t.Run("Existing code, blank password", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, models.EmailConfirmationCode{Code: "ZXqEMubf5DinaTHuOyJIm1z3Dq"}))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := ctl.PasswordConfirm(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "password: cannot be blank")
		}
	})

	t.Run("Existing code, short password", func(t *testing.T) {
		data := models.EmailConfirmationCode{
			Code:     "ZXqEMubf5DinaTHuOyJIm1z3Dq",
			Password: "123",
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := ctl.PasswordConfirm(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "password: the length must be between 8")
		}
	})

	t.Run("Existing code, too long password", func(t *testing.T) {
		data := models.EmailConfirmationCode{
			Code:     "ZXqEMubf5DinaTHuOyJIm1z3Dq",
			Password: "wFbxjwfIEVjTq7YbIGdw0d4u07wFbxjwfIEVjTq7YbIGdw0d4u07wFbxjwfIEVjTq7YbIGdw0d4u07wFbxjwfIEVjTq7YbIGdw0d4u07",
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := ctl.PasswordConfirm(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "password: the length must be between 8")
		}
	})

	t.Run("Existing code, good password", func(t *testing.T) {
		data := models.EmailConfirmationCode{
			Code:     "ZXqEMubf5DinaTHuOyJIm1z3Dq",
			Password: "wFbxjwfIEVjTq7YbIGdw0d4u07",
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, data))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := ctl.PasswordConfirm(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "ok")
		}
	})
}
