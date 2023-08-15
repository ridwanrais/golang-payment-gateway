package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/ridwanrais/golang-payment-gateway/internal/entity"
	"github.com/ridwanrais/golang-payment-gateway/internal/utils"
)

func (s *service) BriCreateBriva(ctx context.Context, RequestData entity.BrivaData) (*entity.BriResponseData, error) {
	timestamp := utils.GenerateCurrentTimestamp()

	token, err := s.BriRetrieveAccessToken(ctx)
	if err != nil {
		return nil, errors.New("failed to retrieve access token: " + err.Error())
	}

	path := "/v1/briva"
	method := "POST"
	signature, err := utils.BriGenerateSignature(path, method, RequestData, token, timestamp)
	if err != nil {
		return nil, errors.New("failed to generate signature: " + err.Error())
	}

	client := &http.Client{Timeout: 10 * time.Second}
	url := os.Getenv("BRI_HOST") + "/v1/briva"
	jsonData, err := json.Marshal(RequestData)
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

	return &responseData, nil
}
