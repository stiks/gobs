package mock

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/repositories"
)

type authRepository struct {
	db      []models.Token
	clients []models.AuthClient
	users   []models.User
}

// NewAuthRepository ...
func NewAuthRepository() repositories.AuthRepository {
	id, err := uuid.Parse("775a5b37-1742-4e54-9439-0357e768b011")
	if err != nil {
		log.Fatalf("Unable to get UUID from string: %s", err.Error())
	}

	return &authRepository{
		db: []models.Token{
			{
				ID:        id,
				ClientID:  id,
				UserID:    id,
				Token:     "sdfsdf5K9QwC6mptVSJVvAuFvA4w245HsiXxfMpOtpzASJ4Rr6E",
				ExpiresAt: int(time.Now().AddDate(0, 0, 1).Unix()),
			},
			{
				ID:       uuid.New(),
				ClientID: id,
				UserID:   uuid.New(),
				Token:    "5K9QwC6mptVSJVvAuFvA4w245sdfsdfHsiXxfMpOtpzASJ4Rr6E",
			},
			{
				ID:       uuid.New(),
				ClientID: id,
				UserID:   id,
				Token:    "5K9QwC6mptVSJVvAuFsdsdvA4w245HsiXxfMpOtpzASJ4Rr6E",
			},
		},
		clients: []models.AuthClient{
			{
				ID:           id,
				ClientID:     "SecRetAuthKey",
				ClientSecret: "SecretSuper",
			},
			{
				ID:           uuid.New(),
				ClientID:     "RandomStuffHere",
				ClientSecret: "RandomKeySecret",
			},
			{
				ID:           uuid.New(),
				ClientID:     "MegaKey",
				ClientSecret: "MegaKeySecretSuper",
			},
		},
		users: []models.User{
			{
				ID:                id,
				Email:             "peter@test.com",
				PasswordHash:      []byte("$2a$10$kPrRofMm9VnE5w9ih6FwtuiuY/fIJ7/pcwvAmvL/3x3t2I144hyyq"),
				PasswordResetHash: "randomhash",
				IsActive:          true,
			},
			{
				ID:                uuid.New(),
				Email:             "oper@test.com",
				PasswordHash:      []byte("$2a$10$bzzov5K9QwC6mptVSJVvAuFvA4w245HsiXxfMpOtpzASJ4Rr6E/DG"),
				PasswordResetHash: "randomh123123ash",
				IsActive:          true,
			},
			{
				ID:           uuid.New(),
				Email:        "admin@test.com",
				PasswordHash: []byte("$2a$10$Dda31WQP2L.pnM4M8F3xZ.yM6vX31mCmb10t76v4ja9WrQ0XRvgDy"),
				IsActive:     false,
			},
			{
				ID:           uuid.New(),
				Email:        "root@test.com",
				PasswordHash: nil,
				IsActive:     false,
			},
		},
	}
}

// FindByClientID ...
func (r *authRepository) FindByClientID(ctx context.Context, id string) (*models.AuthClient, error) {
	for _, key := range r.clients {
		if key.ClientID == id {
			return &key, nil
		}
	}

	return nil, models.ErrAuthClientNotFound
}

// FindUserByUsername ...
func (r *authRepository) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	for _, key := range r.users {
		if key.Email == username {
			return &key, nil
		}
	}

	return nil, models.ErrUserNotFound
}

// FindUserByID ...
func (r *authRepository) FindUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	for _, key := range r.users {
		if key.ID == id {
			return &key, nil
		}
	}

	return nil, models.ErrUserNotFound
}

// UpdateLastLogin ...
func (r *authRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	for k, i := range r.users {
		if i.ID == id {
			r.users[k].LastLogin = time.Now()

			return nil
		}
	}

	return models.ErrUserNotFound
}

// FindByClientUser ...
func (r *authRepository) FindByClientUser(ctx context.Context, clientID uuid.UUID, userID uuid.UUID) (*models.Token, error) {
	for _, key := range r.db {
		if key.UserID == userID && key.ClientID == clientID {
			return &key, nil
		}
	}

	return nil, models.ErrTokenNotFound
}

// FindByHashClient ...
func (r *authRepository) FindByHashClient(ctx context.Context, clientID uuid.UUID, token string) (*models.Token, error) {
	for _, key := range r.db {
		if key.Token == token && key.ClientID == clientID {
			return &key, nil
		}
	}

	return nil, models.ErrTokenNotFound
}

// FindByID ...
func (r *authRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Token, error) {
	for _, key := range r.db {
		if key.ID == id {
			return &key, nil
		}
	}

	return nil, models.ErrTokenNotFound
}

// CreateToken ...
func (r *authRepository) CreateToken(ctx context.Context, data *models.Token) (*models.Token, error) {
	r.db = append(r.db, *data)

	return data, nil
}

// Delete ...
func (r *authRepository) DeleteToken(ctx context.Context, id uuid.UUID) error {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}

	var db []models.Token
	for _, k := range r.db {
		if k.ID != id {
			db = append(db, k)
		}
	}

	r.db = db

	return nil
}
