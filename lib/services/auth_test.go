package services_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/services"
	"github.com/stiks/gobs/pkg/helpers"
)

func _authSrv() services.AuthService {
	os.Setenv("AUTH_SECRET_KEY", "123")
	os.Setenv("AUTH_ACCESS_TOKEN_LIFETIME", "123")
	os.Setenv("AUTH_REFRESH_TOKEN_LIFETIME", "123")

	return services.NewAuthService(mock.NewAuthRepository())
}

func TestService_Auth_NewAuthRepository(t *testing.T) {
	t.Run("AUTH_SECRET_KEY, not set", func(t *testing.T) {
		os.Unsetenv("AUTH_SECRET_KEY")
		os.Unsetenv("AUTH_ACCESS_TOKEN_LIFETIME")
		os.Unsetenv("AUTH_REFRESH_TOKEN_LIFETIME")

		os.Setenv("AUTH_ACCESS_TOKEN_LIFETIME", "123")
		os.Setenv("AUTH_REFRESH_TOKEN_LIFETIME", "123")

		assert.Panics(t, func() { services.NewAuthService(mock.NewAuthRepository()) })
	})

	t.Run("AUTH_ACCESS_TOKEN_LIFETIME", func(t *testing.T) {
		os.Unsetenv("AUTH_SECRET_KEY")
		os.Unsetenv("AUTH_ACCESS_TOKEN_LIFETIME")
		os.Unsetenv("AUTH_REFRESH_TOKEN_LIFETIME")

		os.Setenv("AUTH_SECRET_KEY", "123")
		os.Setenv("AUTH_REFRESH_TOKEN_LIFETIME", "123")

		assert.Panics(t, func() { services.NewAuthService(mock.NewAuthRepository()) })
	})

	t.Run("AUTH_REFRESH_TOKEN_LIFETIME", func(t *testing.T) {
		os.Unsetenv("AUTH_SECRET_KEY")
		os.Unsetenv("AUTH_ACCESS_TOKEN_LIFETIME")
		os.Unsetenv("AUTH_REFRESH_TOKEN_LIFETIME")

		os.Setenv("AUTH_SECRET_KEY", "123")
		os.Setenv("AUTH_ACCESS_TOKEN_LIFETIME", "123")

		assert.Panics(t, func() { services.NewAuthService(mock.NewAuthRepository()) })
	})

	t.Run("All set", func(t *testing.T) {
		assert.Implements(t, (*services.AuthService)(nil), _authSrv())
	})
}

func TestService_Auth_GetClient(t *testing.T) {
	srv := _authSrv()

	t.Run("Ok GradType, empty Client ID and Secret", func(t *testing.T) {
		_, err := srv.GetClient(nil, &models.AuthRequest{GrantType: "password"})
		if assert.Error(t, err) {
			assert.EqualError(t, err, "client ID or secret cannot be empty", "error message %s", "formatted")
		}
	})

	t.Run("Ok GradType, empty Client ID", func(t *testing.T) {
		_, err := srv.GetClient(nil, &models.AuthRequest{GrantType: "password", ClientID: "zzZzz"})
		if assert.Error(t, err) {
			assert.EqualError(t, err, "client ID or secret cannot be empty", "error message %s", "formatted")
		}
	})

	t.Run("Ok GradType, empty Client Secret", func(t *testing.T) {
		_, err := srv.GetClient(nil, &models.AuthRequest{GrantType: "password", ClientSecret: "zzZzz"})
		if assert.Error(t, err) {
			assert.EqualError(t, err, "client ID or secret cannot be empty", "error message %s", "formatted")
		}
	})

	t.Run("Should login", func(t *testing.T) {
		data := models.AuthRequest{
			GrantType:    "password",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
		}

		client, err := srv.GetClient(nil, &data)
		if assert.NoError(t, err) {
			assert.Equal(t, "775a5b37-1742-4e54-9439-0357e768b011", client.ID.String())
		}
	})

	t.Run("Wrong secret", func(t *testing.T) {
		data := models.AuthRequest{
			GrantType:    "password",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "test",
		}

		_, err := srv.GetClient(nil, &data)
		if assert.Error(t, err) {
			assert.EqualError(t, err, "invalid client ID or secret", "error message %s", "formatted")
		}
	})

	t.Run("Wrong client ID", func(t *testing.T) {
		data := models.AuthRequest{
			GrantType:    "password",
			ClientID:     "Nonexisting",
			ClientSecret: "SecretSuper",
		}

		_, err := srv.GetClient(nil, &data)
		if assert.Error(t, err) {
			assert.EqualError(t, err, "invalid client ID or secret", "error message %s", "formatted")
		}
	})
}

func TestService_Auth_PasswordGrant(t *testing.T) {
	srv := _authSrv()

	t.Run("All good", func(t *testing.T) {
		auth := models.AuthRequest{
			GrantType:    "password",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
			Username:     "peter@test.com",
			Password:     "testpass",
		}

		client := models.AuthClient{
			ID:           helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
		}

		token, err := srv.PasswordGrant(nil, &auth, &client)
		if assert.NoError(t, err) {
			assert.Equal(t, "775a5b37-1742-4e54-9439-0357e768b011", token.UserID.String())
		}
	})

	t.Run("Wrong username", func(t *testing.T) {
		auth := models.AuthRequest{
			GrantType:    "password",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
			Username:     "google@test.com",
			Password:     "testpass",
		}

		client := models.AuthClient{
			ID:           helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
		}

		_, err := srv.PasswordGrant(nil, &auth, &client)
		if assert.Error(t, err) {
			assert.EqualError(t, err, "invalid username or password", "error message %s", "formatted")
		}
	})

	t.Run("Wrong password", func(t *testing.T) {
		auth := models.AuthRequest{
			GrantType:    "password",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
			Username:     "peter@test.com",
			Password:     "wrong-pass",
		}

		client := models.AuthClient{
			ID:           helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
		}

		_, err := srv.PasswordGrant(nil, &auth, &client)
		if assert.Error(t, err) {
			assert.EqualError(t, err, "invalid username or password", "error message %s", "formatted")
		}
	})
}

func TestService_Auth_RefreshTokenGrant(t *testing.T) {
	srv := _authSrv()

	t.Run("All good, should login, with refresh token", func(t *testing.T) {
		auth := models.AuthRequest{
			GrantType:    "refresh_token",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
			RefreshToken: "sdfsdf5K9QwC6mptVSJVvAuFvA4w245HsiXxfMpOtpzASJ4Rr6E",
		}

		client := models.AuthClient{
			ID:           helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
		}

		token, err := srv.RefreshTokenGrant(nil, &auth, &client)
		if assert.NoError(t, err) {
			assert.Equal(t, "775a5b37-1742-4e54-9439-0357e768b011", token.UserID.String())
		}
	})

	t.Run("Empty refresh token", func(t *testing.T) {
		auth := models.AuthRequest{
			GrantType:    "refresh_token",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
		}

		client := models.AuthClient{
			ID:           helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
		}

		_, err := srv.RefreshTokenGrant(nil, &auth, &client)
		if assert.Error(t, err) {
			assert.EqualError(t, err, "refresh token is empty or missing", "error message %s", "formatted")
		}
	})

	t.Run("Wrong refresh token", func(t *testing.T) {
		auth := models.AuthRequest{
			GrantType:    "refresh_token",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
			RefreshToken: "wrong",
		}

		client := models.AuthClient{
			ID:           helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
		}

		_, err := srv.RefreshTokenGrant(nil, &auth, &client)
		if assert.Error(t, err) {
			assert.EqualError(t, err, "refresh token not found", "error message %s", "formatted")
		}
	})

	t.Run("Expired refresh token", func(t *testing.T) {
		auth := models.AuthRequest{
			GrantType:    "refresh_token",
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
			RefreshToken: "ExpiredRefreshToken",
		}

		client := models.AuthClient{
			ID:           helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
			ClientID:     "SecRetAuthKey",
			ClientSecret: "SecretSuper",
		}

		_, err := srv.RefreshTokenGrant(nil, &auth, &client)
		if assert.Error(t, err) {
			assert.EqualError(t, err, "refresh token expired", "error message %s", "formatted")
		}
	})
}
