package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	bussinessAuth "zymm/internal/business/auth"
	"zymm/internal/config"
	"zymm/internal/db"
	"zymm/internal/models"
	gymModels "zymm/internal/models/gym_models"
	authRepo "zymm/internal/repository/auth_repo"
	LogService "zymm/internal/service/log_service"
	"zymm/utils"
)

const authApiVersion = "v1"
const authApiPrefix = "auth"

func authRouteAppended(newRoute string) string {
	return utils.ApiRoute(authApiVersion, authApiPrefix, newRoute)
}

func AuthHandlerDelegate(mux *http.ServeMux) {
	mux.HandleFunc(authRouteAppended("getAuthTokenData"), getAuthTokenData)
	mux.HandleFunc(authRouteAppended("login"), loginHandler)
	mux.HandleFunc(authRouteAppended("registration"), registerHandler)
	mux.HandleFunc(authRouteAppended("ownerRegistration"), ownerRegistrationHandler)
}

func getAuthTokenData(w http.ResponseWriter, r *http.Request) {
	if config.IsDebug() {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid Authorization header")
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		claims, err := bussinessAuth.GetDataFromJWTToken(token)
		if err != nil {
			LogService.LogError("AuthMiddleware: token validation failed: ", err)
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid or expired auth token")
			return
		}
		utils.SendSuccessResponse(w, http.StatusOK, claims)
	} else {
		utils.SendErrorResponse(w, http.StatusForbidden, "What are you doing here?")
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("loginHandler was called")
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Decode JSON body
	var req models.LoginRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Check user in DB
	userDataModel, err := authRepo.GetUserDataByEmailOrMobile(req.EmailOrMobile)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials: no user found with email or mobile "+req.EmailOrMobile)
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error retrieving user: "+err.Error())
		return
	}

	isPassSame, passMatchError := bussinessAuth.ComparePasswordArgon2id(userDataModel.Password, req.Password)
	if passMatchError != nil || !isPassSame {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials: password was wrong")
		return
	}

	loginLogsModel := models.LoginLogsRecord{
		UserId:       userDataModel.UserId,
		AppVersion:   req.AppVersion,
		UserAgent:    req.UserAgent,
		LocationLat:  &req.LocationLat,
		LocationLong: &req.LocationLong,
		IpAddress:    &req.IpAddress,
	}

	authRepo.InsertLoginLog(loginLogsModel)
	generatedJwt, jwtError := bussinessAuth.GenerateJWTToken(loginLogsModel.UserId)
	if jwtError != nil || generatedJwt == nil || *generatedJwt == "" {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error generating JWT token: "+jwtError.Error())
		return
	}

	updateTokenError := authRepo.UpdateLoginAuthToken(loginLogsModel.UserId, *generatedJwt)
	if updateTokenError != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error updating auth token: "+updateTokenError.Error())
		return
	}
	data := models.LoginResponseModel{
		DisplayName:  userDataModel.DisplayName,
		AuthCheckSum: *generatedJwt,
	}
	utils.SendSuccessResponse(w, http.StatusOK, data)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Decode JSON body
	var req models.RegistrationApiRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	validationRequestErr := bussinessAuth.ValidateRegistrationRequest(req)
	if validationRequestErr != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, validationRequestErr.Error())
		return
	}

	insertedUserRecord, userInsertionError := authRepo.InsertRegistrationRecord(models.UserRecord{
		DisplayPic: req.DisplayPic,
		UserName:   req.DisplayName,
		Mobile:     req.Mobile,
		Email:      req.Email,
		UserPass:   req.Password,
		UserId:     0,
		Gender:     req.Gender,
		ProfilePic: req.DisplayPic,
		RoleId:     5, // Default role ID for regular users
	})

	if userInsertionError != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error creating user: "+userInsertionError.Error())
		return
	}

	var regLog models.RegistrationLogsRecord = models.RegistrationLogsRecord{
		UserId:     insertedUserRecord.UserId,
		UserName:   insertedUserRecord.UserName,
		UserPass:   insertedUserRecord.UserPass,
		Gender:     insertedUserRecord.Gender,
		Mobile:     insertedUserRecord.Mobile,
		Email:      insertedUserRecord.Email,
		ProfilePic: insertedUserRecord.ProfilePic,
		RoleId:     insertedUserRecord.RoleId,
		AppVersion: req.AppVersion,
		UserAgent:  req.UserAgent,
	}

	registrationLogInsertionError := authRepo.InsertRegistrationLog(regLog)

	if registrationLogInsertionError != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error logging registration: "+registrationLogInsertionError.Error())
		return
	}

	loginLogsModel := models.LoginLogsRecord{
		UserId:       insertedUserRecord.UserId,
		AppVersion:   req.AppVersion,
		UserAgent:    req.UserAgent,
		LocationLat:  &req.LocationLat,
		LocationLong: &req.LocationLong,
		IpAddress:    &req.IpAddress,
	}

	authRepo.InsertLoginLog(loginLogsModel)
	generatedJwt, jwtError := bussinessAuth.GenerateJWTToken(loginLogsModel.UserId)
	if jwtError != nil || generatedJwt == nil || *generatedJwt == "" {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error generating JWT token: "+jwtError.Error())
		return
	}

	resp := models.RegistrationApiResponseModel{
		DisplayName:  insertedUserRecord.UserName,
		AuthCheckSum: *generatedJwt,
	}

	utils.SendSuccessResponse(w, http.StatusOK, resp)
}

func ownerRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("ownerRegistrationHandler was called")
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	// Decode JSON body
	var req models.RegistrationOwnerApiRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	validationRequestErr := bussinessAuth.ValidateOwnerRegistrationRequest(req)
	if validationRequestErr != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, validationRequestErr.Error())
		return
	}

	insertedUserRecord, userInsertionError := authRepo.InsertRegistrationRecord(models.UserRecord{
		DisplayPic: req.DisplayPic,
		UserName:   req.DisplayName,
		Mobile:     req.Mobile,
		Email:      req.OwnerPersonalEmail,
		UserPass:   req.Password,
		UserId:     0,
		Gender:     req.Gender,
		ProfilePic: req.DisplayPic,
		RoleId:     1,
	})

	if userInsertionError != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error creating user: "+userInsertionError.Error())
		return
	}

	var regLog models.RegistrationLogsRecord = models.RegistrationLogsRecord{
		UserId:     insertedUserRecord.UserId,
		UserName:   insertedUserRecord.UserName,
		UserPass:   insertedUserRecord.UserPass,
		Gender:     insertedUserRecord.Gender,
		Mobile:     insertedUserRecord.Mobile,
		Email:      insertedUserRecord.Email,
		ProfilePic: insertedUserRecord.ProfilePic,
		RoleId:     insertedUserRecord.RoleId,
		AppVersion: req.AppVersion,
		UserAgent:  req.UserAgent,
	}

	registrationLogInsertionError := authRepo.InsertRegistrationLog(regLog)

	if registrationLogInsertionError != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error logging registration: "+registrationLogInsertionError.Error())
		return
	}

	gymRecord := gymModels.GymRecord{
		GymName:       req.GymName,
		State:         req.State,
		City:          req.City,
		GymAddress:    req.GymAddress,
		ContactNo:     req.ContactNo,
		OfficialEmail: req.OfficialEmail,
		CreatedBy:     insertedUserRecord.UserId,
		LocationLat:   req.LocationLat,
		LocationLong:  req.LocationLong,
	}

	_, gymInsertionError := authRepo.InsertGym(db.DB, gymRecord)

	if gymInsertionError != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error creating gym: "+gymInsertionError.Error())
		return
	}

	loginLogsModel := models.LoginLogsRecord{
		UserId:       insertedUserRecord.UserId,
		AppVersion:   req.AppVersion,
		UserAgent:    req.UserAgent,
		LocationLat:  &req.LocationLat,
		LocationLong: &req.LocationLong,
		IpAddress:    &req.IpAddress,
	}

	authRepo.InsertLoginLog(loginLogsModel)
	generatedJwt, jwtError := bussinessAuth.GenerateJWTToken(loginLogsModel.UserId)
	if jwtError != nil || generatedJwt == nil || *generatedJwt == "" {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error generating JWT token: "+jwtError.Error())
		return
	}

	updateAuthTokenErr := authRepo.UpdateLoginAuthToken(insertedUserRecord.UserId, *generatedJwt)
	if updateAuthTokenErr != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, updateAuthTokenErr.Error())
		return
	}

	resp := models.RegistrationApiResponseModel{
		DisplayName:  insertedUserRecord.UserName,
		AuthCheckSum: *generatedJwt,
	}

	utils.SendSuccessResponse(w, http.StatusOK, resp)
}
