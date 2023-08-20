package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ridwanrais/golang-payment-gateway/internal/entity"
	"github.com/ridwanrais/golang-payment-gateway/internal/utils"
	"github.com/ridwanrais/golang-payment-gateway/internal/validator"
)

func (c *controllers) BriCreateBriva(ctx *gin.Context) {
	requestData, err := validator.BriCreateBrivaValidator(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, entity.Response{
			Status:       false,
			ResponseCode: http.StatusBadRequest,
			Message:      fmt.Sprintf("validation error: %s", err.Error()),
		})
		return
	}

	response, err := c.service.BriCreateBriva(ctx, *requestData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, entity.Response{
			Status:       false,
			ResponseCode: http.StatusBadRequest,
			Message:      fmt.Sprintf("create briva error: %s", err.Error()),
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
		Data: data,
	})
}

func (c *controllers) BriGetBriva(ctx *gin.Context) {
	vaUuid := ctx.Param("vaUuid")
	if vaUuid == "" {
		ctx.JSON(http.StatusBadRequest, entity.Response{
			Status:       false,
			ResponseCode: http.StatusBadRequest,
			Message:      "validation error: va uuid is required",
		})
		return
	}

	response, err := c.service.BriGetBriva(ctx, vaUuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, entity.Response{
			Status:       false,
			ResponseCode: http.StatusBadRequest,
			Message:      fmt.Sprintf("get briva error: %s", err.Error()),
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
		ResponseCode: http.StatusOK,
		Message:      "ok",
		Data: data,
	})
}

func (c *controllers) BriUpdateBriva(ctx *gin.Context) {
	requestData, err := validator.BriUpdateBrivaValidator(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, entity.Response{
			Status:       false,
			ResponseCode: http.StatusBadRequest,
			Message:      fmt.Sprintf("validation error: %s", err.Error()),
		})
		return
	}

	response, err := c.service.BriUpdateBriva(ctx, *requestData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, entity.Response{
			Status:       false,
			ResponseCode: http.StatusBadRequest,
			Message:      fmt.Sprintf("update briva error: %s", err.Error()),
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
		ResponseCode: http.StatusOK,
		Message:      "ok",
		Data: data,
	})
}

func (c *controllers) BriDeleteBriva(ctx *gin.Context) {
	vaUuid := ctx.Param("vaUuid")
	if vaUuid == "" {
		ctx.JSON(http.StatusBadRequest, entity.Response{
			Status:       false,
			ResponseCode: http.StatusBadRequest,
			Message:      "validation error: va uuid is required",
		})
		return
	}

	err := c.service.BriDeleteBriva(ctx, vaUuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, entity.Response{
			Status:       false,
			ResponseCode: http.StatusBadRequest,
			Message:      fmt.Sprintf("delete briva error: %s", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusCreated, entity.Response{
		Status:       true,
		ResponseCode: http.StatusOK,
		Message:      "ok",
	})
}
