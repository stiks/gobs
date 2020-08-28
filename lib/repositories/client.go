package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/stiks/gobs/lib/models"
)

// ClientRepository ...
type ClientRepository interface {
	CountAll(ctx context.Context, params *models.ClientQueryParams) (int, error)
	FindAll(ctx context.Context, params *models.ClientQueryParams) ([]models.Client, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Client, error)
	Create(ctx context.Context, data *models.Client) (*models.Client, error)
	Update(ctx context.Context, data *models.Client) (*models.Client, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
