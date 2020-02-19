package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/matcornic/hermes/v2"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/services"
	"github.com/stiks/gobs/pkg/env"
	"github.com/stiks/gobs/pkg/xlog"
)

type workerController struct {
	queue services.QueueService
	user  services.UserService
	email services.EmailService
}

// WorkerControllerInterface ...
type WorkerControllerInterface interface {
	Routes(g *echo.Group)
	UserPasswordReset(c echo.Context) error
	UserProfileUpdated(c echo.Context) error
	UserPasswordChanged(c echo.Context) error
}

// NewWorkerController returns a controller
func NewWorkerController(userSrv services.UserService, queueSrv services.QueueService, emailSrv services.EmailService) WorkerControllerInterface {
	return &workerController{
		user:  userSrv,
		queue: queueSrv,
		email: emailSrv,
	}
}

// Routes registers routes
func (ctl *workerController) Routes(g *echo.Group) {
	g.POST("/user-password-reset", ctl.UserPasswordReset)
	g.POST("/user-profile-updated", ctl.UserProfileUpdated)
	g.POST("/user-password-changed", ctl.UserPasswordChanged)
}

// UserPasswordReset ...
func (ctl *workerController) UserPasswordReset(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(models.WorkerRequest)
	if err := c.Bind(req); err != nil {
		xlog.Errorf(ctx, "Bind error: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := ctl.user.GetByID(ctx, req.ID)
	if err != nil {
		xlog.Errorf(ctx, "Unable to find user, err: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	msg := hermes.Email{
		Body: hermes.Body{
			Name: fmt.Sprintf("%s", user.FirstName),
			Intros: []string{
				fmt.Sprintf("You have received this email because a password reset request for %s account was received.", env.MustGetString("PUBLIC_NAME")),
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to reset your password:",
					Button: hermes.Button{
						Color: "#DC4D2F",
						Text:  "Reset your password",
						Link:  fmt.Sprintf("%s/user/reset-password/%s", env.MustGetString("PUBLIC_HOSTNAME"), user.PasswordResetHash),
					},
				},
			},
			Outros: []string{
				"If you did not request a password reset, no further action is required on your part.",
			},
			Signature: "Thanks",
		},
	}

	if err := ctl.email.SendEmail(ctx, user.Email, "Password Recovery", msg); err != nil {
		xlog.Errorf(ctx, "Unable to send email, err: %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// UserProfileUpdated ...
func (ctl *workerController) UserProfileUpdated(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(models.WorkerRequest)
	if err := c.Bind(req); err != nil {
		xlog.Errorf(ctx, "Bind error: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := ctl.user.GetByID(ctx, req.ID)
	if err != nil {
		xlog.Errorf(ctx, "Unable to find user, err: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	msg := hermes.Email{
		Body: hermes.Body{
			Name:      fmt.Sprintf("%s", user.FirstName),
			Intros:    []string{"You have successfully changed your profile."},
			Signature: "Thanks",
		},
	}

	if err := ctl.email.SendEmail(ctx, user.Email, "Profile updated", msg); err != nil {
		xlog.Errorf(ctx, "Unable to send email, err: %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

// UserPasswordChanged ...
func (ctl *workerController) UserPasswordChanged(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(models.WorkerRequest)
	if err := c.Bind(req); err != nil {
		xlog.Errorf(ctx, "Bind error: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := ctl.user.GetByID(ctx, req.ID)
	if err != nil {
		xlog.Errorf(ctx, "Unable to find user, err: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	msg := hermes.Email{
		Body: hermes.Body{
			Name:      fmt.Sprintf("%s", user.FirstName),
			Intros:    []string{"You have successfully changed your password."},
			Signature: "Thanks",
		},
	}

	if err := ctl.email.SendEmail(ctx, user.Email, "Password changed successfully", msg); err != nil {
		xlog.Errorf(ctx, "Unable to send email, err: %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
