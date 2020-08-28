package mock

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/repositories"
	"github.com/stiks/gobs/pkg/helpers"
)

var (
	_clientsList = []models.Client{
		{
			ID:        helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
			Name:      "Test Client",
			Email:     "peter@test.com",
			Status:    models.StatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Email:     "john@snow.com",
			Name:      "John Snow",
			OwnerID:   helpers.UUIDFromString(nil, "775a5b37-1742-4e54-9439-0357e768b011"),
			Status:    models.StatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
)

// NewClientRepository ...
func NewClientRepository() repositories.ClientRepository {
	return &clientRepository{
		db: _clientsList,
	}
}

type clientRepository struct {
	db []models.Client
}

// FindAll ...
func (r *clientRepository) FindAll(ctx context.Context, params *models.ClientQueryParams) ([]models.Client, error) {
	return r.db, nil
}

// CountAll ...
func (r *clientRepository) CountAll(ctx context.Context, params *models.ClientQueryParams) (int, error) {
	return len(r.db), nil
}

// FindByID ...
func (r *clientRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Client, error) {
	for _, key := range r.db {
		if key.ID == id {
			return &key, nil
		}
	}

	return nil, models.ErrClientNotFound
}

// Create ...
func (r *clientRepository) Create(ctx context.Context, data *models.Client) (*models.Client, error) {
	for _, key := range r.db {
		if key.Email == data.Email {
			return nil, models.ErrClientNameTaken
		}
	}

	r.db = append(r.db, *data)

	return data, nil
}

// Update ...
func (r *clientRepository) Update(ctx context.Context, data *models.Client) (*models.Client, error) {
	for k, i := range r.db {
		if i.ID == data.ID {
			r.db[k] = *data

			return data, nil
		}
	}

	return nil, models.ErrClientNotFound
}

// Delete ...
func (r *clientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}

	found := false
	var db []models.Client
	for _, k := range r.db {
		if k.ID != id {
			db = append(db, k)
		}

		found = true
	}

	r.db = db

	if found {
		return nil
	}

	return models.ErrClientNotFound
}
