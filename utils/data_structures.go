package utils

import "strconv"

func ConvertStringToInt(data string) *int {
	val, err := strconv.Atoi(data)
	if err != nil {
		return nil
	}
	return &val
}

func ConvertIntToString(val int) string {
	return strconv.Itoa(val)
}
