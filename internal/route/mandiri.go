package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ridwanrais/golang-payment-gateway/internal/controller"
)

func SetupMandiriRoutes(router *gin.RouterGroup, controller controller.Controllers) {
	auth := router.Group("/mandiri")
	{
		auth.POST("/va", controller.MandiriCreateVA)
	}
}
