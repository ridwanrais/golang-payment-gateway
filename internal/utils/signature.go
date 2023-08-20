package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ridwanrais/golang-payment-gateway/internal/entity"
)

func BriGenerateSignature(path, method string, brivaData *entity.BrivaData, token, timestamp string) (string, error) {
	body := ""
	if brivaData != nil {
		// Convert brivaData to JSON string for body
		bodyBytes, err := json.Marshal(*brivaData)
		if err != nil {
			return "", errors.New("error marshalling request body: " + err.Error())
		}
		body = string(bodyBytes)
	}

	// Construct the payload string
	payload := fmt.Sprintf("path=%s&verb=%s&token=Bearer %s&timestamp=%s&body=%s", path, method, token, timestamp, body)

	// Encrypt the payload string using the SHA256-HMAC algorithm with the Consumer Secret as the key
	h := hmac.New(sha256.New, []byte(os.Getenv("BRI_CLIENT_SECRET")))
	h.Write([]byte(payload))

	// Encode the result with Base64
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature, nil
}

func BriGenerateSignatureForDelete(path, method string, body string, token, timestamp string) (string, error) {
	// Construct the payload string
	payload := fmt.Sprintf("path=%s&verb=%s&token=Bearer %s&timestamp=%s&body=%s", path, method, token, timestamp, body)

	// Encrypt the payload string using the SHA256-HMAC algorithm with the Consumer Secret as the key
	h := hmac.New(sha256.New, []byte(os.Getenv("BRI_CLIENT_SECRET")))
	h.Write([]byte(payload))

	// Encode the result with Base64
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature, nil
}

