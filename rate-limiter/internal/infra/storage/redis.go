package storage

import (
	"context"
	"errors"
	"github.com/DiegoOpenheimer/go/rate-limiter/internal/usecases/ports"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{
		client: client,
	}
}

func (rs *RedisStorage) SaveWithContext(ctx context.Context, key string, value any, options *ports.StorageOptions) error {
	var expire time.Duration
	if options != nil {
		expire = options.Expire
	}
	return rs.client.Set(ctx, key, value, expire).Err()
}

func (rs *RedisStorage) GetWithContext(ctx context.Context, key string) (string, error) {
	result, err := rs.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return result, err
}

func (rs *RedisStorage) IncrementWithContext(ctx context.Context, key string, options *ports.StorageOptions) error {
	pipe := rs.client.Pipeline()
	pipe.Incr(ctx, key)
	if options != nil {
		pipe.Expire(ctx, key, options.Expire)
	}
	_, err := pipe.Exec(ctx)
	return err
}
