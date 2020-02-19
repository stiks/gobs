package models_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/models"
)

func TestModel_Token_NewTokenResponse(t *testing.T) {
	userID := uuid.New()

	accessToken := models.Token{
		UserID:    userID,
		Token:     "AccessToken",
		ExpiresAt: 1,
		User: &models.User{
			Role: models.RoleUser,
		},
	}

	refreshToken := models.Token{
		Token: "RefreshTokenZzz",
	}

	t.Run("Good token", func(t *testing.T) {
		token, err := models.NewTokenResponse(&accessToken, &refreshToken, 1, "Bearer")
		if assert.NoError(t, err) {
			assert.Equal(t, "AccessToken", token.AccessToken)
			assert.Equal(t, "RefreshTokenZzz", token.RefreshToken)
			assert.Equal(t, 1, token.ExpiresIn)
			assert.Equal(t, userID, token.UserID)
			assert.Equal(t, models.RoleUser, token.Authority)
		}
	})

	t.Run("Good token with no refresh token", func(t *testing.T) {
		token, err := models.NewTokenResponse(&accessToken, nil, 1, "Bearer")
		if assert.NoError(t, err) {
			assert.Equal(t, "AccessToken", token.AccessToken)
			assert.Empty(t, token.RefreshToken)
			assert.Equal(t, 1, token.ExpiresIn)
			assert.Equal(t, userID, token.UserID)
			assert.Equal(t, models.RoleUser, token.Authority)
		}
	})
}

func TestModel_Token_NewAccessToken(t *testing.T) {
	client := &models.AuthClient{
		ID: uuid.New(),
	}

	user := &models.User{
		ID: uuid.New(),
	}

	t.Run("Good token", func(t *testing.T) {
		token, err := models.NewAccessToken(client, user, 1, []byte("something"))
		if assert.NoError(t, err) {
			assert.NotEmpty(t, token.Token)
			assert.Equal(t, client.ID, token.ClientID)
			assert.Equal(t, user.ID, token.UserID)
		}
	})
}

func TestModel_Token_NewRefreshToken(t *testing.T) {
	client := &models.AuthClient{
		ID: uuid.New(),
	}

	user := &models.User{
		ID: uuid.New(),
	}

	t.Run("Good token", func(t *testing.T) {
		token := models.NewRefreshToken(client, user, 1)

		assert.NotEmpty(t, token.Token)
		assert.Equal(t, client.ID, token.ClientID)
		assert.Equal(t, user.ID, token.UserID)
	})
}
