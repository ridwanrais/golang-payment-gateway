package controller

import (
	"github.com/gin-gonic/gin"
)

type controllers struct {
	// usecase usecase.Usecases
}

type Controllers interface {
	// health
	GetHealth(c *gin.Context)
}

func NewControllers() Controllers {
	return &controllers{}
	// return &controllers{
	// 	usecase: config.GetUsecase()}
}
