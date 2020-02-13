package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/repositories"
	"github.com/stiks/gobs/pkg/xlog"
)

type userService struct {
	repo  repositories.UserRepository
	queue QueueService
	cache CacheService
}

// UserService ...
type UserService interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByResetHash(ctx context.Context, hash string) (*models.User, error)
	CountAll(ctx context.Context, params *models.UserQueryParams) (int, error)
	GetAll(ctx context.Context, params *models.UserQueryParams) ([]models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Create(ctx context.Context, password string, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, user *models.User) (*models.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, newPassword string) (*models.User, error)
	UpdateUsername(ctx context.Context, id uuid.UUID, newUsername string) (*models.User, error)
	UpdateLogin(ctx context.Context, user *models.User) (*models.User, error)
	ResetPassword(ctx context.Context, username string) (*models.User, error)
}

// NewUserService ...
func NewUserService(repo repositories.UserRepository, queue QueueService, cacheSrv CacheService) UserService {
	return &userService{
		cache: cacheSrv,
		repo:  repo,
		queue: queue,
	}
}

// GetByUsername ...
func (s *userService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	cached := new(models.User)

	key := fmt.Sprintf("user_%s", username)

	// we have cached result
	err := s.cache.GetByKey(ctx, key, cached)
	if err == nil {
		xlog.Infof(ctx, "Returned cached result")

		return cached, nil
	}

	xlog.Infof(ctx, "Missed cache, getting from user by username service, err: %s", err.Error())

	cached, err = s.repo.FindByUsername(ctx, username)
	if err != nil {
		xlog.Errorf(ctx, "User find error: %s", err.Error())

		return nil, err
	}

	xlog.Infof(ctx, "Creating cache for %s", key)

	// adding results to the cache
	if err := s.cache.Create(ctx, key, cached); err != nil {
		xlog.Errorf(ctx, "Unable to create tag cache, err: %s", err.Error())
	}

	return cached, nil
}

// GetAll ...
func (s *userService) GetAll(ctx context.Context, params *models.UserQueryParams) ([]models.User, error) {
	return s.repo.FindAll(ctx, params)
}

// CountAll ...
func (s *userService) CountAll(ctx context.Context, params *models.UserQueryParams) (int, error) {
	return s.repo.CountAll(ctx, params)
}

// GetByID ...
func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	cached := new(models.User)

	if id.String() == "00000000-0000-0000-0000-000000000000" {
		xlog.Infof(ctx, "Empty user ID, existing ...")

		return nil, models.ErrUserNotFound
	}

	key := fmt.Sprintf("user_%s", id.String())

	// we have cached result
	err := s.cache.GetByKey(ctx, key, cached)
	if err == nil {
		xlog.Infof(ctx, "Returned cached result")

		return cached, nil
	}

	xlog.Infof(ctx, "Missed cache, getting from user by ID service, err: %s", err.Error())

	cached, err = s.repo.FindByID(ctx, id)
	if err != nil {
		xlog.Errorf(ctx, "User find error: %s", err.Error())

		return nil, err
	}

	xlog.Infof(ctx, "Creating cache for %s", key)

	// adding results to the cache
	if err := s.cache.Create(ctx, key, cached); err != nil {
		xlog.Errorf(ctx, "Unable to create tag cache, err: %s", err.Error())
	}

	return cached, nil
}

// Create ...
func (s *userService) Create(ctx context.Context, password string, user *models.User) (*models.User, error) {
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	user, err := s.repo.Create(ctx, user)

	// Clear the cache
	if err := s.cache.Flush(ctx); err != nil {
		xlog.Errorf(ctx, "Flushing cache error: %s", err.Error())
	}

	return user, err
}

// Update ...
func (s *userService) Update(ctx context.Context, user *models.User) (*models.User, error) {
	user, err := s.repo.Update(ctx, user)
	if err != nil {
		xlog.Errorf(ctx, "Unable update user, err: %s", err.Error())

		return nil, err
	}

	if err := s.queue.AddObject(ctx, "user-profile-updated", user); err != nil {
		xlog.Errorf(ctx, "Unable to send request into a 'user-password-reset' queue, err: %s", err.Error())
	}

	// Clear the cache
	if err := s.cache.Flush(ctx); err != nil {
		xlog.Errorf(ctx, "Flushing cache error: %s", err.Error())
	}

	return user, nil
}

// Delete ...
func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Clear the cache
	if err := s.cache.Flush(ctx); err != nil {
		xlog.Errorf(ctx, "Flushing cache error: %s", err.Error())
	}

	return nil
}

// UpdatePassword ...
func (s *userService) UpdatePassword(ctx context.Context, id uuid.UUID, newPassword string) (*models.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, models.ErrUserNotFound
	}

	// just in case set password failed
	if err := user.SetPassword(newPassword); err != nil {
		return nil, err
	}

	// drop reset hash
	user.PasswordResetHash = ""

	user, err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	if err := s.queue.AddObject(ctx, "user-password-changed", user); err != nil {
		xlog.Errorf(ctx, "Unable to send request into a 'user-password-changed' queue, err: %s", err.Error())
	}

	// Clear the cache
	if err := s.cache.Flush(ctx); err != nil {
		xlog.Errorf(ctx, "Flushing cache error: %s", err.Error())
	}

	return user, nil
}

// UpdateUsername ...
func (s *userService) UpdateUsername(ctx context.Context, id uuid.UUID, newUsername string) (*models.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, models.ErrUserNotFound
	}

	user.Email = newUsername

	user, err = s.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	// Clear the cache
	if err := s.cache.Flush(ctx); err != nil {
		xlog.Errorf(ctx, "Flushing cache error: %s", err.Error())
	}

	return user, err
}

// UpdateLogin ...
func (s *userService) UpdateLogin(ctx context.Context, user *models.User) (*models.User, error) {
	user.LastLogin = time.Now()

	user, err := s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	// Clear the cache
	if err := s.cache.Flush(ctx); err != nil {
		xlog.Errorf(ctx, "Flushing cache error: %s", err.Error())
	}

	return user, nil
}

// GetByResetHash ...
func (s *userService) GetByResetHash(ctx context.Context, hash string) (*models.User, error) {
	return s.repo.FindByResetHash(ctx, hash)
}

// ResetPassword ...
func (s *userService) ResetPassword(ctx context.Context, username string) (*models.User, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		xlog.Debugf(ctx, "Unable to find user %s, err: %s", username, err.Error())

		return nil, models.ErrUserNotFound
	}

	user.GeneratePasswordResetHash()
	user.PasswordResetAt = time.Now()

	if _, err = s.repo.Update(ctx, user); err != nil {
		xlog.Errorf(ctx, "Unable update user, err: %s", err.Error())

		return nil, err
	}

	if err := s.queue.AddObject(ctx, "user-password-reset", user); err != nil {
		xlog.Criticalf(ctx, "Unable to send request into a 'user-password-reset' queue, err: %s", err.Error())

		return nil, err
	}

	// Clear the cache
	if err := s.cache.Flush(ctx); err != nil {
		xlog.Errorf(ctx, "Flushing cache error: %s", err.Error())
	}

	return user, nil
}
