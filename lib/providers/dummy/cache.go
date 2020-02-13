package dummy

import (
	"context"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/repositories"
)

type cacheRepository struct{}

// NewCacheRepository ...
func NewCacheRepository() repositories.CacheRepository {
	return &cacheRepository{}
}

func (r *cacheRepository) FindByKey(ctx context.Context, key string, obj interface{}) error {
	return models.ErrMissCache
}

func (r *cacheRepository) Create(ctx context.Context, key string, data interface{}) error {
	return nil
}

func (r *cacheRepository) Update(ctx context.Context, key string, data interface{}) error {
	return nil
}

func (r *cacheRepository) Delete(ctx context.Context, key string) error {
	return nil
}

func (r *cacheRepository) Flush(ctx context.Context) error {
	return nil
}
