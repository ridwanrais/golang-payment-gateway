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
	BriCreateBriva(ctx context.Context, RequestData entity.BrivaData) (*entity.BriResponseData, error) 
}

func NewService(r repository.Repository) Service {
	return &service{repo: r}
}
