package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ridwanrais/golang-payment-gateway/internal/entity"
	"github.com/ridwanrais/golang-payment-gateway/internal/utils"
)

func (s *service) MandiriCreateVA(ctx context.Context, request entity.CreateVaRequest) (*entity.CreateVaResponse, error) {
	timestamp := utils.MandiriGenerateCurrentTimestamp()

	token, err := s.MandiriRetrieveAccessToken(ctx)
	if err != nil {
		return nil, errors.New("failed to retrieve access token: " + err.Error())
	}

	path := "/openapi/transaction/v1.0/transfer-va/create-va"
	method := "POST"
	partnerServiceId := os.Getenv("MANDIRI_PARTNER_SERVICE_ID")
	referenceNumber := utils.GenerateReferenceNumber()
	expiryDuration, err := strconv.Atoi(os.Getenv("MANDIRI_VA_EXPIRY_DURATION_IN_HOURS"))
	if err != nil {
		return nil, errors.New("failed to convert MANDIRI_VA_EXPIRY_DURATION_IN_HOURS to integer")
	}

	expiryDate := time.Now().Add(time.Duration(expiryDuration) * time.Hour).Format("2006-01-02T15:04:05-0700")
	mandiriVaData := entity.MandiriVaData{
		PartnerServiceId:   partnerServiceId,
		CustomerNo:         request.PhoneNumber,
		VirtualAccountNo:   partnerServiceId + request.PhoneNumber,
		VirtualAccountName: request.Name,
		TrxId:              referenceNumber,
		TotalAmount: entity.Amount{
			Value:    strconv.Itoa(request.Amount),
			Currency: os.Getenv("MANDIRI_CURRENCY"),
		},
		BillDetails: []entity.BillDetail{
			{
				BillAmount: entity.Amount{
					Value:    strconv.Itoa(request.Amount),
					Currency: os.Getenv("MANDIRI_CURRENCY"),
				},
			},
		},
		ExpiredDate: expiryDate,
	}

	signature, err := utils.MandiriGenerateTransactionSignature(path, method, token, &mandiriVaData, timestamp)
	if err != nil {
		return nil, errors.New("failed to generate signature: " + err.Error())
	}

	fmt.Println("")

	client := &http.Client{Timeout: 10 * time.Second}
	url := os.Getenv("MANDIRI_HOST") + "/openapi/transaction/v1.0/transfer-va/create-va"
	jsonData, err := json.Marshal(mandiriVaData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("MANDIRI-Timestamp", timestamp)
	req.Header.Set("MANDIRI-Signature", signature)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseData entity.MandiriResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}

	successResponse := len(responseData.ResponseMessage) >= 3 && responseData.ResponseMessage[:3] == "200"
	if !successResponse {
		return nil, errors.New(responseData.ResponseMessage)
	}

	refNumber := utils.GenerateReferenceNumber()

	request.ExpiryDate = expiryDate
	request.ReferenceNumber = refNumber
	request.VirtualAccountNumber = mandiriVaData.VirtualAccountNo
	response, err := s.repo.CreateVaTransaction(ctx, request)
	if err != nil {
		return nil, err
	}

	return &entity.CreateVaResponse{
		ReferenceNumber:      refNumber,
		VirtualAccountNumber: mandiriVaData.VirtualAccountNo,
		TransactionUUID:      response.TransactionUUID,
		VaTransactionUUID:    response.VaTransactionUUID,
	}, nil
}
