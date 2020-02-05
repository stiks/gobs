package models

import (
	"errors"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// AuthRequest ...
type AuthRequest struct {
	ClientID     string `json:"client_id"     form:"client_id"     query:"client_id"`
	ClientSecret string `json:"client_secret" form:"client_secret" query:"client_secret"`
	GrantType    string `json:"grant_type"    form:"grant_type"    query:"grant_type" validate:"required"`
	Username     string `json:"username"      form:"username"      query:"username"`
	Password     string `json:"password"      form:"password"      query:"password"`
	RefreshToken string `json:"refresh_token" form:"refresh_token" query:"refresh_token"`
}

var (
	// ErrInvalidGrantType ...
	ErrInvalidGrantType = errors.New("invalid grant type")
	// ErrInvalidClientIDOrSecret ...
	ErrInvalidClientIDOrSecret = errors.New("invalid client ID or secret")
	// ErrEmptyClientIDOrSecret ...
	ErrEmptyClientIDOrSecret = errors.New("client ID or secret cannot be empty")
	// ErrInvalidScope ...
	ErrInvalidScope = errors.New("invalid scope")
)

// Validate users model
func (u *AuthRequest) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.GrantType, validation.Required),
	)
}

// Validate users model
func (u *AuthRequest) ValidateLogin() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Username, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required),
	)
}
