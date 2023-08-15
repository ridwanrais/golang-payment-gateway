package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/ridwanrais/golang-payment-gateway/internal/config"
	"github.com/ridwanrais/golang-payment-gateway/internal/service"
)

type controllers struct {
	service service.Service
}

type Controllers interface {
	// health
	GetHealth(c *gin.Context)

	// BRI
	BriCreateBriva(ctx *gin.Context)
}

func NewControllers() Controllers {
	return &controllers{
		service: config.GetService(),
	}
}
