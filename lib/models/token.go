package models

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var (
	// ErrRefreshTokenEmpty ...
	ErrRefreshTokenEmpty = errors.New("refresh token is empty or missing")
	// ErrRefreshTokenNotFound ...
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	// ErrRefreshTokenExpired ...
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	// ErrTokenNotFound ...
	ErrTokenNotFound = errors.New("token not found")
)

// Token ...
type Token struct {
	ID        uuid.UUID   `json:"id"`
	ClientID  uuid.UUID   `json:"client_id"`
	Client    *AuthClient `json:"-"`
	UserID    uuid.UUID   `json:"user_id"`
	User      *User       `json:"-"`
	Token     string      `json:"token"`
	ExpiresAt int64       `json:"expires_at"`
}

// NewTokenResponse ...
func NewTokenResponse(accessToken *Token, refreshToken *Token, lifetime int, theTokenType string) (*TokenResponse, error) {
	response := &TokenResponse{
		UserID:      accessToken.UserID,
		User:        accessToken.User,
		AccessToken: accessToken.Token,
		ExpiresIn:   lifetime,
		TokenType:   theTokenType,
		Authority:   accessToken.User.Role,
	}

	if refreshToken != nil {
		response.RefreshToken = refreshToken.Token
	}

	return response, nil
}

// NewAccessToken creates new OauthAccessToken instance
func NewAccessToken(client *AuthClient, user *User, expiresIn int, jwtSecret []byte) (*Token, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := make(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["exp"] = time.Now().UTC().Add(time.Duration(expiresIn) * time.Second).Unix()
	claims["iat"] = time.Now().UTC().Unix()
	claims["auth"] = user.Role
	//	claims["iss"]     = "issuer"
	//	claims["sub"]     = "issuer"

	token.Claims = claims

	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	accessToken := &Token{
		ClientID:  client.ID,
		Token:     t,
		ExpiresAt: time.Now().UTC().Add(time.Duration(expiresIn) * time.Second).Unix(),
		UserID:    user.ID,
		User:      user,
	}

	return accessToken, nil
}

// NewRefreshToken creates new Token instance
func NewRefreshToken(client *AuthClient, user *User, expiresIn int) *Token {
	refreshToken := &Token{
		ClientID:  client.ID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().UTC().Add(time.Duration(expiresIn) * time.Second).Unix(),
		UserID:    user.ID,
		User:      user,
	}

	return refreshToken
}
