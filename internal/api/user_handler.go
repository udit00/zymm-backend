package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"zymm/internal/config"
	"zymm/internal/models"
	"zymm/utils"
)

const userApiVersion = "v1"
const userApiPrefix = "user"

func userRouteAppended(newRoute string) string {
	return "/" + config.AppName + "/" + userApiVersion + "/" + userApiPrefix + "/" + newRoute
}

func UserHandlerDelegate(mux *http.ServeMux) {
	mux.HandleFunc(userRouteAppended("getUserData"), getUserData)
}

func getDummyDataForUser() models.LoginResponseModel {
	return models.LoginResponseModel{
		AuthCheckSum: "egoisaedoignaegt",
		DisplayName:  "HaXeR-.",
		LastLoggedIn: utils.GetDbDateTime(),
	}
}

func getUserData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getUserData was called")
	// Pretend login is successful and build response
	resp := getDummyDataForUser()

	apiResp := models.APIResponse{
		Status: 1,
		Data:   resp,
	}

	json.NewEncoder(w).Encode(apiResp)
}
