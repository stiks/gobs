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

		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", body, echo.New())
		ctx.Set("AUTHORISED", true)

		err := ctl.TokenHandler(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "grant_type: must be a valid value.", "error message %s", "formatted")
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

		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", body, echo.New())
		ctx.Set("AUTHORISED", true)

		err := ctl.TokenHandler(ctx)
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

		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", body, echo.New())
		ctx.Set("AUTHORISED", true)

		err := ctl.TokenHandler(ctx)
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

		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", body, echo.New())
		ctx.Set("AUTHORISED", true)

		err := ctl.TokenHandler(ctx)
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

		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", body, echo.New())
		ctx.Set("AUTHORISED", true)

		err := ctl.TokenHandler(ctx)
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

		rec, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", body, echo.New())
		ctx.Set("AUTHORISED", true)

		if assert.NoError(t, ctl.TokenHandler(ctx)) {
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

		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", body, echo.New())
		ctx.Set("AUTHORISED", true)

		err := ctl.TokenHandler(ctx)
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

		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", body, echo.New())
		ctx.Set("AUTHORISED", true)

		err := ctl.TokenHandler(ctx)
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

		rec, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", body, echo.New())
		ctx.Set("AUTHORISED", true)

		if assert.NoError(t, ctl.TokenHandler(ctx)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

	t.Run("Empty refresh token", func(t *testing.T) {
		body := models.AuthRequest{
			GrantType:    "refresh_token",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
		}

		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", body, echo.New())
		ctx.Set("AUTHORISED", true)

		err := ctl.TokenHandler(ctx)
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

		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", body, echo.New())
		ctx.Set("AUTHORISED", true)

		err := ctl.TokenHandler(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "refresh token not found", "error message %s", "formatted")
		}
	})

	t.Run("Expired refresh token", func(t *testing.T) {
		body := models.AuthRequest{
			GrantType:    "refresh_token",
			ClientID:     "RandomStuffHere",
			ClientSecret: "RandomKeySecret",
			RefreshToken: "ExpiredRefreshToken",
		}

		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", body, echo.New())
		ctx.Set("AUTHORISED", true)

		err := ctl.TokenHandler(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "refresh token expired", "error message %s", "formatted")
		}
	})
}
