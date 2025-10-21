package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"zymm/internal/models"
	authRepo "zymm/internal/repository/auth_repo"
	membershipRepo "zymm/internal/repository/membership_repo"
	LogService "zymm/internal/service/log_service"
	"zymm/utils"
)

const userApiVersion = "v1"
const userApiPrefix = "user"

func userRouteAppended(newRoute string) string {
	return utils.ApiRoute(userApiVersion, userApiPrefix, newRoute)
}

func UserHandlerDelegate(mux *http.ServeMux) {
	mux.HandleFunc(userRouteAppended("selfData"), AuthMiddleware(getSelfData))
	mux.HandleFunc(userRouteAppended("getUserContacts"), AuthMiddleware(getUserContacts))
}

func getSelfData(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getSelfData was called")
	// Extract user id from context (set by AuthMiddleware)
	uid := r.Context().Value(ctxUserIDKey)
	if uid == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userId, ok := uid.(int)
	if !ok {
		// safety: try converting from float64 (jwt numeric) or other numeric types
		switch v := uid.(type) {
		case int64:
			userId = int(v)
		case float64:
			userId = int(v)
		default:
			utils.SendErrorResponse(w, http.StatusUnauthorized, "invalid user id in context")
			return
		}
	}

	var planDetails *models.PlanRecord
	var planErr error

	userData, err := authRepo.GetUserByUserId(userId)
	if err != nil {
		LogService.LogError("❌ Error fetching user data: ", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching user data: "+err.Error())
		return
	}

	userActiveMembership, err := membershipRepo.GetUserMembershipByUserId(userId)
	if err != nil {
		if err != sql.ErrNoRows {
			LogService.LogError("❌ Error fetching user active membership data: ", err)
			utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching user active membership data: "+err.Error())
			return
		}
	}

	if userActiveMembership != nil {
		planDetails, planErr = membershipRepo.GetPlanById(userActiveMembership.PlanId)
		if planErr != nil {
			if planErr != sql.ErrNoRows {
				LogService.LogError("❌ Error fetching plan details: ", planErr)
				utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching plan details: "+planErr.Error())
				return
			}
		}
	}

	selfDataResponse := models.SelfDataResponse{
		UserId:       userData.UserId,
		UserName:     userData.UserName,
		Email:        userData.Email,
		Mobile:       userData.Mobile,
		Gender:       userData.Gender,
		RoleId:       userData.RoleId,
		ProfilePic:   userData.ProfilePic,
		VisitedToday: false,
	}
	if planDetails != nil {
		selfDataResponse.PlanId = &planDetails.PlanId
		selfDataResponse.GymId = planDetails.GymId
	}
	if userActiveMembership != nil {
		selfDataResponse.MembershipId = &userActiveMembership.MembershipId
		selfDataResponse.ActiveMembershipDetails = userActiveMembership
	}
	utils.SendSuccessResponse(w, http.StatusOK, selfDataResponse)
}

func getUserContacts(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getUserContacts was called")
	// Extract user id from context (set by AuthMiddleware)
	uid := r.Context().Value(ctxUserIDKey)
	if uid == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userId, ok := uid.(int)
	if !ok {
		// safety: try converting from float64 (jwt numeric) or other numeric types
		switch v := uid.(type) {
		case int64:
			userId = int(v)
		case float64:
			userId = int(v)
		default:
			utils.SendErrorResponse(w, http.StatusUnauthorized, "invalid user id in context")
			return
		}
	}

	utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("aogaog with %v", userId))

	// contacts, err := authRepo.GetContactsByUserId(userId)
	// if err != nil {
	// 	LogService.LogError("❌ Error fetching user contacts: ", err)
	// 	utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching user contacts: "+err.Error())
	// 	return
	// }

	// utils.SendSuccessResponse(w, http.StatusOK, contacts)
}
