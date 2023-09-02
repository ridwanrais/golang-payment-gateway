package utils

import (
	"crypto"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"

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

func MandiriGenerateAccessTokenSignature(timestamp string) (string, error) {
	cwd, _ := os.Getwd()
	privateKeyFile := filepath.Join(cwd, "credentials", "decrypted_Mandiri_API_Portal.pem")

	// password := os.Getenv("MANDIRI_API_PASSWORD")
	data := fmt.Sprintf("%s|%s", os.Getenv("MANDIRI_CLIENT_ID"), timestamp)

	// Load private key
	privateKeyBytes, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return "", errors.New("failed to decode PEM block containing the key")
	}

	decryptedBlock := block.Bytes

	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(decryptedBlock)
	if err != nil {
		return "", err
	}

	// Type assert to *rsa.PrivateKey
	rsaPrivateKey, ok := privateKeyInterface.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("failed to type assert privateKey to *rsa.PrivateKey")
	}

	// Generate Signature
	dataHash := sha256.New()
	_, err = dataHash.Write([]byte(data))
	if err != nil {
		return "", err
	}

	msgHashSum := dataHash.Sum(nil)
	signatureByte, err := rsa.SignPKCS1v15(nil, rsaPrivateKey, crypto.SHA256, msgHashSum)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(signatureByte)
	return signature, nil
}

func MandiriGenerateTransactionSignature(path, method, accessToken string, mandiriVaData *entity.MandiriVaData, timestamp string) (string, error) {
	clientSecret := os.Getenv("MANDIRI_CLIENT_SECRET")

	jsonBodyRequest := ""
	if mandiriVaData != nil {
		// Convert brivaData to JSON string for body
		bodyBytes, err := json.Marshal(*mandiriVaData)
		if err != nil {
			return "", errors.New("error marshalling request body: " + err.Error())
		}
		jsonBodyRequest = string(bodyBytes)
	}

	// // Minify json request body
	// var jsonData map[string]interface{}
	// if err := json.Unmarshal([]byte(jsonBodyRequest), &jsonData); err != nil {
	// 	return "", err
	// }
	// minifiedJSON, err := json.Marshal(jsonData)
	// if err != nil {
	// 	return "", err
	// }

	// Generate the lowercase hex encoded SHA256 hash of minified request body
	hasher := sha256.New()
	hasher.Write([]byte(jsonBodyRequest))
	hashedBody := hex.EncodeToString(hasher.Sum(nil))

	// Construct the data string
	data := fmt.Sprintf("%s:%s:%s:%s:%s", method, path, accessToken, hashedBody, timestamp)

	// Generate the HMAC with SHA512
	h := hmac.New(sha512.New, []byte(clientSecret))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
