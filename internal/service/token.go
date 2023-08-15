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
)

type Response struct {
	AccessToken string `json:"access_token"`
}

func (s *service) BriRetrieveAccessToken(ctx context.Context) (string, error) {
	clientID := os.Getenv("BRI_CLIENT_ID")
	clientSecret := os.Getenv("BRI_CLIENT_SECRET")

	key := fmt.Sprintf("bri_access_token:%s", clientID)
	fmt.Println(key)
	accessToken, err := s.repo.GetCacheValue(ctx, key)
	fmt.Println("accessToken", accessToken)
	fmt.Println("err", err)
	if err != nil {
		fmt.Println("MASUK SINI?")
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