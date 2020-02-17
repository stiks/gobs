package mock_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/repositories"
	"github.com/stiks/gobs/pkg/helpers"
)

func TestMock_Auth_NewAuthRepository(t *testing.T) {
	r := mock.NewAuthRepository()

	assert.Implements(t, (*repositories.AuthRepository)(nil), r)
}

func TestMock_Auth_FindByClientID(t *testing.T) {
	r := mock.NewAuthRepository()

	t.Run("Existing auth client", func(t *testing.T) {
		authClient, err := r.FindByClientID(nil, "SecRetAuthKey")
		if assert.NoError(t, err) {
			assert.Equal(t, "775a5b37-1742-4e54-9439-0357e768b011", authClient.ID.String())
		}
	})

	t.Run("Non-existing auth client", func(t *testing.T) {
		_, err := r.FindByClientID(nil, "NonexstingOne")
		if assert.Error(t, err) {
			assert.EqualError(t, err, "auth client could not be found", "error message %s", "formatted")
		}
	})
}

func TestMock_Auth_FindUserByUsername(t *testing.T) {
	r := mock.NewAuthRepository()

	t.Run("Existing user", func(t *testing.T) {
		user, err := r.FindUserByUsername(nil, "peter@test.com")
		if assert.NoError(t, err) {
			assert.Equal(t, "775a5b37-1742-4e54-9439-0357e768b011", user.ID.String())
		}
	})

	t.Run("Non-existing user", func(t *testing.T) {
		_, err := r.FindUserByUsername(nil, "anything@test.com")
		if assert.Error(t, err) {
			assert.EqualError(t, err, "user not found", "error message %s", "formatted")
		}
	})
}

func TestMock_Auth_FindByID(t *testing.T) {
	r := mock.NewAuthRepository()

	t.Run("Existing user", func(t *testing.T) {
		user, err := r.FindUserByID(nil, helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011"))
		if assert.NoError(t, err) {
			assert.Equal(t, "peter@test.com", user.Email)
		}
	})

	t.Run("Non-existing user", func(t *testing.T) {
		_, err := r.FindUserByID(nil, helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88"))
		if assert.Error(t, err) {
			assert.EqualError(t, err, "user not found", "error message %s", "formatted")
		}
	})
}

func TestMock_Auth_UpdateLastLogin(t *testing.T) {
	r := mock.NewAuthRepository()

	t.Run("Existing user", func(t *testing.T) {
		assert.NoError(t, r.UpdateLastLogin(nil, helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011")))
	})

	t.Run("Non-existing user", func(t *testing.T) {
		err := r.UpdateLastLogin(nil, helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88"))
		if assert.Error(t, err) {
			assert.EqualError(t, err, "user not found", "error message %s", "formatted")
		}
	})
}

func TestMock_Auth_FindByClientUser(t *testing.T) {
	r := mock.NewAuthRepository()

	t.Run("Existing token", func(t *testing.T) {
		token, err := r.FindByClientUser(nil, helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011"), helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011"))
		if assert.NoError(t, err) {
			assert.Equal(t, "775a5b37-1742-4e54-9439-0357e768b011", token.ID.String())
		}
	})

	t.Run("Non-existing token", func(t *testing.T) {
		_, err := r.FindByClientUser(nil, helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88"), helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011"))
		if assert.Error(t, err) {
			assert.EqualError(t, err, "token not found", "error message %s", "formatted")
		}
	})
}

func TestMock_Auth_FindByHashClient(t *testing.T) {
	r := mock.NewAuthRepository()

	t.Run("Existing token", func(t *testing.T) {
		token, err := r.FindByHashClient(nil, helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011"), "sdfsdf5K9QwC6mptVSJVvAuFvA4w245HsiXxfMpOtpzASJ4Rr6E")
		if assert.NoError(t, err) {
			assert.Equal(t, "775a5b37-1742-4e54-9439-0357e768b011", token.ID.String())
		}
	})

	t.Run("Non-existing token", func(t *testing.T) {
		_, err := r.FindByHashClient(nil, helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88"), "sdfsdf5K9QwCasdsd123123123HsiXxfMpOtpzASJ4Rr6E")
		if assert.Error(t, err) {
			assert.EqualError(t, err, "token not found", "error message %s", "formatted")
		}
	})
}

func TestMock_Auth_Create(t *testing.T) {
	r := mock.NewAuthRepository()

	data := models.Token{
		ID:        helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88"),
		ClientID:  helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011"),
		UserID:    helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011"),
		Token:     "tokenhashhere",
		ExpiresAt: time.Now().UTC().Add(time.Duration(100500) * time.Second).Unix(),
	}

	t.Run("Non-existing token", func(t *testing.T) {
		id, err := r.CreateToken(nil, &data)
		if assert.NoError(t, err) {
			assert.Equal(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88", id.ID.String())
		}
	})
}

func TestMock_Auth_Delete(t *testing.T) {
	r := mock.NewAuthRepository()

	t.Run("Existing token", func(t *testing.T) {
		assert.NoError(t, r.DeleteToken(nil, helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011")))
	})

	t.Run("Non-existing token", func(t *testing.T) {
		assert.Error(t, r.DeleteToken(nil, helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88")))
	})
}
