package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ridwanrais/golang-payment-gateway/internal/entity"
	"github.com/ridwanrais/golang-payment-gateway/internal/utils"
	"github.com/ridwanrais/golang-payment-gateway/internal/validator"
)

func (c *controllers) MandiriCreateVA(ctx *gin.Context) {
	requestData, err := validator.CreateVirtualAccountValidator(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, entity.Response{
			Status:       false,
			ResponseCode: http.StatusBadRequest,
			Message:      fmt.Sprintf("validation error: %s", err.Error()),
		})
		return
	}

	response, err := c.service.MandiriCreateVA(ctx, *requestData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, entity.Response{
			Status:       false,
			ResponseCode: http.StatusBadRequest,
			Message:      fmt.Sprintf("create mandiri va error: %s", err.Error()),
		})
		return
	}

	data, err := utils.StructToMap(response)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, entity.Response{
			Status:       false,
			ResponseCode: http.StatusBadRequest,
			Message:      fmt.Sprintf("parsing response error: %s", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusCreated, entity.Response{
		Status:       true,
		ResponseCode: http.StatusCreated,
		Message:      "ok",
		Data:         data,
	})
}
