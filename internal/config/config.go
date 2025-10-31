package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	IsDebug bool
	AppName string
	Port    int
}

var appConfig AppConfig

func Init() {
	appConfig = AppConfig{
		IsDebug: false,
		AppName: "zymm",
		Port:    5000,
	}
	// load variables from .env into the environment
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, falling back to system env")
	}
	isDebugStr := os.Getenv("IS_DEBUG")
	debugParsedBool, debugParsingErr := strconv.ParseBool(isDebugStr)
	if debugParsingErr != nil {
		log.Fatal(" fetching IS_DEBUG caught with an exception. Exiting... -> with err " + debugParsingErr.Error())
	}
	portString := os.Getenv("ZYMM_PORT")
	if portString == "" {
		log.Fatal(" ZYMM_PORT not set in environment. Exiting...")
	}
	portParsedInt, portParsingErr := strconv.Atoi(portString)
	if portParsingErr != nil {
		log.Fatal(" PORT is wrong in environment. Exiting... with " + portParsingErr.Error())
	}
	if portParsedInt <= 0 && portParsedInt > 10000 {
		log.Fatal(" PORT is wrong in environment. Exiting..., port cannot be greater than 10000")
	}

	appConfig.IsDebug = debugParsedBool
	appConfig.Port = portParsedInt
}

func IsProduction() bool {
	return !appConfig.IsDebug
}

func IsDebug() bool {
	return appConfig.IsDebug
}

func GetAppPort() int {
	return appConfig.Port
}

func GetAppPortInString() string {
	return strconv.Itoa(appConfig.Port)
}

func GetAppName() string {
	return appConfig.AppName
}
