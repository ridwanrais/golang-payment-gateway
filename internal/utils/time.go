package utils

import "time"

func GenerateCurrentTimestamp() string {
	currentTime := time.Now().UTC()
	return currentTime.Format("2006-01-02T15:04:05.000Z")
}

func IsLaterThanNow(t time.Time) bool {
	now := time.Now().UTC()

	// Compare only the date components
	return t.After(now)
}
