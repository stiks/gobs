package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/models"
)

func TestModel_AuthRequest_Validate(t *testing.T) {
	t.Run("Grant password", func(t *testing.T) {
		data := &models.AuthRequest{
			ClientID:     "zzZzz",
			ClientSecret: "ZzzZzz",
			GrantType:    "password",
			Username:     "john@snow.com",
			Password:     "testpass",
		}

		assert.NoError(t, data.Validate())
	})

	t.Run("Grant refresh token", func(t *testing.T) {
		data := &models.AuthRequest{
			ClientID:     "zzZzz",
			ClientSecret: "ZzzZzz",
			GrantType:    "refresh_token",
		}

		assert.NoError(t, data.Validate())
	})

	t.Run("Wrong grant", func(t *testing.T) {
		data := &models.AuthRequest{
			ClientID:     "zzZzz",
			ClientSecret: "ZzzZzz",
			GrantType:    "something",
		}

		err := data.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "grant_type: must be a valid value.", "error message %s", "formatted")
		}
	})

	t.Run("Empty grant", func(t *testing.T) {
		data := &models.AuthRequest{
			ClientID:     "zzZzz",
			ClientSecret: "ZzzZzz",
			Username:     "john@snow.com",
			Password:     "testpass",
		}

		err := data.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "grant_type: cannot be blank.", "error message %s", "formatted")
		}
	})
}
