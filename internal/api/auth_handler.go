package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	bussinessAuth "zymm/internal/business/auth"
	"zymm/internal/db"
	"zymm/internal/models"
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
	mux.HandleFunc(authRouteAppended("login"), loginHandler)
	mux.HandleFunc(authRouteAppended("registration"), registerHandler)
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
	var dbPassword string
	err := db.DB.QueryRow("SELECT userPass FROM users WHERE email = @p1 or mobile = @p2", req.EmailOrMobile, req.EmailOrMobile).Scan(&dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials: no user found with email or mobile "+req.EmailOrMobile)
			return
		}
		LogService.LogError("❌ DB error: ", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	isPassSame, passMatchError := bussinessAuth.ComparePasswordArgon2id(dbPassword, req.Password)
	if passMatchError != nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials: pass: "+req.Password+" was wrong, correct Pass is "+dbPassword)
		return
	} else if !isPassSame {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials: pass: "+req.Password+" was wrong, correct Pass is "+dbPassword)
		return
	}
	// Success
	var data = map[string]string{"message": "✅ Login successful for " + req.EmailOrMobile}
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

	// Basic validation
	if req.DisplayName == "" || req.Mobile == "" || req.Password == "" || req.Gender == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Missing required fields")
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
		UserId:      insertedUserRecord.UserId,
		UserName:    insertedUserRecord.UserName,
		UserPass:    insertedUserRecord.UserPass,
		Gender:      insertedUserRecord.Gender,
		Mobile:      insertedUserRecord.Mobile,
		Email:       insertedUserRecord.Email,
		ProfilePic:  insertedUserRecord.ProfilePic,
		RoleId:      insertedUserRecord.RoleId,
		AppVersion:  req.AppVersion,
		AppPlatform: req.AppPlatform,
	}

	registrationLogInsertionError := authRepo.InsertRegistrationLog(regLog)

	if registrationLogInsertionError != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error logging registration: "+registrationLogInsertionError.Error())
		return
	}

	// Insert user into DB and get new userId

	resp := models.RegistrationApiResponseModel{
		UserId:     insertedUserRecord.UserId,
		UserName:   insertedUserRecord.UserName,
		Registered: true,
		Message:    "✅ Registration successful for " + insertedUserRecord.UserName,
	}

	utils.SendSuccessResponse(w, http.StatusOK, resp)
}
