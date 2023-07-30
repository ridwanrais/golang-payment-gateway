package service

import (
	"github.com/ridwanrais/golang-payment-gateway/internal/repository"
)

type service struct {
	repo repository.Repository
}

type Service interface {
}

func NewService(r repository.Repository) Service {
	return &service{repo: r}
}
