package utils

import "time"

func GetCurrentDateTime() time.Time {
	return time.Now()
}

func ConvertTimeToDBDateTime(toChange time.Time) string {
	return toChange.Format(time.DateTime)
}
