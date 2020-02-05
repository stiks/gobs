package models

import (
	"errors"
	"regexp"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUserNotFound ...
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidUsernameOrPassword ...
	ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
	// ErrInvalidUserPassword ...
	ErrInvalidUserPassword = errors.New("invalid user password")
	// ErrCannotSetEmptyUsername ...
	ErrCannotSetEmptyUsername = errors.New("cannot set empty username")
	// ErrUserPasswordNotSet ...
	ErrUserPasswordNotSet = errors.New("user password not set")
	// ErrUsernameTaken ...
	ErrUsernameTaken = errors.New("username taken")
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
	// RoleTranslator ...
	RoleTranslator = "translator"
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
	PasswordResetHash string    `json:"-"          sql:"type:varchar(255),index"`
	Role              string    `json:"role"       sql:"type:varchar(128)"`
	Status            int       `json:"status"`
	IsDeleted         bool      `json:"-"`
	OwnerID           uuid.UUID `json:"ownerId"    sql:",type:uuid"`
	Owner             *User     `json:"owner"`
	Locked            bool      `json:"locked"`
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
		validation.Field(&u.Role, validation.Required),
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
}

// Validate ...
func (u *CreateUser) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"))),
		validation.Field(&u.FirstName, validation.Required),
		validation.Field(&u.LastName, validation.Required),
		validation.Field(&u.Role, validation.Required),
		validation.Field(&u.Password),
		validation.Field(&u.Status, validation.In(0, 1, 2)),
	)
}

// GeneratePassword will generate unique hash for password reset
func (u *CreateUser) GeneratePassword() {
	u.Password = uuid.New().String()
}
