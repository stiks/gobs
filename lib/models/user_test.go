package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stiks/gobs/lib/models"
)

func TestModel_User_SetPassword(t *testing.T) {
	user := new(models.User)

	t.Run("Good password", func(t *testing.T) {
		assert.NoError(t, user.SetPassword("Go0dP4ssword"))
	})

	t.Run("Empty password", func(t *testing.T) {
		err := user.SetPassword("")
		if assert.Error(t, err) {
			assert.EqualError(t, err, "user password not set", "error message %s", "formatted")
		}
	})
}

func TestModel_User_ValidatePassword(t *testing.T) {
	user := models.User{
		PasswordHash: []byte("$2a$10$kPrRofMm9VnE5w9ih6FwtuiuY/fIJ7/pcwvAmvL/3x3t2I144hyyq"),
	}

	t.Run("Good password", func(t *testing.T) {
		assert.Equal(t, true, user.ValidatePassword("testpass"))
	})

	t.Run("Wrong password", func(t *testing.T) {
		assert.Equal(t, false, user.ValidatePassword("wrong-password"))
	})
}

func TestModel_User_GeneratePasswordResetHash(t *testing.T) {
	user := new(models.User)
	user.GeneratePasswordResetHash()

	assert.NotEmpty(t, user.PasswordResetHash)
}

func TestModel_User_Validate(t *testing.T) {
	t.Run("Good model", func(t *testing.T) {
		user := models.User{
			FirstName: "John",
			LastName:  "Snow",
			Email:     "john@snow.com",
			Role:      models.RoleAdmin,
			Status:    models.StatusActive,
		}

		assert.NoError(t, user.Validate())
	})

	t.Run("Bad email", func(t *testing.T) {
		user := models.User{
			FirstName: "John",
			LastName:  "Snow",
			Email:     "john",
			Role:      models.RoleAdmin,
			Status:    models.StatusActive,
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "email: must be in a valid format.", "error message %s", "formatted")
		}
	})

	t.Run("Empty email", func(t *testing.T) {
		user := models.User{
			FirstName: "John",
			LastName:  "Snow",
			Role:      models.RoleAdmin,
			Status:    models.StatusActive,
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "email: cannot be blank.", "error message %s", "formatted")
		}
	})

	t.Run("Empty first name", func(t *testing.T) {
		user := models.User{
			LastName: "Snow",
			Email:    "john@snow.com",
			Role:     models.RoleAdmin,
			Status:   models.StatusActive,
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "firstName: cannot be blank.", "error message %s", "formatted")
		}
	})

	t.Run("Empty last name", func(t *testing.T) {
		user := models.User{
			FirstName: "John",
			Email:     "john@snow.com",
			Role:      models.RoleAdmin,
			Status:    models.StatusActive,
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "lastName: cannot be blank.", "error message %s", "formatted")
		}
	})

	t.Run("Empty role", func(t *testing.T) {
		user := models.User{
			FirstName: "John",
			LastName:  "Snow",
			Email:     "john@snow.com",
			Status:    models.StatusActive,
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "role: cannot be blank.", "error message %s", "formatted")
		}
	})

	t.Run("Wrong role", func(t *testing.T) {
		user := models.User{
			FirstName: "John",
			LastName:  "Snow",
			Email:     "john@snow.com",
			Role:      "zzz",
			Status:    models.StatusActive,
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "role: must be a valid value.", "error message %s", "formatted")
		}
	})

	t.Run("Wrong status", func(t *testing.T) {
		user := models.User{
			FirstName: "John",
			LastName:  "Snow",
			Email:     "john@snow.com",
			Role:      models.RoleAdmin,
			Status:    99,
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "status: must be a valid value.", "error message %s", "formatted")
		}
	})
}

func TestModel_PasswordResetRequest_Validate(t *testing.T) {
	t.Run("Good email", func(t *testing.T) {
		data := &models.PasswordResetRequest{Email: "john@snow.com"}

		assert.NoError(t, data.Validate())
	})

	t.Run("Wrong email", func(t *testing.T) {
		data := &models.PasswordResetRequest{Email: "snow.com"}

		err := data.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "email: must be in a valid format.", "error message %s", "formatted")
		}
	})
}

func TestModel_CreateUser_Validate(t *testing.T) {
	t.Run("Good user", func(t *testing.T) {
		user := models.CreateUser{
			FirstName: "John",
			LastName:  "Snow",
			Email:     "john@snow.com",
			Role:      models.RoleAdmin,
			Status:    models.StatusActive,
			Password:  "G0odP4ssword",
		}

		assert.NoError(t, user.Validate())
	})

	t.Run("Bad email", func(t *testing.T) {
		user := models.CreateUser{
			FirstName: "John",
			LastName:  "Snow",
			Email:     "john",
			Role:      models.RoleAdmin,
			Status:    models.StatusActive,
			Password:  "G0odP4ssword",
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "email: must be in a valid format.", "error message %s", "formatted")
		}
	})

	t.Run("Empty email", func(t *testing.T) {
		user := models.CreateUser{
			FirstName: "John",
			LastName:  "Snow",
			Role:      models.RoleAdmin,
			Status:    models.StatusActive,
			Password:  "G0odP4ssword",
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "email: cannot be blank.", "error message %s", "formatted")
		}
	})

	t.Run("Empty first name", func(t *testing.T) {
		user := models.CreateUser{
			LastName: "Snow",
			Email:    "john@snow.com",
			Role:     models.RoleAdmin,
			Status:   models.StatusActive,
			Password: "G0odP4ssword",
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "firstName: cannot be blank.", "error message %s", "formatted")
		}
	})

	t.Run("Empty last name", func(t *testing.T) {
		user := models.CreateUser{
			FirstName: "John",
			Email:     "john@snow.com",
			Role:      models.RoleAdmin,
			Status:    models.StatusActive,
			Password:  "G0odP4ssword",
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "lastName: cannot be blank.", "error message %s", "formatted")
		}
	})

	t.Run("Empty role", func(t *testing.T) {
		user := models.CreateUser{
			FirstName: "John",
			LastName:  "Snow",
			Email:     "john@snow.com",
			Status:    models.StatusActive,
			Password:  "G0odP4ssword",
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "role: cannot be blank.", "error message %s", "formatted")
		}
	})

	t.Run("Wrong role", func(t *testing.T) {
		user := models.CreateUser{
			FirstName: "John",
			LastName:  "Snow",
			Email:     "john@snow.com",
			Role:      "zzz",
			Status:    models.StatusActive,
			Password:  "G0odP4ssword",
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "role: must be a valid value.", "error message %s", "formatted")
		}
	})

	t.Run("Wrong status", func(t *testing.T) {
		user := models.CreateUser{
			FirstName: "John",
			LastName:  "Snow",
			Email:     "john@snow.com",
			Role:      models.RoleAdmin,
			Status:    99,
			Password:  "G0odP4ssword",
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "status: must be a valid value.", "error message %s", "formatted")
		}
	})

	t.Run("Empty password", func(t *testing.T) {
		user := models.CreateUser{
			FirstName: "John",
			LastName:  "Snow",
			Email:     "john@snow.com",
			Role:      models.RoleAdmin,
			Status:    models.StatusActive,
		}

		err := user.Validate()
		assert.NoError(t, err)
	})

	t.Run("Short password", func(t *testing.T) {
		user := models.CreateUser{
			FirstName: "John",
			LastName:  "Snow",
			Email:     "john@snow.com",
			Role:      models.RoleAdmin,
			Status:    models.StatusActive,
			Password:  "G00d",
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "password: the length must be between 8 and 64.", "error message %s", "formatted")
		}
	})

	t.Run("Long password", func(t *testing.T) {
		user := models.CreateUser{
			FirstName: "John",
			LastName:  "Snow",
			Email:     "john@snow.com",
			Role:      models.RoleAdmin,
			Status:    models.StatusActive,
			Password:  "G00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00d",
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "password: the length must be between 8 and 64.", "error message %s", "formatted")
		}
	})
}

func TestModel_CreateUser_GeneratePassword(t *testing.T) {
	user := new(models.CreateUser)
	user.GeneratePassword()
	assert.NotEmpty(t, user.Password)
}

func TestModel_UpdateUser_Validate(t *testing.T) {
	t.Run("Good user", func(t *testing.T) {
		user := models.UpdateUser{
			FirstName: "John",
			LastName:  "Snow",
			Role:      models.RoleAdmin,
			Status:    models.StatusActive,
		}

		assert.NoError(t, user.Validate())
	})

	t.Run("Empty first name", func(t *testing.T) {
		user := models.UpdateUser{
			LastName: "Snow",
			Role:     models.RoleAdmin,
			Status:   models.StatusActive,
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "firstName: cannot be blank.", "error message %s", "formatted")
		}
	})

	t.Run("Empty last name", func(t *testing.T) {
		user := models.UpdateUser{
			FirstName: "John",
			Role:      models.RoleAdmin,
			Status:    models.StatusActive,
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "lastName: cannot be blank.", "error message %s", "formatted")
		}
	})

	t.Run("Empty role", func(t *testing.T) {
		user := models.UpdateUser{
			FirstName: "John",
			LastName:  "Snow",
			Status:    models.StatusActive,
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "role: cannot be blank.", "error message %s", "formatted")
		}
	})

	t.Run("Wrong role", func(t *testing.T) {
		user := models.UpdateUser{
			FirstName: "John",
			LastName:  "Snow",
			Role:      "zzz",
			Status:    models.StatusActive,
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "role: must be a valid value.", "error message %s", "formatted")
		}
	})

	t.Run("Wrong status", func(t *testing.T) {
		user := models.UpdateUser{
			FirstName: "John",
			LastName:  "Snow",
			Role:      models.RoleAdmin,
			Status:    99,
		}

		err := user.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "status: must be a valid value.", "error message %s", "formatted")
		}
	})
}

func TestModel_UpdateUser_FromUpdate(t *testing.T) {
	user := models.User{
		FirstName: "John",
		LastName:  "Snow",
		Role:      models.RoleAdmin,
		Status:    models.StatusActive,
		IsActive:  true,
	}

	user.FromUpdate(&models.UpdateUser{
		FirstName: "Test",
		LastName:  "Tester",
		Role:      models.RoleUser,
		Status:    models.StatusDisabled,
		Active:    false,
	})

	assert.Equal(t, "Test", user.FirstName)
	assert.Equal(t, "Tester", user.LastName)
	assert.Equal(t, models.RoleUser, user.Role)
	assert.Equal(t, models.StatusDisabled, user.Status)
	assert.Equal(t, false, user.IsActive)
}

func TestModel_EmailConfirmationCode_Validate(t *testing.T) {
	t.Run("Good values", func(t *testing.T) {
		data := models.EmailConfirmationCode{
			Code:     "zzzz",
			Password: "zzZzzzZzzzZzzzZz",
		}

		assert.NoError(t, data.Validate())
	})

	t.Run("Empty code", func(t *testing.T) {
		data := models.EmailConfirmationCode{
			Password: "zzZzzzZzzzZzzzZz",
		}

		err := data.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "code: cannot be blank.", "error message %s", "formatted")
		}

	})

	t.Run("Empty password", func(t *testing.T) {
		data := models.EmailConfirmationCode{
			Code: "zzzz",
		}

		err := data.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "password: cannot be blank.", "error message %s", "formatted")
		}
	})

	t.Run("Short password", func(t *testing.T) {
		data := models.EmailConfirmationCode{
			Code:     "zzzz",
			Password: "G00d",
		}

		err := data.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "password: the length must be between 8 and 64.", "error message %s", "formatted")
		}
	})

	t.Run("Long password", func(t *testing.T) {
		data := models.EmailConfirmationCode{
			Code:     "zzzz",
			Password: "G00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00dG00d",
		}

		err := data.Validate()
		if assert.Error(t, err) {
			assert.EqualError(t, err, "password: the length must be between 8 and 64.", "error message %s", "formatted")
		}
	})
}
