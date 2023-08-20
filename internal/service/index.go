package service

import (
	"context"

	"github.com/ridwanrais/golang-payment-gateway/internal/entity"
	"github.com/ridwanrais/golang-payment-gateway/internal/repository"
)

type service struct {
	repo repository.Repository
}

type Service interface {
	// token
	BriRetrieveAccessToken(ctx context.Context) (string, error)

	// BRI
	BriCreateBriva(ctx context.Context, RequestData entity.CreateBrivaRequest) (*entity.CreateBrivaResponse, error)
	BriGetBriva(ctx context.Context, vaUuid string) (*entity.GetVirtualAccountResponse, error)
	BriUpdateBriva(ctx context.Context, RequestData entity.UpdateVaRequest) (*entity.UpdateVaResponse, error)
	BriDeleteBriva(ctx context.Context, vaUuid string) error
}

func NewService(r repository.Repository) Service {
	return &service{repo: r}
}
