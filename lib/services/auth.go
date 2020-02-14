package services

import (
	"context"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/stiks/gobs/pkg/xlog"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/repositories"
	"github.com/stiks/gobs/pkg/env"
)

type authService struct {
	repo                 repositories.AuthRepository
	AccessTokenLifetime  int
	RefreshTokenLifetime int
	JWTSecretCode        []byte
}

// AuthService ...
type AuthService interface {
	RefreshTokenGrant(ctx context.Context, r *models.AuthRequest, client *models.AuthClient) (*models.TokenResponse, error)
	PasswordGrant(ctx context.Context, r *models.AuthRequest, client *models.AuthClient) (*models.TokenResponse, error)
	GetClient(ctx context.Context, r *models.AuthRequest) (*models.AuthClient, error)
}

// NewAuthService ...
func NewAuthService(repo repositories.AuthRepository) AuthService {
	if !env.MustPresent("AUTH_SECRET_KEY") {
		log.Fatalf("'AUTH_SECRET_KEY' must be set")
	}

	if !env.MustPresent("AUTH_ACCESS_TOKEN_LIFETIME") {
		log.Fatalf("'AUTH_ACCESS_TOKEN_LIFETIME' must be set")
	}

	if !env.MustPresent("AUTH_REFRESH_TOKEN_LIFETIME") {
		log.Fatalf("'AUTH_REFRESH_TOKEN_LIFETIME' must be set")
	}

	return &authService{
		JWTSecretCode:        []byte(env.MustGetString("AUTH_SECRET_KEY")),
		AccessTokenLifetime:  env.MustGetInt("AUTH_ACCESS_TOKEN_LIFETIME"),
		RefreshTokenLifetime: env.MustGetInt("AUTH_REFRESH_TOKEN_LIFETIME"),
		repo:                 repo,
	}
}

// GetClient from request
func (s *authService) GetClient(ctx context.Context, r *models.AuthRequest) (*models.AuthClient, error) {
	// Get client credentials from request
	if r.ClientID == "" || r.ClientSecret == "" {
		return nil, models.ErrEmptyClientOrSecret
	}

	// Looking for the client by Client ID
	client, err := s.repo.FindByClientID(ctx, r.ClientID)
	if err != nil {
		xlog.Errorf(ctx, "Error getting Client ID: %s", err.Error())

		// For security reasons, return a general error message
		return nil, models.ErrInvalidClientOrSecret
	}

	// Validate client secret
	if !client.ValidateSecret(r.ClientSecret) {
		xlog.Errorf(ctx, "Client secret is invalid")

		// For security reasons, return a general error message
		return nil, models.ErrInvalidClientOrSecret
	}

	return client, nil
}

// PasswordGrant login using password
func (s *authService) PasswordGrant(ctx context.Context, req *models.AuthRequest, client *models.AuthClient) (*models.TokenResponse, error) {
	// Find user by username
	user, err := s.repo.FindUserByUsername(ctx, req.Username)
	if err != nil && err == models.ErrUserNotFound {
		// For security reason
		return nil, models.ErrInvalidUsernameOrPassword
	}

	// Any other get user error
	if err != nil {
		xlog.Errorf(ctx, "Find user details error: %s", err.Error())

		return nil, err
	}

	// Check that the password is set
	if len(user.PasswordHash) <= 0 {
		xlog.Errorf(ctx, "User password hash field is empty")

		return nil, models.ErrUserPasswordNotSet
	}

	// Verify the password
	if !user.ValidatePassword(req.Password) {
		xlog.Errorf(ctx, "User password is wrong")

		return nil, models.ErrInvalidUsernameOrPassword
	}

	// create a new access token
	accessToken, err := models.NewAccessToken(client, user, s.AccessTokenLifetime, s.JWTSecretCode)
	if err != nil {
		xlog.Errorf(ctx, "Unable to create access token, err: %s", err.Error())

		return nil, err
	}

	// create or retrieve a refresh token
	refreshToken, err := s.getOrCreateRefreshToken(ctx, client, user)
	if err != nil {
		xlog.Errorf(ctx, "Unable to create or get refresh token, err: %s", err.Error())

		return nil, err
	}

	if err := s.repo.UpdateLastLogin(ctx, user.ID); err != nil {
		xlog.Errorf(ctx, "Unable to set users last login, err: %s", err.Error())
	}

	// create response
	return models.NewTokenResponse(accessToken, refreshToken, s.AccessTokenLifetime, "Bearer")
}

// getOrCreateRefreshToken retrieves an existing refresh token, if expired,
// the token gets deleted and new refresh token is created
func (s *authService) getOrCreateRefreshToken(ctx context.Context, client *models.AuthClient, user *models.User) (*models.Token, error) {
	// Try to fetch an existing refresh token first
	refreshToken, err := s.repo.FindByClientUser(ctx, client.ID, user.ID)
	if err != nil {
		xlog.Errorf(ctx, "Unable to find token, err: %s", err.Error())

		// We assume token already expired
		return s.generateNewRefreshToken(ctx, client, user)
	}

	// If the refresh token has expired, delete it
	if time.Now().UTC().After(time.Unix(int64(refreshToken.ExpiresAt), 0)) {
		xlog.Errorf(ctx, "Token %d expired, deleting", refreshToken.ID)

		if err := s.repo.DeleteToken(ctx, refreshToken.ID); err != nil {
			xlog.Errorf(ctx, "Unable delete token %d, err: %s", refreshToken.ID, err.Error())
		}

		return s.generateNewRefreshToken(ctx, client, user)
	}

	// All other cases, we just return token
	return refreshToken, nil
}

// generateNewToken generates new token
func (s *authService) generateNewRefreshToken(ctx context.Context, client *models.AuthClient, user *models.User) (*models.Token, error) {
	// We assume token already expired
	refreshToken, err := s.repo.CreateToken(ctx, models.NewRefreshToken(client, user, s.RefreshTokenLifetime))
	if err != nil {
		xlog.Errorf(ctx, "Unable to create token, err: %s", err.Error())

		return nil, err
	}

	refreshToken.Client = client
	refreshToken.User = user

	return refreshToken, nil
}

func (s *authService) RefreshTokenGrant(ctx context.Context, r *models.AuthRequest, client *models.AuthClient) (*models.TokenResponse, error) {
	// Fetch the refresh token
	refreshToken, err := s.getValidRefreshToken(ctx, r.RefreshToken, client)
	if err != nil {
		return nil, err
	}

	// Find user by User ID
	user, err := s.repo.FindUserByID(ctx, refreshToken.UserID)
	if err != nil {
		xlog.Errorf(ctx, "User not found, err: %s", err.Error())

		return nil, err
	}

	// create a new access token
	accessToken, err := models.NewAccessToken(client, user, s.AccessTokenLifetime, s.JWTSecretCode)
	if err != nil {
		xlog.Errorf(ctx, "Unable to create access token, err: %s", err.Error())

		return nil, err
	}

	// create response
	return models.NewTokenResponse(accessToken, refreshToken, s.AccessTokenLifetime, "Bearer")
}

// getValidRefreshToken returns a valid non expired refresh token
func (s *authService) getValidRefreshToken(ctx context.Context, token string, client *models.AuthClient) (*models.Token, error) {
	// Fetch the refresh token from the database
	refreshToken, err := s.repo.FindByHashClient(ctx, client.ID, token)
	if err != nil {
		xlog.Errorf(ctx, "Unable to find client by Client ID, err: %s", err.Error())

		return nil, models.ErrRefreshTokenNotFound
	}

	// Check the refresh token hasn't expired
	if time.Now().UTC().After(time.Unix(int64(refreshToken.ExpiresAt), 0)) {
		xlog.Errorf(ctx, "Token already expired")

		return nil, models.ErrRefreshTokenExpired
	}

	return refreshToken, nil
}
