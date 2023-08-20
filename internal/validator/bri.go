package validator

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ridwanrais/golang-payment-gateway/internal/entity"
)

func BriCreateBrivaValidator(c *gin.Context) (*entity.CreateBrivaRequest, error) {
	var requestData entity.CreateBrivaRequest
	if err := c.ShouldBind(&requestData); err != nil {
		// Validation failed, handle the error
		if verr, ok := err.(validation.Errors); ok {
			// Validation errors occurred
			return nil, verr
		}
		// Other errors occurred
		return nil, err
	}

	if err := validation.ValidateStruct(&requestData,
		validation.Field(&requestData.Name, validation.Required),
		validation.Field(&requestData.PhoneNumber, validation.Required),
		validation.Field(&requestData.Amount, validation.Required),
	); err != nil {
		// Validation failed, handle the error
		return nil, err
	}

	return &requestData, nil
}

func BriUpdateBrivaValidator(c *gin.Context) (*entity.UpdateVaRequest, error) {
	var requestData entity.UpdateVaRequest

	requestData.VaTransactionUUID = c.Param("vaUuid")

	if err := c.ShouldBind(&requestData); err != nil {
		// Validation failed, handle the error
		if verr, ok := err.(validation.Errors); ok {
			// Validation errors occurred
			return nil, verr
		}
		// Other errors occurred
		return nil, err
	}

	if err := validation.ValidateStruct(&requestData,
		validation.Field(&requestData.VaTransactionUUID, validation.Required),
		validation.Field(&requestData.Name, validation.Required),
		validation.Field(&requestData.PhoneNumber, validation.Required),
		validation.Field(&requestData.Amount, validation.Required),
		// validation.Field(&requestData.PaymentStatus, validation.Required, validation.In(constants.PAYMENT_PENDING, constants.PAYMENT_COMPLETED, constants.PAYMENT_FAILED)),
		validation.Field(&requestData.ExpiryDate, validation.Required),
	); err != nil {
		// Validation failed, handle the error
		return nil, err
	}

	return &requestData, nil
}
