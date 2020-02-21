package controllers_test

import (
	"net/http"
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

func TestControllers_Account_Routes(t *testing.T) {
	t.Run("Get profile", func(t *testing.T) {
		e := echo.New()
		controllers.NewAccountController(_userSrv).Routes(e.Group("api"))

		c, _ := helpers.RequestTest(http.MethodGet, "/api/account/profile", e)
		assert.Equal(t, 400, c)
	})

	t.Run("Post reset confirmation", func(t *testing.T) {
		e := echo.New()
		controllers.NewAccountController(_userSrv).Routes(e.Group("api"))

		c, _ := helpers.RequestTest(http.MethodPost, "/api/account/reset-confirm", e)
		assert.Equal(t, 400, c)
	})

	t.Run("Post reset", func(t *testing.T) {
		e := echo.New()
		controllers.NewAccountController(_userSrv).Routes(e.Group("api"))

		c, _ := helpers.RequestTest(http.MethodPost, "/api/account/reset", e)
		assert.Equal(t, 400, c)
	})
}

func TestControllers_Account_GetProfile(t *testing.T) {
	ctl := controllers.NewAccountController(_userSrv)

	t.Run("Invalid UUID", func(t *testing.T) {
		_, ctx := helpers.RequestWithBody(http.MethodGet, "/", nil, echo.New())

		err := ctl.GetProfile(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid UUID", "error message %s", "formatted")
		}
	})

	t.Run("Non-existing user", func(t *testing.T) {
		_, ctx := helpers.RequestWithBody(http.MethodGet, "/", nil, echo.New())
		ctx.Set("USER_ID", "921c3683-e8e6-41fd-8adb-cdb54429ad51")

		err := ctl.GetProfile(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "user not found", "error message %s", "formatted")
		}
	})

	t.Run("Existing user", func(t *testing.T) {
		rec, ctx := helpers.RequestWithBody(http.MethodGet, "/", nil, echo.New())
		ctx.Set("USER_ID", "775a5b37-1742-4e54-9439-0357e768b011")

		err := ctl.GetProfile(ctx)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "peter@test.com")
		}
	})
}

func TestControllers_Account_ResetRequest(t *testing.T) {
	ctl := controllers.NewAccountController(_userSrv)

	t.Run("Non-existing user", func(t *testing.T) {
		rec, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", models.PasswordResetRequest{Email: "test@google.com"}, echo.New())

		err := ctl.ResetRequest(ctx)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "ok")
		}
	})

	t.Run("Existing user", func(t *testing.T) {
		rec, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", models.PasswordResetRequest{Email: "peter@test.com"}, echo.New())

		err := ctl.ResetRequest(ctx)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "ok")
		}
	})
}

func TestControllers_Account_PasswordConfirm(t *testing.T) {
	ctl := controllers.NewAccountController(_userSrv)

	t.Run("Blank code", func(t *testing.T) {
		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", models.PasswordResetRequest{Email: "test@google.com"}, echo.New())

		err := ctl.PasswordConfirm(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "email confirmation code cannot be blank")
		}
	})

	t.Run("Non-existing user", func(t *testing.T) {
		data := models.EmailConfirmationCode{
			Code:     "zzZzzZzz",
			Password: "ASdkjnw3rwdf234sdf",
		}

		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", data, echo.New())

		err := ctl.PasswordConfirm(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "invalid email confirmation code supplied")
		}
	})

	t.Run("Locked account", func(t *testing.T) {
		data := models.EmailConfirmationCode{
			Code:     "2e4EHSsVkledZxWwU7j3BnNBYo",
			Password: "R0otIsG0od",
		}

		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", data, echo.New())

		err := ctl.PasswordConfirm(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "user account is locked")
		}
	})

	t.Run("Expired token", func(t *testing.T) {
		data := models.EmailConfirmationCode{
			Code:     "5zQVfk8aQlZgQiW0vd2PA8kyj4",
			Password: "R0otIsG0od",
		}

		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", data, echo.New())

		err := ctl.PasswordConfirm(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "email confirmation code already used or expired")
		}
	})

	t.Run("Existing code, blank password", func(t *testing.T) {
		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", models.EmailConfirmationCode{Code: "ZXqEMubf5DinaTHuOyJIm1z3Dq"}, echo.New())

		err := ctl.PasswordConfirm(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "password: cannot be blank")
		}
	})

	t.Run("Existing code, short password", func(t *testing.T) {
		data := models.EmailConfirmationCode{
			Code:     "ZXqEMubf5DinaTHuOyJIm1z3Dq",
			Password: "123",
		}

		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", data, echo.New())

		err := ctl.PasswordConfirm(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "password: the length must be between 8")
		}
	})

	t.Run("Existing code, too long password", func(t *testing.T) {
		data := models.EmailConfirmationCode{
			Code:     "ZXqEMubf5DinaTHuOyJIm1z3Dq",
			Password: "wFbxjwfIEVjTq7YbIGdw0d4u07wFbxjwfIEVjTq7YbIGdw0d4u07wFbxjwfIEVjTq7YbIGdw0d4u07wFbxjwfIEVjTq7YbIGdw0d4u07",
		}

		_, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", data, echo.New())

		err := ctl.PasswordConfirm(ctx)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "password: the length must be between 8")
		}
	})

	t.Run("Existing code, good password", func(t *testing.T) {
		data := models.EmailConfirmationCode{
			Code:     "ZXqEMubf5DinaTHuOyJIm1z3Dq",
			Password: "wFbxjwfIEVjTq7YbIGdw0d4u07",
		}

		rec, ctx := helpers.RequestObjectWithBody(t, http.MethodPost, "/", data, echo.New())

		err := ctl.PasswordConfirm(ctx)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "ok")
		}
	})
}
