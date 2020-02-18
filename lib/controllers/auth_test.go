package controllers_test

import (
	"net/http"
	"net/http/httptest"
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
	_        = os.Setenv("AUTH_SECRET_KEY", "123")
	_        = os.Setenv("AUTH_ACCESS_TOKEN_LIFETIME", "123")
	_        = os.Setenv("AUTH_REFRESH_TOKEN_LIFETIME", "123")
	_authSrv = services.NewAuthService(mock.NewAuthRepository())
)

func TestControllers_Auth_NewAuthController(t *testing.T) {
	assert.NotNil(t, controllers.NewAuthController(_authSrv))
}

func TestControllers_Auth_Routes(t *testing.T) {
	e := echo.New()
	controllers.NewAuthController(_authSrv).Routes(e.Group("api"))

	c, _ := helpers.RequestTest(http.MethodPost, "/api/auth/token", e)

	assert.Equal(t, 400, c)
}

func TestControllers_Auth_TokenHandler(t *testing.T) {
	ctl := controllers.NewAuthController(_authSrv)

	t.Run("Wrong GrandType", func(t *testing.T) {
		body := models.AuthRequest{
			GrantType:    "something",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
			Username:     "peter@test.com",
			Password:     "testpass",
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("AUTHORISED", true)

		err := ctl.TokenHandler(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid grant type", "error message %s", "formatted")
		}
	})

	t.Run("Wrong Client ID", func(t *testing.T) {
		body := models.AuthRequest{
			GrantType:    "password",
			ClientID:     "Wrong",
			ClientSecret: "SecretSuper",
			Username:     "peter@test.com",
			Password:     "wrong-pass",
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("AUTHORISED", true)

		err := ctl.TokenHandler(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid client ID or secret", "error message %s", "formatted")
		}
	})

	t.Run("Empty Client ID", func(t *testing.T) {
		body := models.AuthRequest{
			GrantType:    "password",
			ClientSecret: "SecretSuper",
			Username:     "peter@test.com",
			Password:     "wrong-pass",
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("AUTHORISED", true)

		err := ctl.TokenHandler(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "client ID or secret cannot be empty", "error message %s", "formatted")
		}
	})

	t.Run("Wrong Client Secret", func(t *testing.T) {
		body := models.AuthRequest{
			GrantType:    "password",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "Wrong",
			Username:     "google@test.com",
			Password:     "testpass",
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("AUTHORISED", true)

		err := ctl.TokenHandler(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid client ID or secret", "error message %s", "formatted")
		}
	})

	t.Run("Empty Client Secret", func(t *testing.T) {
		body := models.AuthRequest{
			GrantType: "password",
			ClientID:  "SecRetAuthKey",
			Username:  "google@test.com",
			Password:  "testpass",
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("AUTHORISED", true)

		err := ctl.TokenHandler(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "client ID or secret cannot be empty", "error message %s", "formatted")
		}
	})
}

func TestControllers_Auth_TokenHandler_Password(t *testing.T) {
	ctl := controllers.NewAuthController(_authSrv)

	t.Run("Can login", func(t *testing.T) {
		body := models.AuthRequest{
			GrantType:    "password",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
			Username:     "peter@test.com",
			Password:     "testpass",
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("AUTHORISED", true)

		if assert.NoError(t, ctl.TokenHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Wrong password", func(t *testing.T) {
		body := models.AuthRequest{
			GrantType:    "password",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
			Username:     "peter@test.com",
			Password:     "wrong-pass",
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("AUTHORISED", true)

		err := ctl.TokenHandler(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid username or password", "error message %s", "formatted")
		}
	})

	t.Run("Wrong username", func(t *testing.T) {
		body := models.AuthRequest{
			GrantType:    "password",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
			Username:     "google@test.com",
			Password:     "testpass",
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("AUTHORISED", true)

		err := ctl.TokenHandler(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid username or password", "error message %s", "formatted")
		}
	})
}

func TestControllers_Auth_TokenHandler_RefreshToken(t *testing.T) {
	ctl := controllers.NewAuthController(_authSrv)

	t.Run("Can login", func(t *testing.T) {
		body := models.AuthRequest{
			GrantType:    "refresh_token",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
			RefreshToken: "sdfsdf5K9QwC6mptVSJVvAuFvA4w245HsiXxfMpOtpzASJ4Rr6E",
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("AUTHORISED", true)

		if assert.NoError(t, ctl.TokenHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Empty refresh token", func(t *testing.T) {
		body := models.AuthRequest{
			GrantType:    "refresh_token",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("AUTHORISED", true)

		err := ctl.TokenHandler(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "refresh token is empty or missing", "error message %s", "formatted")
		}
	})

	t.Run("Wrong refresh token", func(t *testing.T) {
		body := models.AuthRequest{
			GrantType:    "refresh_token",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
			RefreshToken: "wrong",
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("AUTHORISED", true)

		err := ctl.TokenHandler(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "refresh token not found", "error message %s", "formatted")
		}
	})

	t.Run("Expired refresh token", func(t *testing.T) {
		body := models.AuthRequest{
			GrantType:    "refresh_token",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
			RefreshToken: "ExpiredRefreshToken",
		}

		req := httptest.NewRequest(http.MethodPost, "/", helpers.ObjectToByte(t, body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set("AUTHORISED", true)

		err := ctl.TokenHandler(c)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "refresh token expired", "error message %s", "formatted")
		}
	})
}
