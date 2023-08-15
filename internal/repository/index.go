package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

type Repository interface {
	// caching
	SetCacheValue(ctx context.Context, key, value string, expiration time.Duration) error
	GetCacheValue(ctx context.Context, key string) (string, error)
}

func NewRepository(d *pgxpool.Pool, r *redis.Client) Repository {
	return &repository{db: d, redis: r}
}
