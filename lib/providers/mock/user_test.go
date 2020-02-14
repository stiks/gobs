package mock_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/repositories"
	"github.com/stiks/gobs/pkg/helpers"
)

func TestMock_User_NewUserRepository(t *testing.T) {
	r := mock.NewUserRepository()

	assert.Implements(t, (*repositories.UserRepository)(nil), r)
}

func TestMock_User_FindByUsername(t *testing.T) {
	r := mock.NewUserRepository()

	t.Run("Existing user", func(t *testing.T) {
		user, err := r.FindByUsername(nil, "peter@test.com")
		if assert.NoError(t, err) {
			assert.Equal(t, "775a5b37-1742-4e54-9439-0357e768b011", user.ID.String())
		}
	})

	t.Run("Non-existing user", func(t *testing.T) {
		_, err := r.FindByUsername(nil, "anything@test.com")
		if assert.Error(t, err) {
			assert.EqualError(t, err, "user not found", "error message %s", "formatted")
		}
	})
}

func TestMock_User_FindByResetHash(t *testing.T) {
	r := mock.NewUserRepository()

	t.Run("Existing user", func(t *testing.T) {
		user, err := r.FindByResetHash(nil, "random")
		if assert.NoError(t, err) {
			assert.Equal(t, "775a5b37-1742-4e54-9439-0357e768b011", user.ID.String())
		}
	})

	t.Run("Non-existing user", func(t *testing.T) {
		_, err := r.FindByResetHash(nil, "anything random")
		if assert.Error(t, err) {
			assert.EqualError(t, err, "user not found", "error message %s", "formatted")
		}
	})
}

func TestMock_User_CountAll(t *testing.T) {
	r := mock.NewUserRepository()

	count, err := r.CountAll(nil, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, 4, count)
	}
}

func TestMock_User_FindAll(t *testing.T) {
	r := mock.NewUserRepository()

	users, err := r.FindAll(nil, nil)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, users)
	}
}

func TestMock_User_FindByID(t *testing.T) {
	r := mock.NewUserRepository()

	t.Run("Existing user", func(t *testing.T) {
		user, err := r.FindByID(nil, helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011"))
		if assert.NoError(t, err) {
			assert.Equal(t, "peter@test.com", user.Email)
		}
	})

	t.Run("Non-existing user", func(t *testing.T) {
		_, err := r.FindByID(nil, helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88"))
		if assert.Error(t, err) {
			assert.EqualError(t, err, "user not found", "error message %s", "formatted")
		}
	})
}

func TestMock_User_Update(t *testing.T) {
	r := mock.NewUserRepository()

	t.Run("Existing user", func(t *testing.T) {
		data := models.User{
			ID:                helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011"),
			Email:             "update@test.com",
			PasswordHash:      []byte("hash"),
			PasswordResetHash: "resethash",
			Status:            10,
		}

		user, err := r.Update(nil, &data)
		if assert.NoError(t, err) {
			assert.Equal(t, "update@test.com", user.Email)
		}
	})

	t.Run("Non-existing user", func(t *testing.T) {
		data := models.User{
			ID:                helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88"),
			Email:             "nonexisting@test.com",
			PasswordHash:      []byte("hash"),
			PasswordResetHash: "resethash",
			Status:            10,
		}

		_, err := r.Update(nil, &data)
		if assert.Error(t, err) {
			assert.EqualError(t, err, "user not found", "error message %s", "formatted")
		}
	})
}

func TestMock_User_Create(t *testing.T) {
	r := mock.NewUserRepository()

	t.Run("Non-existing user", func(t *testing.T) {
		data := models.User{
			ID:                helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88"),
			Email:             "test@test.com",
			PasswordHash:      []byte("hash"),
			PasswordResetHash: "resethash",
			Status:            10,
		}

		id, err := r.Create(nil, &data)
		if assert.NoError(t, err) {
			assert.Equal(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88", id.ID.String())
		}
	})

	t.Run("Existing user", func(t *testing.T) {
		data := models.User{
			ID:                helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011"),
			Email:             "test@test.com",
			PasswordHash:      []byte("hash"),
			PasswordResetHash: "resethash",
			Status:            10,
		}

		_, err := r.Create(nil, &data)
		assert.Error(t, err)
	})
}

func TestMock_User_Delete(t *testing.T) {
	r := mock.NewUserRepository()

	t.Run("Existing user", func(t *testing.T) {
		assert.NoError(t, r.Delete(nil, helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011")))
	})

	t.Run("Non-existing user", func(t *testing.T) {
		assert.Error(t, r.Delete(nil, helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88")))
	})
}
