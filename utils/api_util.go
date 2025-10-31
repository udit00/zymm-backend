package utils

import (
	"encoding/json"
	"net/http"
	"zymm/internal/config"
	"zymm/internal/models"
)

func ApiRoute(apiVersion string, apiPrefix string, newRoute string) string {
	return "/" + config.GetAppName() + "/" + apiVersion + "/" + apiPrefix + "/" + newRoute
}

func SendErrorResponse(writer http.ResponseWriter, status int, errMsg string, data ...interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	var finalApiResponse models.APIResponse
	if data != nil {
		finalApiResponse = models.APIResponse{
			Status: status,
			Error:  errMsg,
			Data:   data,
		}
	} else {
		finalApiResponse = models.APIResponse{
			Status: status,
			Error:  errMsg,
		}
	}

	json.NewEncoder(writer).Encode(finalApiResponse)
}

func SendSuccessResponse(writer http.ResponseWriter, status int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	finalApiResponse := models.APIResponse{
		Status: status,
		Data:   data,
	}
	json.NewEncoder(writer).Encode(finalApiResponse)
}
