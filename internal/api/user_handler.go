package api

import (
	"net/http"
	authRepo "zymm/internal/repository/auth_repo"
	LogService "zymm/internal/service/log_service"
	"zymm/utils"
)

const userApiVersion = "v1"
const userApiPrefix = "user"

func userRouteAppended(newRoute string) string {
	return utils.ApiRoute(userApiVersion, userApiPrefix, newRoute)
}

func UserHandlerDelegate(mux *http.ServeMux) {
	mux.HandleFunc(userRouteAppended("getUserData"), getUserData)
}

func getUserData(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getUserData was called")
	// Pretend login is successful and build response
	userData, err := authRepo.GetUserByUserId(1)
	if err != nil {
		LogService.LogError("❌ Error fetching user data: ", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching user data: "+err.Error())
		return
	}
	// utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching user data: ", userData)

	utils.SendSuccessResponse(w, http.StatusOK, userData)
}
