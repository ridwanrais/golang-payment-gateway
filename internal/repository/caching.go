package repository

import (
	"context"
	"fmt"
	"time"
)

func (r *repository) setAccessToken(ctx context.Context, clientID, token string, expirationTime time.Time) error {
	duration := expirationTime.Sub(time.Now())
	return r.redis.Set(ctx, fmt.Sprintf("access_token:%s", clientID), token, duration).Err()
}

func (r *repository) getAccessToken(ctx context.Context, clientID string) (string, error) {
	return r.redis.Get(ctx, fmt.Sprintf("access_token:%s", clientID)).Result()
}

func (r *repository) SetCacheValue(ctx context.Context, key, value string, expiration time.Duration) error {
	return r.redis.Set(ctx, key, value, expiration).Err()
}

func (r *repository) GetCacheValue(ctx context.Context, key string) (string, error) {
	return r.redis.Get(ctx, key).Result()
}
