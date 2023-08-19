package repository

import (
	"context"
	"time"
)

func (r *repository) SetCacheValue(ctx context.Context, key, value string, expiration time.Duration) error {
	return r.redis.Set(ctx, key, value, expiration).Err()
}

func (r *repository) GetCacheValue(ctx context.Context, key string) (string, error) {
	return r.redis.Get(ctx, key).Result()
}
