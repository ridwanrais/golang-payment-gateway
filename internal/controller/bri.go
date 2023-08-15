package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ridwanrais/golang-payment-gateway/internal/validator"
)

func (c *controllers) BriCreateBriva(ctx *gin.Context) {
	brivaData, err := validator.BriCreateBrivaValidator(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"validation error": err.Error()})
		return
	}

	response, err := c.service.BriCreateBriva(ctx, *brivaData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"create briva error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "ok",
		"data":    response,
	})
}
