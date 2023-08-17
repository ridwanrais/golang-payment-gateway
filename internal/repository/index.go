package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ridwanrais/golang-payment-gateway/internal/entity"
)

type repository struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

type Repository interface {
	// caching
	SetCacheValue(ctx context.Context, key, value string, expiration time.Duration) error
	GetCacheValue(ctx context.Context, key string) (string, error)

	// BRI
	InsertBrivaTransaction(ctx context.Context, data entity.BrivaData, referenceNumber, vaNumber string) (int, error)
}

func NewRepository(d *pgxpool.Pool, r *redis.Client) Repository {
	return &repository{db: d, redis: r}
}
