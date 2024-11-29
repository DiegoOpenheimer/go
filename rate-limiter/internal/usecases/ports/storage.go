package ports

import (
	"context"
	"time"
)

type StorageOptions struct {
	Expire time.Duration
}

type Storage interface {
	SaveWithContext(ctx context.Context, key string, value any, options *StorageOptions) error
	GetWithContext(ctx context.Context, key string) (string, error)
	IncrementWithContext(ctx context.Context, key string, options *StorageOptions) error
}
