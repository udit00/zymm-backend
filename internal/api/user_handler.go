package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	bussinessAuth "zymm/internal/business/auth"
	"zymm/internal/models"
	authRepo "zymm/internal/repository/auth_repo"
	employeesRepo "zymm/internal/repository/employees_repo"
	gymRepo "zymm/internal/repository/gym_repo"
	membershipRepo "zymm/internal/repository/membership_repo"
	"zymm/internal/repository/userRepo"
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
	mux.HandleFunc(userRouteAppended("changePassword"), AuthMiddleware(changePassword))
	mux.HandleFunc(userRouteAppended("deleteProfile"), AuthMiddleware(deleteProfile))
}

func getSelfData(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getSelfData was called")

	claims, ok := r.Context().Value(ctxClaimDataKey).(*bussinessAuth.MyCustomClaims)
	if !ok || claims == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	currentUserId := claims.UserId
	if currentUserId < 0 {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "invalid user id in context")
		return
	}

	var planDetails *models.PlanRecord
	var planErr error

	userData, err := userRepo.GetUserByUserId(currentUserId)
	if err != nil {
		LogService.LogError("❌ Error fetching user data: ", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching user data: "+err.Error())
		return
	}

	userActiveMembership, err := membershipRepo.GetUserMembershipByUserId(currentUserId)
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
		selfDataResponse.PlanDetails = planDetails
	}
	if userActiveMembership != nil {
		selfDataResponse.MembershipId = &userActiveMembership.MembershipId
		selfDataResponse.ActiveMembershipDetails = userActiveMembership
	}
	if userData.RoleId < 5 {
		if userData.RoleId == 1 {
			gymDetails, gymDetailsErr := gymRepo.GetGymWithAdditionalDataByOwnerId(userData.UserId)
			if gymDetailsErr != nil {
				utils.SendErrorResponse(w, http.StatusInternalServerError, gymDetailsErr.Error())
				return
			}
			selfDataResponse.GymId = &gymDetails.GymId
		} else {
			employeeDetails, employeeDetailsErr := employeesRepo.GetEmployeeByUserId(userData.UserId)
			if employeeDetailsErr != nil {
				utils.SendErrorResponse(w, http.StatusInternalServerError, employeeDetailsErr.Error())
				return
			}
			selfDataResponse.GymId = &employeeDetails.GymId
		}
	}
	utils.SendSuccessResponse(w, http.StatusOK, selfDataResponse)
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("changePassword was called")

	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.ChangePasswordRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	claims, ok := r.Context().Value(ctxClaimDataKey).(*bussinessAuth.MyCustomClaims)
	if !ok || claims == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	currentUserId := claims.UserId
	if currentUserId < 0 {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "invalid user id in context")
		return
	}

	userDetails, userDetailsErr := userRepo.GetUserByUserId(currentUserId)
	if userDetailsErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, userDetailsErr.Error())
		return
	}

	loginUserData, loginUserDataErr := userRepo.GetUserDataByEmailOrMobile(userDetails.Mobile)
	if loginUserDataErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, loginUserDataErr.Error())
		return
	}

	isPassSame, passMatchError := bussinessAuth.ComparePasswordArgon2id(loginUserData.Password, req.OldPassword)
	if passMatchError != nil || !isPassSame {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials: Old password was wrong")
		return
	}

	if req.OldPassword == req.NewPassword {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Password cannot be same.")
		return
	}

	if len(req.NewPassword) < 6 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Password must be at least 6 characters.")
		return
	}

	changePasswordErr := authRepo.UpdateUserPassword(claims.UserId, req.NewPassword)
	if changePasswordErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, changePasswordErr.Error())
		return
	}

	generatedJwt, jwtError := bussinessAuth.GenerateJWTToken(claims.UserId, claims.RoleId)
	if jwtError != nil || generatedJwt == nil || *generatedJwt == "" {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Error generating JWT token: "+jwtError.Error())
		return
	}

	updateTokenError := authRepo.UpdateLoginAuthToken(claims.UserId, *generatedJwt)
	if updateTokenError != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Error updating auth token: "+updateTokenError.Error())
		return
	}

	responseModel := &models.ChangePasswordResponseModel{
		AuthToken: *generatedJwt,
	}

	utils.SendSuccessResponse(w, http.StatusAccepted, responseModel)
}

func deleteProfile(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("deleteProfile was called")

	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.DeleteProfileRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Password == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	claims, ok := r.Context().Value(ctxClaimDataKey).(*bussinessAuth.MyCustomClaims)
	if !ok || claims == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	currentUserId := claims.UserId
	if currentUserId < 0 {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "invalid user id in context")
		return
	}

	userDetails, userDetailsErr := userRepo.GetUserByUserId(currentUserId)
	if userDetailsErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, userDetailsErr.Error())
		return
	}

	loginUserData, loginUserDataErr := userRepo.GetUserDataByEmailOrMobile(userDetails.Mobile)
	if loginUserDataErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, loginUserDataErr.Error())
		return
	}

	isPassSame, passMatchError := bussinessAuth.ComparePasswordArgon2id(loginUserData.Password, req.Password)
	if passMatchError != nil || !isPassSame {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials: Old password was wrong")
		return
	}

	deleteProfileError := authRepo.InActiveUser(claims.UserId)
	if deleteProfileError != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Error delete your account: "+deleteProfileError.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusAccepted, nil)
}
