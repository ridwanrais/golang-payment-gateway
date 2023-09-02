package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/ridwanrais/golang-payment-gateway/internal/entity"
	"github.com/ridwanrais/golang-payment-gateway/internal/utils"
)

type Response struct {
	AccessToken string `json:"access_token"`
}

func (s *service) BriRetrieveAccessToken(ctx context.Context) (string, error) {
	clientID := os.Getenv("BRI_CLIENT_ID")
	clientSecret := os.Getenv("BRI_CLIENT_SECRET")

	key := fmt.Sprintf("bri_access_token:%s", clientID)
	accessToken, err := s.repo.GetCacheValue(ctx, key)

	if err != nil {
		// Token not found or expired, fetch a new one from the API
		accessToken, err = s.briGetAccessTokenFromAPI(ctx, clientID, clientSecret)
		if err != nil {
			// Handle error while fetching a new access token
			fmt.Println("Error:", err)
			return "", err
		}
	}

	return accessToken, nil
}

func (s *service) briGetAccessTokenFromAPI(ctx context.Context, clientID, clientSecret string) (string, error) {
	briHost := os.Getenv("BRI_HOST")

	// Set the request URL
	url := fmt.Sprintf("%s/oauth/client_credential/accesstoken?grant_type=client_credentials", briHost)

	// Set the request headers
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	// Set the request data (client_id and client_secret)
	data := map[string]string{
		"client_id":     os.Getenv("BRI_CLIENT_ID"),
		"client_secret": os.Getenv("BRI_CLIENT_SECRET"),
	}

	// Create a new Resty client
	client := resty.New()

	// Set the request headers
	client.SetHeaders(headers)

	// Send the POST request and get the response
	resp, err := client.R().
		SetFormData(data).
		Post(url)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}

	var response Response
	err = json.Unmarshal([]byte(resp.String()), &response)
	if err != nil {
		fmt.Println("Error unmarshalling:", err)
		return "", err
	}

	newToken := response.AccessToken
	cacheExpiryTimeStr := os.Getenv("BRI_ACCESS_TOKEN_CACHE_EXPIRY")
	if cacheExpiryTimeStr == "" {
		return "", errors.New("BRI_ACCESS_TOKEN_CACHE_EXPIRY environment variable is not set")
	}

	cacheExpiryTime, err := strconv.Atoi(cacheExpiryTimeStr)
	if err != nil {
		fmt.Println("Error parsing BRI_ACCESS_TOKEN_CACHE_EXPIRY:", err)
		return "", errors.New("Invalid value for BRI_ACCESS_TOKEN_CACHE_EXPIRY. It should be a valid integer.")
	}

	expirationTime := time.Duration(cacheExpiryTime) * time.Hour

	err = s.repo.SetCacheValue(ctx, fmt.Sprintf("bri_access_token:%s", clientID), newToken, expirationTime)
	if err != nil {
		return "", err
	}
	return newToken, nil
}

func (s *service) MandiriRetrieveAccessToken(ctx context.Context) (string, error) {
	clientID := os.Getenv("MANDIRI_CLIENT_ID")

	key := fmt.Sprintf("mandiri_access_token:%s", clientID)
	accessToken, err := s.repo.GetCacheValue(ctx, key)

	// if err != nil {
	if true {
		// Token not found or expired, fetch a new one from the API
		accessToken, err = s.mandiriGetAccessTokenFromAPI(ctx)
		if err != nil {
			// Handle error while fetching a new access token
			fmt.Println("Error:", err)
			return "", err
		}
	}

	return accessToken, err
}

func (s *service) mandiriGetAccessTokenFromAPI(ctx context.Context) (string, error) {
	// mandiriHost := os.Getenv("MANDIRI_HOST")
	clientKey := os.Getenv("MANDIRI_CLIENT_ID")

	timestamp := utils.MandiriGenerateCurrentTimestamp()
	signature, err := utils.MandiriGenerateAccessTokenSignature(timestamp)
	if err != nil {
		// Handle error while fetching a new access token
		return "", err
	}

	fmt.Println("signature: ", signature)

	// Set the request URL
	// url := fmt.Sprintf("%s/openapi/auth/v2.0/access-token/b2b", mandiriHost)
	url := fmt.Sprintf("https://sandbox.bankmandiri.co.id/openapi/auth/v2.0/access-token/b2b")

	fmt.Println("timestamp", timestamp)

	// Set the request headers
	headers := map[string]string{
		"X-CLIENT-KEY": clientKey,
		"X-TIMESTAMP":  timestamp,
		"X-SIGNATURE":  signature,
		"Content-Type": "application/x-www-form-urlencoded",
	}

	// Set the request data
	data := map[string]string{
		"grantType": "client_credentials",
	}

	// Create a new Resty client
	client := resty.New()

	// Set the request headers
	client.SetHeaders(headers)

	// Send the POST request and get the response
	resp, err := client.R().
		SetFormData(data).
		Post(url)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}

	var response entity.MandiriAccessTokenResponse
	err = json.Unmarshal([]byte(resp.String()), &response)
	if err != nil {
		fmt.Println("Error unmarshalling:", err)
		return "", err
	}

	successResponse := len(response.ResponseMessage) >= 3 && response.ResponseMessage[:3] == "200"
	if !successResponse {
		return "", errors.New(response.ResponseMessage)
	}

	fmt.Println("response:", response)

	newToken := response.AccessToken
	cacheExpiryTimeStr := os.Getenv("MANDIRI_ACCESS_TOKEN_CACHE_EXPIRY_IN_MINUTES")
	if cacheExpiryTimeStr == "" {
		return "", errors.New("MANDIRI_ACCESS_TOKEN_CACHE_EXPIRY_IN_MINUTES environment variable is not set")
	}

	cacheExpiryTime, err := strconv.Atoi(cacheExpiryTimeStr)
	if err != nil {
		fmt.Println("Error parsing MANDIRI_ACCESS_TOKEN_CACHE_EXPIRY_IN_MINUTES:", err)
		return "", errors.New("Invalid value for MANDIRI_ACCESS_TOKEN_CACHE_EXPIRY_IN_MINUTES. It should be a valid integer.")
	}

	expirationTime := time.Duration(cacheExpiryTime) * time.Minute

	err = s.repo.SetCacheValue(ctx, fmt.Sprintf("mandiri_access_token:%s", clientKey), newToken, expirationTime)
	if err != nil {
		return "", err
	}
	return newToken, nil
}
