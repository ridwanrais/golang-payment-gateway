package config

import (
	"log"
	"sync"

	"github.com/ridwanrais/golang-payment-gateway/internal/repository"
	"github.com/ridwanrais/golang-payment-gateway/internal/service"
)

var oneUc sync.Once
var uc service.Service

func GetService() service.Service {
	oneUc.Do(func() {
		uc = service.NewService(
			getRepository(),
		)
	})

	return uc
}

var repo repository.Repository
var oneRepo sync.Once

func getRepository() repository.Repository {
	pgPool, err := ConnectToPostgreSQL()
	if err != nil {
		log.Fatal(err.Error())
	}

	oneRepo.Do(func() {
		repo = repository.NewRepository(pgPool)
	})

	return repo
}
