package mock

import (
	"context"
	"log"
	"net/url"

	"github.com/stiks/gobs/lib/repositories"
)

type queueRepository struct {
}

// NewQueueRepository ...
func NewQueueRepository() repositories.QueueRepository {
	return &queueRepository{}
}

// CountAll ...
func (r *queueRepository) Add(ctx context.Context, queue string, data []byte) error {
	log.Printf("Mock queue service")

	return nil
}

// AddObject ...
func (r *queueRepository) AddObject(ctx context.Context, queue string, data interface{}) error {
	log.Printf("Mock queue service")

	return nil
}

// AddToURL ...
func (r *queueRepository) AddToURL(ctx context.Context, queue string, data url.Values) error {
	log.Printf("Mock queue service")

	return nil
}
