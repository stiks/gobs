package models

import (
	"errors"

	"github.com/google/uuid"
)

var (
	// ErrAuthClientNotFound ...
	ErrAuthClientNotFound = errors.New("auth client could not be found")
	// ErrAuthClientAlreadyExist ...
	ErrAuthClientAlreadyExist = errors.New("auth client already exist")
)

// AuthClient ...
type AuthClient struct {
	ID           uuid.UUID `json:"id"`
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
}

// ValidateSecret ...
func (u *AuthClient) ValidateSecret(secret string) bool {
	if u.ClientSecret == secret {
		return true
	}

	return false
}
