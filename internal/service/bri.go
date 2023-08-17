package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

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
	signature, err := utils.BriGenerateSignature(path, method, brivaData, token, timestamp)
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

	// Generate a random UUID for reference_number
	refNumber := utils.GenerateReferenceNumber()

	// Construct virtual account number
	virtualAccountNumber := brivaData.BrivaNo + brivaData.CustCode

	_, err = s.repo.InsertBrivaTransaction(ctx, brivaData, refNumber, virtualAccountNumber)
	if err != nil {
		return nil, err
	}

	return &entity.CreateBrivaResponse{
		ReferenceNumber:      refNumber,
		VirtualAccountNumber: brivaData.BrivaNo + brivaData.CustCode,
	}, nil
}
