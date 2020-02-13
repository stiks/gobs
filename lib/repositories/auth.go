package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/stiks/gobs/lib/models"
)

// AuthRepository ...
type AuthRepository interface {
	FindByClientUser(ctx context.Context, clientID uuid.UUID, userID uuid.UUID) (*models.Token, error)
	FindByClientID(ctx context.Context, clientID string) (*models.AuthClient, error)
	FindByHashClient(ctx context.Context, clientID uuid.UUID, token string) (*models.Token, error)
	FindUserByUsername(ctx context.Context, username string) (*models.User, error)
	FindUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	CreateToken(ctx context.Context, data *models.Token) (*models.Token, error)
	DeleteToken(ctx context.Context, id uuid.UUID) error
}
