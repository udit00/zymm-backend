package api

import (
	"encoding/json"
	"net/http"
	"zymm/internal/config"
)

const userApiVersion = "v1"
const userApiPrefix = "user"

func userRouteAppended(newRoute string) string {
	return config.AppName + "/" + userApiVersion + "/" + userApiPrefix + "/" + newRoute
}

func UserHandlerDelegate(mux *http.ServeMux) {
	mux.HandleFunc(userRouteAppended("getUserData"), getUserData)
}

func getUserData(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "✅ Retrieve User Data Successful"})
}
