package services

import (
	"context"
	"net/url"

	"github.com/stiks/gobs/lib/repositories"
)

type queueService struct {
	repo repositories.QueueRepository
}

// QueueService ...
type QueueService interface {
	Add(ctx context.Context, queue string, data []byte) error
	AddObject(ctx context.Context, queue string, data interface{}) error
	AddToURL(ctx context.Context, queue string, data url.Values) error
}

// NewQueueService ...
func NewQueueService(repo repositories.QueueRepository) QueueService {
	return &queueService{
		repo: repo,
	}
}

// Add ...
func (s *queueService) Add(ctx context.Context, queue string, data []byte) error {
	return s.repo.Add(ctx, queue, data)
}

// AddObject ...
func (s *queueService) AddObject(ctx context.Context, queue string, data interface{}) error {
	return s.repo.AddObject(ctx, queue, data)
}

// AddToURL ...
func (s *queueService) AddToURL(ctx context.Context, queue string, data url.Values) error {
	return s.repo.AddToURL(ctx, queue, data)
}
