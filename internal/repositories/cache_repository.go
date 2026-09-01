package repositories

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

type redisCacheRepository struct {
	client *redis.Client
}

// @inject
func NewRedisCacheRepository(client *redis.Client) CacheRepository {
	return &redisCacheRepository{client: client}
}

func (r *redisCacheRepository) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *redisCacheRepository) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}
