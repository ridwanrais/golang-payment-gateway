package validator

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ridwanrais/golang-payment-gateway/internal/entity"
)

func BriCreateBrivaValidator(c *gin.Context) (*entity.BrivaData, error) {
	var requestData entity.BrivaData
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
		validation.Field(&requestData.InstitutionCode, validation.Required),
		validation.Field(&requestData.BrivaNo, validation.Required),
		validation.Field(&requestData.CustCode, validation.Required),
		validation.Field(&requestData.Nama, validation.Required),
		validation.Field(&requestData.Amount, validation.Required),
	); err != nil {
		// Validation failed, handle the error
		return nil, err
	}

	return &requestData, nil
}
