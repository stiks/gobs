package models

import (
	"errors"
	"regexp"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUserNotFound ...
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidUsernameOrPassword ...
	ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
	// ErrCannotSetEmptyUsername ...
	ErrCannotSetEmptyUsername = errors.New("cannot set empty username")
	// ErrUserPasswordNotSet ...
	ErrUserPasswordNotSet = errors.New("user password not set")
	// ErrUsernameTaken ...
	ErrUsernameTaken = errors.New("username taken")
	// ErrInvalidUUID ...
	ErrInvalidUUID = errors.New("invalid UUID")
	// ErrUserIsLocked ...
	ErrUserIsLocked = errors.New("user account is locked")
	// ErrEmailInvalidCode ...
	ErrEmailInvalidCode = errors.New("invalid email confirmation code supplied")
	// ErrEmailCodeIsEmpty ...
	ErrEmailCodeIsEmpty = errors.New("email confirmation code cannot be blank")
	// ErrEmailCodeExpired ...
	ErrEmailCodeExpired = errors.New("email confirmation code already used or expired")
)

const (
	// RoleSuperUser ...
	RoleSuperUser = "super"
	// RoleAdmin ...
	RoleAdmin = "admin"
	// RoleClient ...
	RoleClient = "client"
	// RoleManager ...
	RoleManager = "manager"
	// RoleUser ...
	RoleUser = "user"
)

// UserQueryParams ...
type UserQueryParams struct {
	Page    int    `query:"page"`
	PerPage int    `query:"perPage"`
	Role    string `query:"role"`
	Status  *int   `query:"status"`
	Query   string `query:"query"`
}

// User model
type User struct {
	ID                uuid.UUID `json:"id"         sql:"type:uuid,pk"`
	FirstName         string    `json:"firstName"  sql:"type:varchar(255)"`
	LastName          string    `json:"lastName"   sql:"type:varchar(255)"`
	Email             string    `json:"email"      sql:",unique,index"`
	Verified          bool      `json:"verified"`
	PasswordHash      []byte    `json:"-"          sql:",index"`
	PasswordResetHash string    `json:"-"          sql:"type:varchar(128),index"`
	ValidationHash    string    `json:"-"          sql:"type:varchar(128),index"`
	Role              string    `json:"role"       sql:"type:varchar(128)"`
	Status            int       `json:"status"`
	IsDeleted         bool      `json:"-"`
	OwnerID           uuid.UUID `json:"ownerId"    sql:",type:uuid"`
	Owner             *User     `json:"owner"`
	Locked            bool      `json:"locked"`
	IsActive          bool      `json:"active"`
	PasswordResetAt   time.Time `json:"-"`
	CreatedAt         time.Time `json:"createdAt"  sql:"default:now()"`
	UpdatedAt         time.Time `json:"updatedAt"  sql:"default:now()"`
	LastLogin         time.Time `json:"lastLogin"`
}

// SetPassword will set users password
func (u *User) SetPassword(userPassword string) error {
	var err error

	if len(userPassword) <= 0 {
		return ErrUserPasswordNotSet
	}

	u.PasswordHash, err = bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return nil
}

// ValidatePassword ...
func (u *User) ValidatePassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return false
	}

	return true
}

// GeneratePasswordResetHash will generate unique hash for password reset
func (u *User) GeneratePasswordResetHash() {
	u.PasswordResetHash = uuid.New().String()
}

// Validate user model
func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"))),
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.LastName, validation.Required),
		validation.Field(&u.Role, validation.Required, validation.In(
			RoleAdmin,
			RoleClient,
			RoleManager,
			RoleSuperUser,
			RoleUser,
		)),
		validation.Field(&u.Status, validation.In(
			StatusInit,
			StatusDraft,
			StatusActive,
		)),
	)
}

// PasswordResetRequest ...
type PasswordResetRequest struct {
	Email string `json:"email" form:"email" query:"email" validate:"email"`
}

// Validate admin auth model
func (u *PasswordResetRequest) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"))),
	)
}

// CreateUser ...
type CreateUser struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Password  string    `json:"password"`
	Status    int       `json:"status"`
	Active    bool      `json:"active"`
}

// Validate ...
func (u *CreateUser) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"))),
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.LastName, validation.Required),
		validation.Field(&u.Password, validation.Length(8, 64)),
		validation.Field(&u.Role, validation.Required, validation.In(
			RoleAdmin,
			RoleClient,
			RoleManager,
			RoleSuperUser,
			RoleUser,
		)),
		validation.Field(&u.Status, validation.In(
			StatusInit,
			StatusDraft,
			StatusActive,
		)),
	)
}

// GeneratePassword will generate unique hash for password reset
func (u *CreateUser) GeneratePassword() {
	u.Password = uuid.New().String()
}

// ToUser ...
func (u *CreateUser) ToUser(id *uuid.UUID) *User {
	user := &User{
		ID:        uuid.New(),
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Status:    u.Status,
		Role:      u.Role,
		IsActive:  u.Active,
	}

	if id != nil {
		user.OwnerID = *id
	}

	// If password is empty, auto generate something
	if len(u.Password) <= 0 {
		u.GeneratePassword()
	}

	// Set user's password
	user.SetPassword(u.Password)

	return user
}

// UpdateUser ...
type UpdateUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
	Status    int    `json:"status"`
	Active    bool   `json:"active"`
}

// Validate ...
func (u *UpdateUser) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.LastName, validation.Required),
		validation.Field(&u.Role, validation.Required, validation.In(
			RoleAdmin,
			RoleClient,
			RoleManager,
			RoleSuperUser,
			RoleUser,
		)),
		validation.Field(&u.Status, validation.In(
			StatusInit,
			StatusDraft,
			StatusActive,
		)),
	)
}

// FromUpdate ...
func (u *User) FromUpdate(data *UpdateUser) {
	u.FirstName = data.FirstName
	u.LastName = data.LastName
	u.Role = data.Role
	u.IsActive = data.Active
	u.Status = data.Status
}

// EmailConfirmationCode ...
type EmailConfirmationCode struct {
	Code     string `json:"code"     form:"code"     query:"code"`
	Password string `json:"password" form:"password" query:"password"`
}

// Validate ...
func (u *EmailConfirmationCode) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Code, validation.Required),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 64)),
	)
}

// ConfirmEmail ...
type ConfirmEmail struct {
	UserID uuid.UUID `json:"id"     form:"id"     query:"id"`
	Code   string    `json:"code" form:"code" query:"code"`
}

// Validate ...
func (u *ConfirmEmail) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.UserID, validation.Required, is.UUID),
		validation.Field(&u.Code, validation.Required, validation.Length(8, 64)),
	)
}
