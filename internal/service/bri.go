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
	"strings"
	"time"

	"github.com/ridwanrais/golang-payment-gateway/internal/constants"
	"github.com/ridwanrais/golang-payment-gateway/internal/entity"
	"github.com/ridwanrais/golang-payment-gateway/internal/utils"
)

func (s *service) BriCreateBriva(ctx context.Context, RequestData entity.CreateBrivaRequest) (*entity.CreateBrivaResponse, error) {
	timestamp := utils.GenerateCurrentTimestamp()

	token, err := s.BriRetrieveAccessToken(ctx)
	if err != nil {
		return nil, errors.New("failed to retrieve access token: " + err.Error())
	}

	path := "/v1/briva"
	method := "POST"
	expiryDuration, err := strconv.Atoi(os.Getenv("BRI_BRIVA_EXPIRY_DURATION_IN_HOURS"))
	if err != nil {
		return nil, errors.New("failed to convert BRI_BRIVA_EXPIRY_DURATION_IN_HOURS to integer")
	}

	brivaData := entity.BrivaData{
		InstitutionCode: os.Getenv("BRI_INSTITUTION_CODE"),
		BrivaNo:         os.Getenv("BRI_BRIVA_NUMBER"),
		CustCode:        RequestData.PhoneNumber,
		Nama:            RequestData.Name,
		Amount:          strconv.Itoa(RequestData.Amount),
		Keterangan:      RequestData.Note,
		ExpiredDate:     time.Now().Add(time.Duration(expiryDuration) * time.Hour).Format("2006-01-02 15:04:05"),
	}
	signature, err := utils.BriGenerateSignature(path, method, &brivaData, token, timestamp)
	if err != nil {
		return nil, errors.New("failed to generate signature: " + err.Error())
	}

	client := &http.Client{Timeout: 10 * time.Second}
	url := os.Getenv("BRI_HOST") + "/v1/briva"
	jsonData, err := json.Marshal(brivaData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("BRI-Timestamp", timestamp)
	req.Header.Set("BRI-Signature", signature)
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

	var responseData entity.BriResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}

	if !responseData.Status {
		return nil, errors.New(responseData.ErrDesc)
	}

	refNumber := utils.GenerateReferenceNumber()

	// Construct virtual account number
	virtualAccountNumber := brivaData.BrivaNo + brivaData.CustCode

	response, err := s.repo.InsertBrivaTransaction(ctx, brivaData, refNumber, virtualAccountNumber)
	if err != nil {
		return nil, err
	}

	return &entity.CreateBrivaResponse{
		ReferenceNumber:      refNumber,
		VirtualAccountNumber: brivaData.BrivaNo + brivaData.CustCode,
		TransactionUUID:      response.TransactionUUID,
		VaTransactionUUID:    response.VaTransactionUUID,
	}, nil
}

func (s *service) BriGetBriva(ctx context.Context, vaUuid string) (*entity.GetVirtualAccountResponse, error) {
	trx, vaTrx, err := s.repo.GetVaTransaction(ctx, vaUuid)
	if err != nil {
		return nil, errors.New("failed to get va transaction: " + err.Error())
	}

	timestamp := utils.GenerateCurrentTimestamp()

	token, err := s.BriRetrieveAccessToken(ctx)
	if err != nil {
		return nil, errors.New("failed to retrieve access token: " + err.Error())
	}

	client := &http.Client{Timeout: 10 * time.Second}
	institutionCode := os.Getenv("BRI_INSTITUTION_CODE")
	customerBrivaNumber := vaTrx.VirtualAccountNumber
	brivaNo := customerBrivaNumber[:5]
	customerCode := customerBrivaNumber[5:]
	path := "/v1/briva/" + institutionCode + "/" + brivaNo + "/" + customerCode
	url := os.Getenv("BRI_HOST") + path
	method := "GET"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("BRI-Timestamp", timestamp)
	signature, err := utils.BriGenerateSignature(path, method, nil, token, timestamp)
	if err != nil {
		return nil, errors.New("failed to generate signature: " + err.Error())
	}
	req.Header.Set("BRI-Signature", signature)
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

	var responseData entity.BriResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}

	if !responseData.Status {
		return nil, errors.New(responseData.ErrDesc)
	}

	name := ""
	if trx.Description.Valid {
		name = trx.Name.String
	}
	description := ""
	if trx.Description.Valid {
		description = trx.Description.String
	}
	metadata := ""
	if vaTrx.Metadata.Valid {
		metadata = vaTrx.Metadata.String
	}
	paymentStatus := ""
	if responseData.Data.StatusBayar == "Y" {
		paymentStatus = constants.PAYMENT_COMPLETED
	} else if utils.IsLaterThanNow(vaTrx.ExpiryDate.Time) {
		paymentStatus = constants.PAYMENT_PENDING
	} else {
		paymentStatus = constants.PAYMENT_FAILED
	}

	return &entity.GetVirtualAccountResponse{
		PaymentStatus:        paymentStatus,
		Name:                 name,
		ReferenceNumber:      trx.ReferenceNumber,
		TransactionDate:      trx.TransactionDate,
		TransactionAmount:    int(trx.TransactionAmount),
		Description:          description,
		BankName:             vaTrx.BankName,
		VirtualAccountNumber: vaTrx.VirtualAccountNumber,
		ExpiryDate:           &vaTrx.ExpiryDate.Time,
		Metadata:             metadata,
	}, nil
}

func (s *service) BriUpdateBriva(ctx context.Context, RequestData entity.UpdateVaRequest) (*entity.UpdateVaResponse, error) {
	trx, vaTrx, err := s.repo.GetVaTransaction(ctx, RequestData.VaTransactionUUID)
	if err != nil {
		return nil, errors.New("failed to get va transaction: " + err.Error())
	}

	timestamp := utils.GenerateCurrentTimestamp()

	token, err := s.BriRetrieveAccessToken(ctx)
	if err != nil {
		return nil, errors.New("failed to retrieve access token: " + err.Error())
	}

	path := "/v1/briva"
	method := "PUT"
	expiryDuration, err := strconv.Atoi(os.Getenv("BRI_BRIVA_EXPIRY_DURATION_IN_HOURS"))
	if err != nil {
		return nil, errors.New("failed to convert BRI_BRIVA_EXPIRY_DURATION_IN_HOURS to integer")
	}
	customerBrivaNumber := vaTrx.VirtualAccountNumber
	brivaNo := customerBrivaNumber[:5]
	custCode := customerBrivaNumber[5:]

	brivaData := entity.BrivaData{
		InstitutionCode: os.Getenv("BRI_INSTITUTION_CODE"),
		BrivaNo:         brivaNo,
		CustCode:        custCode,
		Nama:            RequestData.Name,
		Amount:          strconv.Itoa(RequestData.Amount),
		Keterangan:      RequestData.Note,
		ExpiredDate:     time.Now().Add(time.Duration(expiryDuration) * time.Hour).Format("2006-01-02 15:04:05"),
	}
	signature, err := utils.BriGenerateSignature(path, method, &brivaData, token, timestamp)
	if err != nil {
		return nil, errors.New("failed to generate signature: " + err.Error())
	}

	client := &http.Client{Timeout: 10 * time.Second}
	url := os.Getenv("BRI_HOST") + "/v1/briva"
	jsonData, err := json.Marshal(brivaData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("BRI-Timestamp", timestamp)
	req.Header.Set("BRI-Signature", signature)
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

	var responseData entity.BriResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}

	if !responseData.Status {
		return nil, errors.New(responseData.ErrDesc)
	}

	// Construct virtual account number
	virtualAccountNumber := brivaData.BrivaNo + brivaData.CustCode

	_, err = s.repo.UpdateVaTransaction(ctx, entity.UpdateVaRequest{
		TransactionUUID:      trx.UUID,
		VaTransactionUUID:    vaTrx.UUID,
		Name:                 RequestData.Name,
		Amount:               RequestData.Amount,
		Note:                 RequestData.Note,
		BankName:             constants.BANK_NAME_BRI,
		VirtualAccountNumber: virtualAccountNumber,
		ExpiryDate:           RequestData.ExpiryDate,
	})
	if err != nil {
		return nil, err
	}

	return &entity.UpdateVaResponse{
		ReferenceNumber:      trx.ReferenceNumber,
		VirtualAccountNumber: virtualAccountNumber,
		TransactionUUID:      trx.UUID,
		VaTransactionUUID:    vaTrx.UUID,
	}, nil
}

func (s *service) BriDeleteBriva(ctx context.Context, vaUuid string) error {
	_, vaTrx, err := s.repo.GetVaTransaction(ctx, vaUuid)
	if err != nil {
		return errors.New("failed to get va transaction: " + err.Error())
	}

	timestamp := utils.GenerateCurrentTimestamp()

	token, err := s.BriRetrieveAccessToken(ctx)
	if err != nil {
		return errors.New("failed to retrieve access token: " + err.Error())
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Set your parameters
	institutionCode := os.Getenv("BRI_INSTITUTION_CODE")
	customerBrivaNumber := vaTrx.VirtualAccountNumber
	brivaNo := customerBrivaNumber[:5]
	custCode := customerBrivaNumber[5:]
	path := "/v1/briva"
	url := os.Getenv("BRI_HOST") + path
	method := "DELETE"

	// Construct the request payload
	payload := fmt.Sprintf("institutionCode=%s&brivaNo=%s&custCode=%s", institutionCode, brivaNo, custCode)
	payloadReader := strings.NewReader(payload)

	req, err := http.NewRequest(method, url, payloadReader)
	if err != nil {
		return err
	}

	req.Header.Set("BRI-Timestamp", timestamp)
	signature, err := utils.BriGenerateSignatureForDelete(path, method, payload, token, timestamp)
	if err != nil {
		return errors.New("failed to generate signature: " + err.Error())
	}
	req.Header.Set("BRI-Signature", signature)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var responseData entity.BriResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return  err
	}

	if !responseData.Status {
		return errors.New(responseData.ErrDesc)
	}

	fmt.Println(responseData)

	return s.repo.DeleteVaTransaction(ctx, vaUuid, false)
}
