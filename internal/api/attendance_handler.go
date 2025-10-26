package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	bussinessAuth "zymm/internal/business/auth"
	businessRoles "zymm/internal/business/roles"
	businessRoleType "zymm/internal/business/roles/roles_type"
	"zymm/internal/models"
	attendanceRepo "zymm/internal/repository/attendance_repo"
	LogService "zymm/internal/service/log_service"
	"zymm/utils"
)

const attendanceApiVersion = "v1"
const attendanceApiPrefix = "attendance"

func attendanceRouteAppended(newRoute string) string {
	return utils.ApiRoute(attendanceApiVersion, attendanceApiPrefix, newRoute)
}

func AttendanceApiPrefixHandlerDelegate(mux *http.ServeMux) {
	mux.HandleFunc(attendanceRouteAppended("punchIn"), AuthMiddleware(punchInAttendance))
	mux.HandleFunc(attendanceRouteAppended("punchOut"), AuthMiddleware(punchOutAttendance))
	mux.HandleFunc(attendanceRouteAppended("deleteAttendance"), AuthMiddleware(deleteAttendance))
	mux.HandleFunc(attendanceRouteAppended("getAllAttendance"), AuthMiddleware(getAllAttendanceForUser))
}

func punchInAttendance(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("punchInAttendance was called")
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.PunchInAttendanceRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.LocationLat == nil || req.LocationLong == nil || *req.LocationLat == "" || *req.LocationLong == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Request was not proper.")
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

	req.UserId = currentUserId

	lastInsertedRecord, lastInsertedRecordErr := attendanceRepo.LastAttendanceRecordOfTheUser(req.UserId)
	if lastInsertedRecordErr != nil && lastInsertedRecordErr != sql.ErrNoRows {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, lastInsertedRecordErr.Error())
		return
	}

	if lastInsertedRecord != nil && (lastInsertedRecord.PunchOutTime == nil || *lastInsertedRecord.PunchOutTime == "") {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Last punch out was not recorded, Please punch out first.")
		return
	}

	insertAttendanceId, insertAttendanceErr := attendanceRepo.InsertSubmitAttendance(req)
	if insertAttendanceErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, insertAttendanceErr.Error())
		return
	}

	if insertAttendanceId == nil || (insertAttendanceId != nil && *insertAttendanceId <= 0) {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Something went wrong while trying to insert attendance")
		return
	}

	utils.SendSuccessResponse(w, http.StatusAccepted, nil)

}

func punchOutAttendance(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("punchOutAttendance was called")
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.PunchOutAttendanceRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.LocationLat == nil || req.LocationLong == nil || *req.LocationLat == "" || *req.LocationLong == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Request was not proper.")
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

	req.UserId = currentUserId

	lastInsertedRecord, lastInsertedRecordErr := attendanceRepo.LastAttendanceRecordOfTheUser(req.UserId)
	if lastInsertedRecordErr != nil && lastInsertedRecordErr != sql.ErrNoRows {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, lastInsertedRecordErr.Error())
		return
	}

	if lastInsertedRecordErr == sql.ErrNoRows {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "No record found, Please punch in first.")
		return
	}

	if lastInsertedRecord != nil && (lastInsertedRecord.PunchOutTime == nil || *lastInsertedRecord.PunchOutTime == "") {
		updatePunchOutErr := attendanceRepo.UpdatePunchOutTime(req, lastInsertedRecord.AttendanceId)
		if updatePunchOutErr != nil {
			utils.SendErrorResponse(w, http.StatusExpectationFailed, updatePunchOutErr.Error())
			return
		}
		utils.SendSuccessResponse(w, http.StatusAccepted, nil)
	} else {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Punch out was already recorded for the last punch in, Please punch in again.")
		return
	}

}

func getAllAttendanceForUser(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getAllAttendanceForUser was called")
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	queryParams := r.URL.Query()
	userIdRawStr := queryParams.Get("userId")

	userIdConverted := utils.ConvertStringToInt(userIdRawStr)
	if userIdConverted == nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Request was not proper.")
		return
	}

	userId := *userIdConverted

	if userId <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Request was not proper.")
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

	if userId != currentUserId && businessRoles.IsNotAllowedToViewAttendance(businessRoleType.RoleType(claims.RoleId)) {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Access Denied.")
		return
	}

	allAttendanceData, allAttendanceDataErr := attendanceRepo.GetAllAttendanceByUserId(userId)
	if allAttendanceDataErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, allAttendanceDataErr.Error())
		return
	}
	if allAttendanceData == nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Data not found.")
		return
	}

	utils.SendSuccessResponse(w, http.StatusAccepted, allAttendanceData)

}

func deleteAttendance(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("deleteAttendance was called")
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.DeleteAttendanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.AttendanceId <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Request was not proper.")
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

	attendanceDetails, attendanceDetailsErr := attendanceRepo.GetAttendanceRecord(req.AttendanceId)
	if attendanceDetailsErr != nil {
		if attendanceDetailsErr == sql.ErrNoRows {
			utils.SendErrorResponse(w, http.StatusExpectationFailed, "No data found for Attendance ID.")
			return
		}
		utils.SendErrorResponse(w, http.StatusExpectationFailed, attendanceDetailsErr.Error())
		return
	}
	if attendanceDetails == nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "No data found for Attendance ID.")
		return
	}

	if currentUserId != attendanceDetails.UserId {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Access Denied")
		return
	}

	deleteAttendanceRecordErr := attendanceRepo.DeleteAttendanceRecord(attendanceDetails.AttendanceId)
	if deleteAttendanceRecordErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, deleteAttendanceRecordErr.Error())
		return
	}
	utils.SendSuccessResponse(w, http.StatusAccepted, nil)
}
