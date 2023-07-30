package controller

import (
	"github.com/gin-gonic/gin"
)

func (c *controllers) GetHealth(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "ok",
	})
}