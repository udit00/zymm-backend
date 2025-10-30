package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	bussinessAuth "zymm/internal/business/auth"
	"zymm/internal/models"
	authRepo "zymm/internal/repository/auth_repo"
	employeesRepo "zymm/internal/repository/employees_repo"
	gymRepo "zymm/internal/repository/gym_repo"
	membershipRepo "zymm/internal/repository/membership_repo"
	notificationRepo "zymm/internal/repository/notification_repo"
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
	mux.HandleFunc(userRouteAppended("uploadProfilePicture"), AuthMiddleware(uploadProfilePicture))

	// Serve static images
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
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

	userData, err := userRepo.GetUserDetailsForSelfByUserId(currentUserId)
	if err != nil {
		LogService.LogError("❌ Error fetching user data: ", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching user data: "+err.Error())
		return
	}

	if !userData.IsActive {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Forcing Log out")
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

	// Get unread notification count
	unreadCount, unreadErr := notificationRepo.GetUnreadNotificationCount(currentUserId)
	if unreadErr != nil {
		LogService.LogError("❌ Error fetching unread notification count: ", unreadErr)
		// Don't fail the request, just set count to 0
		unreadCount = 0
	}

	selfDataResponse := models.SelfDataResponse{
		UserId:                  userData.UserId,
		UserName:                userData.UserName,
		Email:                   userData.Email,
		Mobile:                  userData.Mobile,
		Gender:                  userData.Gender,
		RoleId:                  userData.RoleId,
		ProfilePic:              userData.ProfilePic,
		VisitedToday:            false,
		UnreadNotificationCount: unreadCount,
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

func uploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("uploadProfilePicture was called")

	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get authenticated user from context
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

	// Parse multipart form data (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Unable to parse form: "+err.Error())
		return
	}

	// Get the file from the request
	file, handler, err := r.FormFile("profilePicture")
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Error retrieving file: "+err.Error())
		return
	}
	defer file.Close()

	// Log the content type for debugging
	contentType := handler.Header.Get("Content-Type")
	LogService.LogMessage(fmt.Sprintf("📸 Received file: %s, Content-Type: %s, Size: %d bytes", handler.Filename, contentType, handler.Size))

	// Validate file type - check both content type and extension
	fileExt := strings.ToLower(filepath.Ext(handler.Filename))
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	isValidContentType := strings.HasPrefix(contentType, "image/")
	isValidExtension := validExtensions[fileExt]

	if !isValidContentType && !isValidExtension {
		utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Only image files are allowed. Received: %s (ext: %s)", contentType, fileExt))
		return
	}

	// Validate file size (max 10MB)
	if handler.Size > 10<<20 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "File size exceeds 10MB limit")
		return
	}

	// Create images directory if it doesn't exist
	imagesDir := "./images/profile_pictures"
	if err := os.MkdirAll(imagesDir, os.ModePerm); err != nil {
		LogService.LogError("❌ Error creating images directory: ", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error creating upload directory")
		return
	}

	// Use the already validated fileExt, or default to .jpg if empty
	if fileExt == "" {
		// Try to get extension from content type if filename had no extension
		switch contentType {
		case "image/jpeg", "image/jpg":
			fileExt = ".jpg"
		case "image/png":
			fileExt = ".png"
		case "image/gif":
			fileExt = ".gif"
		case "image/webp":
			fileExt = ".webp"
		default:
			fileExt = ".jpg"
		}
	}

	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("user_%d_%d%s", currentUserId, timestamp, fileExt)
	filepath := filepath.Join(imagesDir, filename)

	// Create the file
	dst, err := os.Create(filepath)
	if err != nil {
		LogService.LogError("❌ Error creating file: ", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error saving file")
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	bytesWritten, err := io.Copy(dst, file)
	if err != nil {
		LogService.LogError("❌ Error copying file: ", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error saving file")
		return
	}

	LogService.LogMessage(fmt.Sprintf("✅ Image saved: %s (%.2f MB, %d bytes)", filename, float64(bytesWritten)/1024/1024, bytesWritten))

	// Generate the URL for the image
	// Get the host from the request
	host := r.Host
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	imageUrl := fmt.Sprintf("%s://%s/images/profile_pictures/%s", scheme, host, filename)

	// Update user's profile picture in database
	err = userRepo.UpdateUserProfilePicture(currentUserId, imageUrl)
	if err != nil {
		// If database update fails, delete the uploaded file
		os.Remove(filepath)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error updating profile picture: "+err.Error())
		return
	}

	// Return success response with the image URL
	response := map[string]interface{}{
		"message":    "Profile picture uploaded successfully",
		"profilePic": imageUrl,
	}

	utils.SendSuccessResponse(w, http.StatusOK, response)
}
