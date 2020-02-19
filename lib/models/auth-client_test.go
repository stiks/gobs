package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/models"
)

func TestModel_AuthClient_ValidateSecret(t *testing.T) {
	data := &models.AuthClient{
		ClientID:     "zzZzz",
		ClientSecret: "ZzzZzz",
	}

	t.Run("Good password", func(t *testing.T) {
		assert.Equal(t, true, data.ValidateSecret("ZzzZzz"))
	})

	t.Run("Empty password", func(t *testing.T) {
		assert.Equal(t, false, data.ValidateSecret(""))
	})

	t.Run("Wrong password", func(t *testing.T) {
		assert.Equal(t, false, data.ValidateSecret("WrongPass"))
	})
}
