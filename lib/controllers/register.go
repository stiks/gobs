package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/services"
	"github.com/stiks/gobs/pkg/xlog"
)

// RegisterControllerInterface ...
type RegisterControllerInterface interface {
	User(c echo.Context) error
	Routes(g *echo.Group)
}

type registerController struct {
	user services.UserService
}

// NewRegisterController returns a new Service instance
func NewRegisterController(userSrv services.UserService) RegisterControllerInterface {
	return &registerController{
		user: userSrv,
	}
}

// Routes registers routes
func (ctl *registerController) Routes(g *echo.Group) {
	g.POST("/register", ctl.User)
}

// ResetRequest ...
func (ctl *registerController) User(c echo.Context) error {
	ctx := c.Request().Context()

	xlog.Debugf(ctx, "New user registration request")

	req := new(models.CreateUser)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// set manually role to user
	req.Role = models.RoleUser
	req.Status = models.StatusInit

	if err := req.Validate(); err != nil {
		xlog.Errorf(ctx, "Unable to validate query, err: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Checking if users already exist
	if _, err := ctl.user.GetByUsername(ctx, req.Email); err == nil {
		xlog.Infof(ctx, "User already exist")

		return echo.NewHTTPError(http.StatusConflict, models.ErrUsernameTaken.Error())
	}

	user := req.ToUser(nil)
	user.ValidationHash = random.String(32, random.Alphanumeric)

	if _, err := ctl.user.Create(ctx, req.Password, user); err != nil {
		xlog.Errorf(ctx, "Unable to create user, err: %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": "ok"})
}
