package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ridwanrais/golang-payment-gateway/internal/controller"
)

func SetupBriRoutes(router *gin.RouterGroup, controller controller.Controllers) {
	auth := router.Group("/bri")
	{
		auth.POST("/briva", controller.BriCreateBriva)
	}
}
