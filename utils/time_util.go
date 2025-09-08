package utils

import "time"

func GetCurrentDateTime() time.Time {
	return time.Now()
}

func GetDbDateTime() string {
	return GetCurrentDateTime().UTC().Format(time.RFC3339)
}
