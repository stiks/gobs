package repositories

import (
	"context"
	"net/url"
)

// QueueRepository ...
type QueueRepository interface {
	Add(ctx context.Context, queue string, data []byte) error
	AddObject(ctx context.Context, queue string, data interface{}) error
	AddToURL(ctx context.Context, queue string, data url.Values) error
}
