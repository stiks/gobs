package mock

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stiks/gobs/pkg/helpers"

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
	return &authRepository{
		db: []models.Token{
			{
				ID:        helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
				ClientID:  helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
				UserID:    helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
				Token:     "sdfsdf5K9QwC6mptVSJVvAuFvA4w245HsiXxfMpOtpzASJ4Rr6E",
				ExpiresAt: int(time.Now().AddDate(0, 0, 1).Unix()),
			},
			{
				ID:       uuid.New(),
				ClientID: helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
				UserID:   uuid.New(),
				Token:    "5K9QwC6mptVSJVvAuFvA4w245sdfsdfHsiXxfMpOtpzASJ4Rr6E",
			},
			{
				ID:       uuid.New(),
				ClientID: helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
				UserID:   helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
				Token:    "5K9QwC6mptVSJVvAuFsdsdvA4w245HsiXxfMpOtpzASJ4Rr6E",
			},
		},
		clients: []models.AuthClient{
			{
				ID:           helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
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
		users: _usersList,
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
