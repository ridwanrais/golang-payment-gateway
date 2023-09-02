package route

import (
	"github.com/gin-gonic/gin"
	"github.com/ridwanrais/golang-payment-gateway/internal/controller"
)

func SetupRoutes(router *gin.Engine) {
	controller := controller.NewControllers()

	v1 := router.Group("/v1")

	{
		SetupHealthsRoutes(v1, controller)
		SetupBriRoutes(v1, controller)
		// SetupMandiriRoutes(v1, controller) // not ready yet	
	}
}
