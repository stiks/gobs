package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/repositories"
)

type clientService struct {
	repo repositories.ClientRepository
}

// ClientService ...
type ClientService interface {
	CountAll(ctx context.Context, params *models.ClientQueryParams) (int, error)
	GetAll(ctx context.Context, params *models.ClientQueryParams) ([]models.Client, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Client, error)
	Create(ctx context.Context, client *models.Client) (*models.Client, error)
	Update(ctx context.Context, client *models.Client) (*models.Client, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// NewClientService ...
func NewClientService(repo repositories.ClientRepository) ClientService {
	return &clientService{
		repo: repo,
	}
}

// GetAll ...
func (s *clientService) GetAll(ctx context.Context, params *models.ClientQueryParams) ([]models.Client, error) {
	return s.repo.FindAll(ctx, params)
}

// CountAll ...
func (s *clientService) CountAll(ctx context.Context, params *models.ClientQueryParams) (int, error) {
	return s.repo.CountAll(ctx, params)
}

// GetByID ...
func (s *clientService) GetByID(ctx context.Context, id uuid.UUID) (*models.Client, error) {
	return s.repo.FindByID(ctx, id)
}

// Create ...
func (s *clientService) Create(ctx context.Context, client *models.Client) (*models.Client, error) {
	client.CreatedAt = time.Now()
	client.UpdatedAt = time.Now()

	return s.repo.Create(ctx, client)
}

// Update ...
func (s *clientService) Update(ctx context.Context, client *models.Client) (*models.Client, error) {
	client.UpdatedAt = time.Now()

	return s.repo.Update(ctx, client)
}

// Delete ...
func (s *clientService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
