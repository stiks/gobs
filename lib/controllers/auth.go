package controllers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/services"
	"github.com/stiks/gobs/pkg/xlog"
)

type authController struct {
	auth services.AuthService
}

// AuthControllerInterface ...
type AuthControllerInterface interface {
	Routes(g *echo.Group)
	TokenHandler(c echo.Context) error
}

// NewAuthController returns a new Service instance
func NewAuthController(authSrv services.AuthService) AuthControllerInterface {
	return &authController{
		auth: authSrv,
	}
}

// Routes registers routes
func (ctl *authController) Routes(g *echo.Group) {
	g.POST("/auth/token", ctl.TokenHandler)
}

// TokenHandler ...
func (ctl *authController) TokenHandler(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(models.AuthRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Map of grant types against handler functions
	grantTypes := map[string]func(ctx context.Context, r *models.AuthRequest, client *models.AuthClient) (*models.TokenResponse, error){
		"password":      ctl.auth.PasswordGrant,
		"refresh_token": ctl.auth.RefreshTokenGrant,
	}

	// Check the grant type
	grantHandler, ok := grantTypes[req.GrantType]
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, models.ErrInvalidGrantType.Error())
	}

	// Get auth client from request
	client, err := ctl.auth.GetClient(ctx, req)
	if err != nil {
		xlog.Infof(ctx, "Info: Trying to login with ClientID: %s, err: %s", req.ClientID, err.Error())

		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	xlog.Debugf(ctx, "User is trying to login with ClientID: %s and ClientSecret: %s", client.ClientID, client.ClientSecret)

	// Grant processing
	resp, err := grantHandler(ctx, req, client)
	if err != nil {
		xlog.Errorf(ctx, "Login error, %s", err.Error())
		xlog.Debugf(ctx, "Response, %+v", resp)

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	xlog.Infof(ctx, "User login successful")

	return c.JSON(http.StatusOK, resp)
}
