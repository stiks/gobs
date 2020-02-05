package mock

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/repositories"
)

// NewUserRepository ...
func NewUserRepository() repositories.UserRepository {
	id, err := uuid.Parse("775a5b37-1742-4e54-9439-0357e768b011")
	if err != nil {
		log.Fatalf("Unable to get UUID from string: %s", err.Error())
	}

	return &userRepository{
		db: []models.User{
			{
				ID:                id,
				Email:             "peter@test.com",
				PasswordHash:      []byte("$2a$10$kPrRofMm9VnE5w9ih6FwtuiuY/fIJ7/pcwvAmvL/3x3t2I144hyyq"),
				PasswordResetHash: "randomhash",
				Status:            models.StatusActive,
			},
			{
				ID:                uuid.New(),
				Email:             "oper@test.com",
				PasswordHash:      []byte("$2a$10$bzzov5K9QwC6mptVSJVvAuFvA4w245HsiXxfMpOtpzASJ4Rr6E/DG"),
				PasswordResetHash: "randomh123123ash",
				Status:            models.StatusActive,
			},
			{
				ID:           uuid.New(),
				Email:        "admin@test.com",
				PasswordHash: []byte("$2a$10$Dda31WQP2L.pnM4M8F3xZ.yM6vX31mCmb10t76v4ja9WrQ0XRvgDy"),
				Status:       models.StatusInit,
			},
			{
				ID:           uuid.New(),
				Email:        "root@test.com",
				PasswordHash: nil,
				Status:       models.StatusInit,
			},
		},
	}
}

type userRepository struct {
	db []models.User
}

// FindByUsername ...
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	for _, key := range r.db {
		if key.Email == username {
			return &key, nil
		}
	}

	return nil, models.ErrUserNotFound
}

// FindByResetHash ...
func (r *userRepository) FindByResetHash(ctx context.Context, hash string) (*models.User, error) {
	for _, key := range r.db {
		if string(key.PasswordResetHash) == hash {
			return &key, nil
		}
	}

	return nil, models.ErrUserNotFound
}

// FindAll ...
func (r *userRepository) FindAll(ctx context.Context, params *models.UserQueryParams) ([]models.User, error) {
	return r.db, nil
}

// CountAll ...
func (r *userRepository) CountAll(ctx context.Context, params *models.UserQueryParams) (int, error) {
	return len(r.db), nil
}

// FindByID ...
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	for _, key := range r.db {
		if key.ID == id {
			return &key, nil
		}
	}

	return nil, models.ErrUserNotFound
}

// Create ...
func (r *userRepository) Create(ctx context.Context, data *models.User) (*models.User, error) {
	for _, key := range r.db {
		if key.Email == data.Email {
			return nil, models.ErrUsernameTaken
		}
	}

	r.db = append(r.db, *data)

	return data, nil
}

// Update ...
func (r *userRepository) Update(ctx context.Context, data *models.User) (*models.User, error) {
	for k, i := range r.db {
		if i.ID == data.ID {
			r.db[k] = *data

			return data, nil
		}
	}

	return nil, models.ErrUserNotFound
}

// Delete ...
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}

	var db []models.User
	for _, k := range r.db {
		if k.ID != id {
			db = append(db, k)
		}
	}

	r.db = db

	return nil
}
