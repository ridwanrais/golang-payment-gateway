package utils

import (
	"fmt"
	"math/rand"
	"time"
)

var rng *rand.Rand

func init() {
	s := rand.NewSource(time.Now().UnixNano())
	rng = rand.New(s)
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}

func GenerateReferenceNumber() string {
	currentDate := time.Now().Format("20060102") // YYYYMMDD
	randomChars := GenerateRandomString(4)
	return fmt.Sprintf("%s-%s", currentDate, randomChars)
}
