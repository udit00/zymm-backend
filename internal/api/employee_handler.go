package api

import (
	"encoding/json"
	"net/http"
	bussinessAuth "zymm/internal/business/auth"
	businessRoles "zymm/internal/business/roles"
	businessRoleType "zymm/internal/business/roles/roles_type"
	"zymm/internal/models"
	authRepo "zymm/internal/repository/auth_repo"
	employeesRepo "zymm/internal/repository/employees_repo"
	gymRepo "zymm/internal/repository/gym_repo"
	"zymm/internal/repository/userRepo"
	LogService "zymm/internal/service/log_service"
	"zymm/utils"
)

const apiVersion = "v1"
const apiPrefix = "employee"

func routeAppended(newRoute string) string {
	return utils.ApiRoute(apiVersion, apiPrefix, newRoute)
}

func EmployeeHandlerDelegate(mux *http.ServeMux) {
	mux.HandleFunc(routeAppended("create"), AuthMiddleware(createEmployee))
	mux.HandleFunc(routeAppended("getEmployeeById"), AuthMiddleware(getEmployeeById))
	mux.HandleFunc(routeAppended("getAllEmployeeByGymId"), AuthMiddleware(getAllEmployeeByGymId))
	// mux.HandleFunc(routeAppended("removeEmployee"), AuthMiddleware(removeEmployee))
}

func createEmployee(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("createEmployee was called")
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
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

	// Decode JSON body
	var req models.EmployeeRegistrationApiRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	validationRequestErr := bussinessAuth.ValidateEmployeeRegistrationRequest(req)
	if validationRequestErr != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, validationRequestErr.Error())
		return
	}

	userExistsWithMobile, userExistsWithMobileErr := userRepo.CheckUserExistsByMobile(req.Mobile)
	if userExistsWithMobileErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, userExistsWithMobileErr.Error())
		return
	}

	if userExistsWithMobile {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "User already exists with mobile number "+req.Mobile)
		return
	}

	if req.Email != nil {
		userExistsWithEmail, userExistsWithEmailErr := userRepo.CheckUserExistsByEmail(*req.Email)
		if userExistsWithEmailErr != nil {
			utils.SendErrorResponse(w, http.StatusExpectationFailed, userExistsWithEmailErr.Error())
			return
		}

		if userExistsWithEmail {
			utils.SendErrorResponse(w, http.StatusExpectationFailed, "User already exists with email "+*req.Email)
			return
		}
	}

	roleTypePtr := businessRoleType.GetRoleTypeFromInt(claims.RoleId)
	if roleTypePtr == nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Invalid role")
		return
	}

	roleType := *roleTypePtr

	if businessRoles.IsNotAllowedToManageEmployees(roleType) {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Access Denied.")
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
		RoleId:     roleType.Int(),
	})

	if userInsertionError != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error creating employee: "+userInsertionError.Error())
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

	var gymId int = 0
	if roleType == businessRoleType.RoleOwner {
		gymDetails, gymDetailsErr := gymRepo.GetGymWithAdditionalDataByOwnerId(currentUserId)
		if gymDetailsErr != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, gymDetailsErr.Error())
			return
		}
		gymId = gymDetails.GymId
	} else {
		employeeDetails, employeeDetailsErr := employeesRepo.GetEmployeeByUserId(currentUserId)
		if employeeDetailsErr != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, employeeDetailsErr.Error())
			return
		}
		gymId = employeeDetails.GymId
	}

	empModel := &models.EmployeeModel{
		EmployeeId: 0,
		UserId:     insertedUserRecord.UserId,
		GymId:      gymId,
		CreatedBy:  currentUserId,
	}

	employeeInsert, employeeInsertErr := employeesRepo.InsertEmployees(*empModel)
	if employeeInsertErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, employeeInsertErr.Error())
		return
	}

	if employeeInsert == nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Employee didn't insert.")
		return
	}

	response := &models.EmployeeRegistrationApiResponseModel{
		Mobile:   req.Mobile,
		Password: req.Password,
	}

	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func getAllEmployeeByGymId(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getAllEmployeeByGymId was called")
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
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

	queryParams := r.URL.Query()
	gymIdRawStr := queryParams.Get("gymId")

	gymIdConverted := utils.ConvertStringToInt(gymIdRawStr)
	if gymIdConverted == nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Request was not proper.")
		return
	}

	gymId := *gymIdConverted

	if gymId <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Request was not proper.")
		return
	}

	roleTypePtr := businessRoleType.GetRoleTypeFromInt(claims.RoleId)
	if roleTypePtr == nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Invalid role")
		return
	}

	roleType := *roleTypePtr

	if businessRoles.IsNotAllowedToManageEmployees(roleType) {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Access Denied.")
		return
	}

	employeesData, employeesDataErr := employeesRepo.GetAllEmployees(gymId)
	if employeesDataErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, employeesDataErr.Error())
		return
	}

	if employeesData == nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "No data found.")
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, employeesData)
}

func getEmployeeById(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getEmployeeById was called")
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
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

	queryParams := r.URL.Query()
	employeeIdRawStr := queryParams.Get("employeeId")

	employeeIdConverted := utils.ConvertStringToInt(employeeIdRawStr)
	if employeeIdConverted == nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Request was not proper.")
		return
	}

	employeeId := *employeeIdConverted

	if employeeId <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Request was not proper.")
		return
	}

	roleTypePtr := businessRoleType.GetRoleTypeFromInt(claims.RoleId)
	if roleTypePtr == nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Invalid role")
		return
	}

	roleType := *roleTypePtr

	if businessRoles.IsNotAllowedToManageEmployees(roleType) {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Access Denied.")
		return
	}

	employeeData, employeeDataErr := employeesRepo.GetEmployeeByEmployeeId(employeeId)
	if employeeDataErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, employeeDataErr.Error())
		return
	}

	if employeeData == nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "No data found.")
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, employeeData)
}
