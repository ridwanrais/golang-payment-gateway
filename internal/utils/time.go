package utils

import "time"

func GenerateCurrentTimestamp() string {
	currentTime := time.Now().UTC()
	return currentTime.Format("2006-01-02T15:04:05.000Z")
}

func IsLaterThanNow(t time.Time) bool {
	now := time.Now().UTC()

	return t.After(now)
}

func MandiriGenerateCurrentTimestamp() string {
	location, _ := time.LoadLocation("Asia/Jakarta") // For UTC+7
	t := time.Now().In(location)

	return t.Format("2006-01-02T15:04:05.000T-0700")
}
