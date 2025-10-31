package LogService

import "fmt"

func LogMessage(msg string) {
	fmt.Println("[LOG]: ", msg)
}

func LogError(msg string, err error) {
	if err != nil {
		fmt.Println("[ERROR]: ", msg, err.Error())
	}
}
