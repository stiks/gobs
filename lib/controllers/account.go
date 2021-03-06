package controllers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/services"
	"github.com/stiks/gobs/pkg/auth"
	"github.com/stiks/gobs/pkg/xlog"
)

// AccountControllerInterface ...
type AccountControllerInterface interface {
	PasswordConfirm(c echo.Context) error
	ResetRequest(c echo.Context) error
	EmailConfirm(c echo.Context) error
	GetProfile(c echo.Context) error
	Routes(g *echo.Group)
}

type accountController struct {
	user services.UserService
}

// NewAccountController returns a new Service instance
func NewAccountController(userSrv services.UserService) AccountControllerInterface {
	return &accountController{
		user: userSrv,
	}
}

// Routes registers routes
func (ctl *accountController) Routes(g *echo.Group) {
	g.POST("/account/reset-confirm", ctl.PasswordConfirm)
	g.POST("/account/reset", ctl.ResetRequest)
	g.POST("/account/email-confirm", ctl.EmailConfirm)
	g.GET("/account/profile", ctl.GetProfile, auth.EnableAuthorisation(), auth.RequiredAuth())
}

// GetProfile ...
func (ctl *accountController) GetProfile(c echo.Context) error {
	ctx := c.Request().Context()

	userID, err := auth.GetUserID(c)
	if err != nil {
		xlog.Errorf(ctx, "Unable to parse ID, %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, models.ErrInvalidUUID.Error())
	}

	user, err := ctl.user.GetByID(ctx, userID)
	if err != nil {
		xlog.Errorf(ctx, "Unable to find user, %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	xlog.Infof(ctx, "User %d login successful", userID)

	return c.JSON(http.StatusOK, user)
}

// ResetRequest ...
func (ctl *accountController) ResetRequest(c echo.Context) error {
	ctx := c.Request().Context()

	xlog.Debugf(ctx, "User requesting password reset")

	req := new(models.PasswordResetRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		xlog.Errorf(ctx, "Unable to validate query, err: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err := ctl.user.ResetPassword(ctx, req.Email)
	if err != nil {
		xlog.Debugf(ctx, "Password reset error: %s", err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

// PasswordConfirm ...
func (ctl *accountController) PasswordConfirm(c echo.Context) error {
	ctx := c.Request().Context()

	xlog.Debugf(ctx, "User going to password confirmation page")

	req := new(models.EmailConfirmationCode)
	if err := c.Bind(req); err != nil {
		xlog.Errorf(ctx, "Unable to bind, err: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		xlog.Errorf(ctx, "Unable to validate query, err: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := ctl.user.GetByResetHash(ctx, req.Code)
	if err != nil {
		xlog.Debugf(ctx, "Getting user by hash error: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, models.ErrEmailInvalidCode.Error())
	}

	// Cannot change password on locked account
	if user.Locked {
		xlog.Debugf(ctx, "user account is locked")

		return echo.NewHTTPError(http.StatusUnprocessableEntity, models.ErrUserIsLocked.Error())
	}

	// Forgot password code can be used only once
	if time.Now().Sub(user.PasswordResetAt) > time.Hour*24 {
		xlog.Debugf(ctx, "Forgot password code already used or expired")

		return echo.NewHTTPError(http.StatusUnprocessableEntity, models.ErrEmailCodeExpired.Error())
	}

	user, err = ctl.user.UpdatePassword(ctx, user.ID, req.Password)
	if err != nil {
		xlog.Debugf(ctx, "Unable update password, err: %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

// EmailConfirm ...
func (ctl *accountController) EmailConfirm(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(models.ConfirmEmail)
	if err := c.Bind(req); err != nil {
		xlog.Errorf(ctx, "Unable to bind, err: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		xlog.Errorf(ctx, "Unable to validate query, err: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := ctl.user.GetByID(ctx, req.UserID)
	if err != nil {
		xlog.Debugf(ctx, "Unable to fid user, error: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Cannot change password on locked account
	if user.Locked {
		xlog.Debugf(ctx, "user account is locked")

		return echo.NewHTTPError(http.StatusUnprocessableEntity, models.ErrUserIsLocked.Error())
	}

	if user.IsActive {
		xlog.Debugf(ctx, "user account already activated")

		return echo.NewHTTPError(http.StatusUnprocessableEntity, models.ErrEmailAlreadyConfirmed.Error())
	}

	if user.ValidationHash != req.Code {
		xlog.Debugf(ctx, "Incorrect code supplier")

		return echo.NewHTTPError(http.StatusUnprocessableEntity, models.ErrEmailConfirmationCode.Error())
	}

	user.IsActive = true

	if _, err := ctl.user.Update(ctx, user); err != nil {
		xlog.Debugf(ctx, "Getting user by hash error: %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
