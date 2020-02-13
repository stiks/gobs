package services

import (
	"context"

	"github.com/stiks/gobs/lib/repositories"
)

type cacheService struct {
	repo repositories.CacheRepository
}

// CacheService ...
type CacheService interface {
	GetByKey(ctx context.Context, key string, obj interface{}) error
	Create(ctx context.Context, key string, data interface{}) error
	Update(ctx context.Context, key string, data interface{}) error
	Delete(ctx context.Context, key string) error
	Flush(ctx context.Context) error
}

// NewCacheService ...
func NewCacheService(repo repositories.CacheRepository) CacheService {
	return &cacheService{
		repo: repo,
	}
}

func (s *cacheService) GetByKey(ctx context.Context, key string, obj interface{}) error {
	return s.repo.FindByKey(ctx, key, obj)
}

func (s *cacheService) Create(ctx context.Context, key string, data interface{}) error {
	return s.repo.Create(ctx, key, data)
}

func (s *cacheService) Update(ctx context.Context, key string, data interface{}) error {
	return s.repo.Update(ctx, key, data)
}

func (s *cacheService) Delete(ctx context.Context, key string) error {
	return s.repo.Delete(ctx, key)
}

func (s *cacheService) Flush(ctx context.Context) error {
	return s.repo.Flush(ctx)
}
