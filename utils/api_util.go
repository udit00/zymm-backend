package utils

import (
	"zymm/internal/config"
)

func ApiRoute(apiVersion string, apiPrefix string, newRoute string) string {
	return "/" + config.AppName + "/" + apiVersion + "/" + apiPrefix + "/" + newRoute
}
