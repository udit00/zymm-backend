package utils

import "time"

func GetCurrentDateTime() time.Time {
	return time.Now()
}

func ConvertTimeToDBDateTime(toChange time.Time) string {
	return toChange.UTC().Format("2006-01-02 15:04:05")
}
