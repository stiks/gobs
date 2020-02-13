package appengine

import (
	"context"

	"google.golang.org/appengine/memcache"

	"github.com/stiks/gobs/lib/models"
	"github.com/stiks/gobs/lib/repositories"
	"github.com/stiks/gobs/pkg/xlog"
)

type cacheRepository struct {
}

// NewCacheRepository ...
func NewCacheRepository() repositories.CacheRepository {
	return &cacheRepository{}
}

func (r *cacheRepository) FindByKey(ctx context.Context, key string, obj interface{}) error {
	_, err := memcache.Gob.Get(ctx, key, obj)
	// item not found
	if err != nil && err == memcache.ErrCacheMiss {
		xlog.Infof(ctx, "Cache miss")

		return models.ErrMissCache
	}

	// any other error
	if err != nil {
		xlog.Errorf(ctx, "Find by key error: %s", err.Error())

		return err
	}

	return nil
}

func (r *cacheRepository) Create(ctx context.Context, key string, data interface{}) error {
	return r.Update(ctx, key, data)
}

func (r *cacheRepository) Update(ctx context.Context, key string, data interface{}) error {
	item := &memcache.Item{
		Key:        key,
		Object:     data,
		Expiration: 0,
	}

	return memcache.Gob.Set(ctx, item)
}

func (r *cacheRepository) Delete(ctx context.Context, key string) error {
	return memcache.Delete(ctx, key)
}

func (r *cacheRepository) Flush(ctx context.Context) error {
	return memcache.Flush(ctx)
}
