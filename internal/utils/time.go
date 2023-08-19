package utils

import "time"

func GenerateCurrentTimestamp() string {
	currentTime := time.Now().UTC()
	return currentTime.Format("2006-01-02T15:04:05.000Z")
}

func IsLaterThanToday(t time.Time) bool {
	today := time.Now().UTC()

	// Compare only the date components
	return t.After(time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC))
}
