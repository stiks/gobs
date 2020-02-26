package mock

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/repositories"
	"github.com/stiks/gobs/pkg/helpers"
)

var (
	_usersList = []models.User{
		{
			ID:                helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
			Email:             "peter@test.com",
			PasswordHash:      []byte("$2a$10$kPrRofMm9VnE5w9ih6FwtuiuY/fIJ7/pcwvAmvL/3x3t2I144hyyq"),
			PasswordResetHash: "random",
			FirstName:         "Apple",
			LastName:          "Appleton",
			IsActive:          true,
			Verified:          true,
			Role:              models.RoleSuperUser,
			Status:            models.StatusActive,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			ID:                uuid.New(),
			Email:             "oper@test.com",
			PasswordHash:      []byte("$2a$10$bzzov5K9QwC6mptVSJVvAuFvA4w245HsiXxfMpOtpzASJ4Rr6E/DG"),
			PasswordResetHash: "5zQVfk8aQlZgQiW0vd2PA8kyj4",
			PasswordResetAt:   time.Now().AddDate(0, 0, -2),
			OwnerID:           helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
			Status:            models.StatusActive,
			Role:              models.RoleUser,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			ID:                helpers.UUIDFromString(nil, "3ab1ba2a-6031-4e34-aae3-dcd43a987775"),
			Email:             "admin@test.com",
			FirstName:         "Admin",
			LastName:          "Example",
			PasswordHash:      []byte("$2a$10$Dda31WQP2L.pnM4M8F3xZ.yM6vX31mCmb10t76v4ja9WrQ0XRvgDy"),
			OwnerID:           helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
			Status:            models.StatusInit,
			Role:              models.RoleUser,
			PasswordResetHash: "ZXqEMubf5DinaTHuOyJIm1z3Dq",
			PasswordResetAt:   time.Now(),
			IsActive:          false,
			ValidationHash:    "SomeHash123",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
		{
			ID:                uuid.New(),
			Email:             "root@test.com",
			PasswordHash:      nil,
			Role:              models.RoleUser,
			PasswordResetHash: "2e4EHSsVkledZxWwU7j3BnNBYo",
			Status:            models.StatusInit,
			Locked:            true,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		},
	}
)

// NewUserRepository ...
func NewUserRepository() repositories.UserRepository {
	return &userRepository{
		db: _usersList,
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
		log.Printf("PWD: %s HASH: %s", key.PasswordResetHash, hash)

		if key.PasswordResetHash == hash {
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
