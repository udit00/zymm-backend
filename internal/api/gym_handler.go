package api

import (
	"database/sql"
	"net/http"
	"strconv"

	bussinessAuth "zymm/internal/business/auth"
	businessRoles "zymm/internal/business/roles"
	businessRoleType "zymm/internal/business/roles/roles_type"
	employeesrepo "zymm/internal/repository/employees_repo"
	gymRepo "zymm/internal/repository/gym_repo"
	LogService "zymm/internal/service/log_service"
	"zymm/utils"
)

const gymApiVersion = "v1"
const gymApiPrefix = "gym"

func gymRouteAppended(newRoute string) string {
	return utils.ApiRoute(gymApiVersion, gymApiPrefix, newRoute)
}

func GymHandlerDelegate(mux *http.ServeMux) {
	mux.HandleFunc(gymRouteAppended("getAllGyms"), AuthMiddleware(getAllGyms))
	mux.HandleFunc(gymRouteAppended("searchGyms"), AuthMiddleware(searchGyms))
	mux.HandleFunc(gymRouteAppended("getGymData"), AuthMiddleware(getGymData))
	mux.HandleFunc(gymRouteAppended("activeMembers"), AuthMiddleware(getActiveGymMembers))
	mux.HandleFunc(gymRouteAppended("getGymMembers"), AuthMiddleware(getGymMembers))
}

func getAllGyms(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getAllGyms was called")
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	claims, ok := r.Context().Value(ctxClaimDataKey).(*bussinessAuth.MyCustomClaims)
	if !ok || claims == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if claims.UserId <= 0 {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "invalid user id in context")
		return
	}

	allGyms, allGymErr := gymRepo.GetAllGymWithAdditionalData()
	if allGymErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, allGymErr.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, allGyms)

}

func searchGyms(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("searchGyms was called")
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	claims, ok := r.Context().Value(ctxClaimDataKey).(*bussinessAuth.MyCustomClaims)
	if !ok || claims == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if claims.UserId <= 0 {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "invalid user id in context")
		return
	}

	searchQuery := r.URL.Query().Get("query")
	// Empty search query returns all gyms
	if searchQuery == "" {
		allGyms, allGymErr := gymRepo.GetAllGymWithAdditionalData()
		if allGymErr != nil {
			utils.SendErrorResponse(w, http.StatusExpectationFailed, allGymErr.Error())
			return
		}
		utils.SendSuccessResponse(w, http.StatusOK, allGyms)
		return
	}

	// Search gyms by name
	gyms, err := gymRepo.SearchGymsByName(searchQuery)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error searching gyms: "+err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, gyms)
}

func getGymData(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getGymData was called")
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	claims, ok := r.Context().Value(ctxClaimDataKey).(*bussinessAuth.MyCustomClaims)
	if !ok || claims == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if claims.UserId <= 0 {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "invalid user id in context")
		return
	}

	gymIdStr := r.URL.Query().Get("gymId")
	if gymIdStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "gymId is required")
		return
	}

	gymId, err := strconv.Atoi(gymIdStr)
	if err != nil || gymId <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid gymId")
		return
	}

	gymRecord, err := gymRepo.GetGymWithAdditionalDataByGymId(gymId)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendErrorResponse(w, http.StatusNotFound, "Gym not found")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching gym: "+err.Error())
		return
	}
	utils.SendSuccessResponse(w, http.StatusOK, gymRecord)

}

func getActiveGymMembers(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getActiveGymMembers was called")
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	claims, ok := r.Context().Value(ctxClaimDataKey).(*bussinessAuth.MyCustomClaims)
	if !ok || claims == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if claims.UserId <= 0 {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "invalid user id in context")
		return
	}

	roleTypePtr := businessRoleType.GetRoleTypeFromInt(claims.RoleId)
	if roleTypePtr == nil {
		utils.SendErrorResponse(w, http.StatusForbidden, "Invalid role type")
		return
	}

	roleType := *roleTypePtr
	if businessRoles.IsNotAllowedToManagePlans(roleType) {
		utils.SendErrorResponse(w, http.StatusForbidden, "Access Denied")
		return
	}

	gymIdStr := r.URL.Query().Get("gymId")
	if gymIdStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "gymId is required")
		return
	}

	gymId, err := strconv.Atoi(gymIdStr)
	if err != nil || gymId <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid gymId")
		return
	}

	gymRecord, err := gymRepo.GetGymById(gymId)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendErrorResponse(w, http.StatusNotFound, "Gym not found")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching gym: "+err.Error())
		return
	}

	switch roleType {
	case businessRoleType.RoleOwner:
		if gymRecord.CreatedBy != claims.UserId {
			utils.SendErrorResponse(w, http.StatusForbidden, "Access Denied")
			return
		}
	case businessRoleType.RoleManager:
		employee, err := employeesrepo.GetEmployeeByUserId(claims.UserId)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.SendErrorResponse(w, http.StatusForbidden, "Access Denied")
				return
			}
			utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching employee info: "+err.Error())
			return
		}
		if employee.GymId != gymId {
			utils.SendErrorResponse(w, http.StatusForbidden, "Access Denied")
			return
		}
	default:
		utils.SendErrorResponse(w, http.StatusForbidden, "Access Denied")
		return
	}

	members, err := gymRepo.GetActiveMembershipDetailsForGym(gymId)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching memberships: "+err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, members)
}

// getGymMembers returns all gym members with their active plans and feedback ratings
// Available to owners and managers only
func getGymMembers(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getGymMembers was called")

	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get claims from context
	claims, ok := r.Context().Value(ctxClaimDataKey).(*bussinessAuth.MyCustomClaims)
	if !ok || claims == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	currentUserId := claims.UserId
	roleType := businessRoleType.RoleType(claims.RoleId)

	// Check if user has permission to view members
	if businessRoles.IsNotAllowedToManageEmployees(roleType) {
		utils.SendErrorResponse(w, http.StatusForbidden, "Access Denied: Only owners and managers can view members")
		return
	}

	// Get gymId based on user role
	var gymId int
	switch roleType {
	case businessRoleType.RoleOwner:
		// Owner - get gym from gym table
		gym, err := gymRepo.GetGymWithAdditionalDataByOwnerId(currentUserId)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.SendErrorResponse(w, http.StatusNotFound, "No gym found for this owner")
				return
			}
			utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching gym info: "+err.Error())
			return
		}
		gymId = gym.GymId
	case businessRoleType.RoleManager:
		// Manager - get gym from employee table
		employee, err := employeesrepo.GetEmployeeWithUserDetailsByEmployeeId(currentUserId)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.SendErrorResponse(w, http.StatusNotFound, "No employee record found for this manager")
				return
			}
			utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching employee info: "+err.Error())
			return
		}
		if !employee.IsActive {
			utils.SendErrorResponse(w, http.StatusExpectationFailed, "You are not an active user.")
			return
		}
		gymId = employee.GymId
	default:
		utils.SendErrorResponse(w, http.StatusForbidden, "Access Denied")
		return
	}

	// Get gym members with their active plans and ratings
	members, err := gymRepo.GetGymMembersWithDetails(gymId)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching gym members: "+err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, members)
}
