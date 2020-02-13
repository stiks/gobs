package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/providers/mock"
	"github.com/stiks/gobs/lib/services"
	"github.com/stiks/gobs/pkg/helpers"
)

func _userSrv() services.UserService {
	return services.NewUserService(mock.NewUserRepository(), services.NewQueueService(mock.NewQueueRepository()), services.NewCacheService(mock.NewCacheRepository()))
}

func TestService_User_NewUserService(t *testing.T) {
	assert.Implements(t, (*services.UserService)(nil), _userSrv())
}

func TestService_User_GetByUsername(t *testing.T) {
	srv := _userSrv()

	t.Run("Get existing user by username", func(t *testing.T) {
		user, err := srv.GetByUsername(nil, "peter@test.com")
		if assert.NoError(t, err) {
			assert.Equal(t, "775a5b37-1742-4e54-9439-0357e768b011", user.ID.String())
		}
	})

	t.Run("Get non-existing user by username", func(t *testing.T) {
		_, err := srv.GetByUsername(nil, "wrong@test.com")
		if assert.Error(t, err) {
			assert.EqualError(t, err, "user not found")
		}
	})
}

func TestService_User_GetAll(t *testing.T) {
	srv := _userSrv()

	_, err := srv.GetAll(nil, &models.UserQueryParams{})
	assert.NoError(t, err)
}

func TestService_User_CountAll(t *testing.T) {
	srv := _userSrv()

	total, err := srv.CountAll(nil, &models.UserQueryParams{})
	if assert.NoError(t, err) {
		assert.Equal(t, 4, total)
	}
}

func TestService_User_GetByID(t *testing.T) {
	srv := _userSrv()

	t.Run("Get user by ID", func(t *testing.T) {
		user, err := srv.GetByID(nil, helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011"))
		if assert.NoError(t, err) {
			assert.Equal(t, "peter@test.com", user.Email)
		}
	})

	t.Run("Get user by ID, non-existing", func(t *testing.T) {
		_, err := srv.GetByID(nil, helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88"))
		if assert.Error(t, err) {
			assert.EqualError(t, err, "user not found")
		}
	})
}

func TestService_User_Create(t *testing.T) {
	srv := _userSrv()

	newUser := models.User{
		ID:     helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88"),
		Email:  "new@friend.com",
		Status: 10,
	}

	t.Run("Create new user", func(t *testing.T) {
		user, err := srv.Create(nil, "testpass", &newUser)
		if assert.NoError(t, err) {
			assert.Equal(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88", user.ID.String())
		}
	})

	t.Run("Create new user with empty password", func(t *testing.T) {
		_, err := srv.Create(nil, "", &newUser)
		assert.Error(t, err)
	})
}

func TestService_User_Update(t *testing.T) {
	srv := _userSrv()

	t.Run("Update user", func(t *testing.T) {
		ret, err := srv.Update(nil, &models.User{
			ID:           helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011"),
			Email:        "peter@test.com",
			PasswordHash: []byte("$2a$10$kPrRofMm9VnE5w9ih6FwtuiuY/fIJ7/pcwvAmvL/3x3t2I144hyyq"),
			Status:       20,
		})
		if assert.NoError(t, err) {
			assert.Equal(t, 20, ret.Status)
		}
	})

	t.Run("Update non-existing user", func(t *testing.T) {
		_, err := srv.Update(nil, &models.User{
			ID:           helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88"),
			Email:        "peter@test.com",
			PasswordHash: []byte("$2a$10$kPrRofMm9VnE5w9ih6FwtuiuY/fIJ7/pcwvAmvL/3x3t2I144hyyq"),
			Status:       20,
		})
		if assert.Error(t, err) {
			assert.EqualError(t, err, "user not found")
		}
	})
}

func TestService_User_UpdatePassword(t *testing.T) {
	srv := _userSrv()

	t.Run("Update user's password", func(t *testing.T) {
		user, err := srv.UpdatePassword(nil, helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011"), "testpwd")
		if assert.NoError(t, err) {
			assert.NotEqual(t, []byte("$2a$10$kPrRofMm9VnE5w9ih6FwtuiuY/fIJ7/pcwvAmvL/3x3t2I144hyyq"), user.PasswordHash)
		}
	})

	t.Run("Update password non-existing password", func(t *testing.T) {
		_, err := srv.UpdatePassword(nil, helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88"), "newpwd")
		if assert.Error(t, err) {
			assert.EqualError(t, err, "user not found")
		}
	})
}

func TestService_User_UpdateUsername(t *testing.T) {
	srv := _userSrv()

	t.Run("Update user username", func(t *testing.T) {
		user, err := srv.UpdateUsername(nil, helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011"), "uber@yn.ee")
		if assert.NoError(t, err) {
			assert.Equal(t, "uber@yn.ee", user.Email)
		}
	})

	t.Run("Update user username, non-existing", func(t *testing.T) {
		_, err := srv.UpdateUsername(nil, helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88"), "new@username.com")
		if assert.Error(t, err) {
			assert.EqualError(t, err, "user not found")
		}
	})
}

func TestService_User_GetByPwdResetHash(t *testing.T) {
	srv := _userSrv()

	t.Run("Get user by reset hash", func(t *testing.T) {
		user, err := srv.GetByResetHash(nil, "randomhash")
		if assert.NoError(t, err) {
			assert.Equal(t, "775a5b37-1742-4e54-9439-0357e768b011", user.ID.String())
		}
	})

	t.Run("Get user by reset hash, non-existing", func(t *testing.T) {
		_, err := srv.GetByResetHash(nil, "nonexistingResetHash")
		if assert.Error(t, err) {
			assert.EqualError(t, err, "user not found")
		}
	})
}

func TestService_User_ResetPassword(t *testing.T) {
	srv := _userSrv()

	t.Run("Reset password by username", func(t *testing.T) {
		user, err := srv.ResetPassword(nil, "peter@test.com")
		if assert.NoError(t, err) {
			assert.NotEqual(t, []byte(""), user.PasswordResetHash)
		}
	})

	t.Run("Non-existing username", func(t *testing.T) {
		_, err := srv.ResetPassword(nil, "new@username.com")
		if assert.Error(t, err) {
			assert.EqualError(t, err, "user not found")
		}
	})
}

func TestService_User_UpdateLogin(t *testing.T) {
	srv := _userSrv()

	user, err := srv.UpdateLogin(nil, &models.User{ID: helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011")})
	if assert.NoError(t, err) {
		assert.InDelta(t, time.Now().Unix(), user.LastLogin.Unix(), 5)
	}
}

func TestService_User_Delete(t *testing.T) {
	srv := _userSrv()

	t.Run("Delete existing user", func(t *testing.T) {
		assert.NoError(t, srv.Delete(nil, helpers.UUIDFromString(t, "775a5b37-1742-4e54-9439-0357e768b011")))
	})

	t.Run("Delete non-existing user", func(t *testing.T) {
		err := srv.Delete(nil, helpers.UUIDFromString(t, "5fcc94e5-c6aa-4320-8469-f5021af54b88"))
		if assert.Error(t, err) {
			assert.EqualError(t, err, "user not found")
		}
	})
}
