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

	// Transaction
	GetVaTransaction(ctx context.Context, virtualAccountUuuid string) (*entity.Transaction, *entity.VirtualAccountTransaction, error)
	UpdateVaTransaction(ctx context.Context, updateData entity.UpdateVaRequest) (*entity.UpdateVaResponse, error)
	DeleteVaTransaction(ctx context.Context, vaTransactionUUID string, softDelete bool) error

	// BRI
	InsertBrivaTransaction(ctx context.Context, data entity.BrivaData, referenceNumber, vaNumber string) (*entity.CreateBrivaResponse, error)
}

func NewRepository(d *pgxpool.Pool, r *redis.Client) Repository {
	return &repository{db: d, redis: r}
}
