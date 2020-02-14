package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/services"
	"github.com/stiks/gobs/pkg/auth"
	"github.com/stiks/gobs/pkg/xlog"
)

type userController struct {
	user services.UserService
}

// UserControllerInterface ...
type UserControllerInterface interface {
	View(c echo.Context) error
	List(c echo.Context) error
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	Routes(g *echo.Group)
}

// NewUserController ...
func NewUserController(service services.UserService) UserControllerInterface {
	return &userController{
		user: service,
	}
}

// Routes registers route handlers for the health service
func (ctl *userController) Routes(g *echo.Group) {
	g.Use(auth.EnableAuthorisation())

	g.GET("/users", ctl.List, auth.RequiredAuth())
	g.GET("/users/:id", ctl.View, auth.RequiredAuth())
	g.POST("/users", ctl.Create, auth.RequiredAuth())
	g.PUT("/users/:id", ctl.Update, auth.RequiredAuth())
	g.DELETE("/users/:id", ctl.Delete, auth.RequiredAuth())
}

// List ...
func (ctl *userController) List(c echo.Context) error {
	ctx := c.Request().Context()

	params := new(models.UserQueryParams)
	if err := c.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	users, err := ctl.user.GetAll(ctx, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	total, err := ctl.user.CountAll(ctx, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// hack to get non-empty list
	if len(users) <= 0 {
		users = []models.User{}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"list": users,
		"pagination": echo.Map{
			"current":  params.Page,
			"pageSize": params.PerPage,
			"total":    total,
		},
	})
}

// View ...
func (ctl *userController) View(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := ctl.user.GetByID(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.ErrUserNotFound.Error())
	}

	return c.JSON(http.StatusOK, user)
}

// Create ...
func (ctl *userController) Create(c echo.Context) error {
	ctx := c.Request().Context()

	u := new(models.CreateUser)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := u.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Checking if users already exist
	if _, err := ctl.user.GetByUsername(ctx, u.Email); err == nil {
		return echo.NewHTTPError(http.StatusConflict, models.ErrUsernameTaken.Error())
	}

	if len(u.Password) <= 0 {
		u.GeneratePassword()
	}

	user, err := ctl.user.Create(ctx, u.Password, &models.User{
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Status:    u.Status,
		Role:      u.Role,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}

// Update ...
func (ctl *userController) Update(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := ctl.user.GetByID(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.ErrUserNotFound.Error())
	}

	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := user.Validate(); err != nil {
		xlog.Errorf(ctx, "Unable to validate user query, err: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Overwrite request user
	user.ID = id

	user, err = ctl.user.Update(ctx, user)
	if err != nil {
		xlog.Errorf(ctx, "Unable to update user, err: %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusAccepted, user)
}

// Delete ...
func (ctl *userController) Delete(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = ctl.user.Delete(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}