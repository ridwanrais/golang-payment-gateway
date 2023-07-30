package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ridwanrais/golang-payment-gateway/internal/controller"
)

func SetupHealthsRoutes(router *gin.RouterGroup, controller controller.Controllers) {
	health := router.Group("/health")
	{
		health.GET("", controller.GetHealth)
	}
}
