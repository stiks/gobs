package repositories

import "context"

// CacheRepository ...
type CacheRepository interface {
	FindByKey(ctx context.Context, key string, obj interface{}) error
	Create(ctx context.Context, key string, data interface{}) error
	Update(ctx context.Context, key string, data interface{}) error
	Delete(ctx context.Context, key string) error
	Flush(ctx context.Context) error
}
