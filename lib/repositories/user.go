package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/stiks/gobs/lib/models"
)

// UserRepository ...
type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByResetHash(ctx context.Context, hash string) (*models.User, error)
	CountAll(ctx context.Context, params *models.UserQueryParams) (int, error)
	FindAll(ctx context.Context, params *models.UserQueryParams) ([]models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Create(ctx context.Context, data *models.User) (*models.User, error)
	Update(ctx context.Context, data *models.User) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
