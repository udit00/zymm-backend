package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	bussinessAuth "zymm/internal/business/auth"
	businessMembership "zymm/internal/business/membership"
	bussinessMembershipRequest "zymm/internal/business/membership/membership_request_action_type"
	businessRoles "zymm/internal/business/roles"
	businessRoleType "zymm/internal/business/roles/roles_type"
	"zymm/internal/models"
	employeesRepo "zymm/internal/repository/employees_repo"
	gymRepo "zymm/internal/repository/gym_repo"
	membershipRepo "zymm/internal/repository/membership_repo"
	notificationRepo "zymm/internal/repository/notification_repo"
	LogService "zymm/internal/service/log_service"
	"zymm/utils"
)

const membershipApiVersion = "v1"
const membershipApiPrefix = "membership"

func membershipRouteAppended(newRoute string) string {
	return utils.ApiRoute(membershipApiVersion, membershipApiPrefix, newRoute)
}

func MembershipHandlerDelegate(mux *http.ServeMux) {
	mux.HandleFunc(membershipRouteAppended("upsertPlan"), AuthMiddleware(upsertMembershipPlan))
	mux.HandleFunc(membershipRouteAppended("requestPlan"), AuthMiddleware(requestMembershipByUserToGym))
	mux.HandleFunc(membershipRouteAppended("planHistory"), AuthMiddleware(userMembershipHistory))
	mux.HandleFunc(membershipRouteAppended("getAllMemberships"), AuthMiddleware(allMembershipHistoryForGymOwnerAndManagers))
	mux.HandleFunc(membershipRouteAppended("takeActionOnMembership"), AuthMiddleware(takeActionOnMembership))
	mux.HandleFunc(membershipRouteAppended("getPlanDetails"), AuthMiddleware(getPlanDetails))
	mux.HandleFunc(membershipRouteAppended("getAllPlansByGymId"), AuthMiddleware(getAllPlansByGymId))
	mux.HandleFunc(membershipRouteAppended("cancelMembershipRequest"), AuthMiddleware(cancelMembershipRequestByUserToGym))
}

func upsertMembershipPlan(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("createMembershipPlan was called")
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.UpsertPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

	// Validate request and check if user and gym exist in database
	validationErr := businessMembership.ValidateUpsertPlanWithDBChecks(req, currentUserId)
	if validationErr != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	var planChangeLog *models.PlanChangeLog
	if *req.PlanId == 0 {
		isActive := true
		// build plan record
		plan := models.PlanRecord{
			PlanBanner:   req.PlanBanner,
			PlanName:     req.PlanName,
			PlanDesc:     req.PlanDesc,
			PlanPrice:    req.PlanPrice,
			PlanDuration: req.PlanDuration,
			IsActive:     &isActive,
			CreatedBy:    &currentUserId,
			GymId:        &req.GymId,
		}

		insertedId, err := membershipRepo.InsertPlan(plan)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, "Error creating plan: "+err.Error())
			return
		}

		// insert plan change log - for creation, old fields can be empty/defaults
		planBannerStr := derefString(req.PlanBanner)
		planChangeLog = &models.PlanChangeLog{
			PlanId:          insertedId,
			ChangedBy:       &currentUserId,
			ChangeType:      "create",
			ChangeDetails:   "Created",
			OldPlanBanner:   req.PlanBanner,
			OldPlanName:     req.PlanName,
			OldPlanDesc:     req.PlanDesc,
			OldPlanPrice:    req.PlanPrice,
			OldPlanDuration: req.PlanDuration,
			OldIsActive:     true,
			NewPlanBanner:   &planBannerStr,
			NewPlanName:     req.PlanName,
			NewPlanDesc:     req.PlanDesc,
			NewPlanPrice:    req.PlanPrice,
			NewPlanDuration: req.PlanDuration,
			NewIsActive:     true,
		}

		resp := map[string]interface{}{"planId": insertedId}
		utils.SendSuccessResponse(w, http.StatusOK, resp)
	} else {
		oldPlanDetails, oldPlanDetailsErr := membershipRepo.GetPlanById(*req.PlanId)
		if oldPlanDetailsErr != nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, oldPlanDetailsErr.Error())
			return
		}
		if oldPlanDetails == nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, "Plan Details not found.")
			return
		}

		// planBannerStr := derefString(req.PlanBanner)

		newPlanDetails := models.PlanRecord{
			PlanId:       *req.PlanId,
			PlanBanner:   req.PlanBanner,
			PlanName:     req.PlanName,
			PlanDesc:     req.PlanDesc,
			PlanPrice:    req.PlanPrice,
			PlanDuration: req.PlanDuration,
			IsActive:     &req.IsActive,
			GymId:        &req.GymId,
			CreatedBy:    oldPlanDetails.CreatedBy,
			CreatedAt:    oldPlanDetails.CreatedAt,
		}

		updatePlanErr := membershipRepo.UpdatePlan(newPlanDetails)
		if updatePlanErr != nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, "Error updating plan. "+updatePlanErr.Error())
			return
		}

		planChangeLog = &models.PlanChangeLog{
			PlanId:          req.PlanId,
			ChangedBy:       &currentUserId,
			ChangeType:      "update",
			ChangeDetails:   businessMembership.GetPlanChangeDiffDetails(*oldPlanDetails, newPlanDetails),
			OldPlanBanner:   oldPlanDetails.PlanBanner,
			OldPlanName:     oldPlanDetails.PlanName,
			OldPlanDesc:     oldPlanDetails.PlanDesc,
			OldPlanPrice:    oldPlanDetails.PlanPrice,
			OldPlanDuration: oldPlanDetails.PlanDuration,
			OldIsActive:     *oldPlanDetails.IsActive,
			NewPlanBanner:   newPlanDetails.PlanBanner,
			NewPlanName:     newPlanDetails.PlanName,
			NewPlanDesc:     newPlanDetails.PlanDesc,
			NewPlanPrice:    newPlanDetails.PlanPrice,
			NewPlanDuration: newPlanDetails.PlanDuration,
			NewIsActive:     *newPlanDetails.IsActive,
		}
		utils.SendSuccessResponse(w, http.StatusOK, nil)
	}

	if planChangeLog != nil {
		_, planChangeInsertionLogErr := membershipRepo.InsertPlanChangeLog(*planChangeLog)
		if planChangeInsertionLogErr != nil {
			// log error but do not fail the creation (or choose to fail depending on desired semantics)
			LogService.LogError("Failed to insert plan change log: ", planChangeInsertionLogErr)
		}
	}

}

func requestMembershipByUserToGym(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("requestMembershipByUserToGym was called")
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.RequestPlanFromUserToGymModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	claims, ok := r.Context().Value(ctxClaimDataKey).(*bussinessAuth.MyCustomClaims)
	if !ok || claims == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	createdBy := claims.UserId
	if createdBy < 0 {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "invalid user id in context")
		return
	}

	roleType := businessRoleType.GetRoleTypeFromInt(claims.RoleId)
	if roleType == nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Role type not found.")
		return
	}

	if *roleType != businessRoleType.RoleMember {
		utils.SendErrorResponse(w, http.StatusBadRequest, "You cannot request a membership.")
		return
	}

	planDetails, planDetailErr := membershipRepo.GetPlanById(req.PlanId)
	if planDetailErr != nil || planDetails == nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, planDetailErr.Error())
		return
	}

	latestUserMembershipDetails, userMembershipError := membershipRepo.GetUserMembershipByUserId(createdBy)
	if latestUserMembershipDetails != nil {
		layout := time.RFC3339
		startTime, err := time.Parse(layout, latestUserMembershipDetails.StartDate)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		endTime, err := time.ParseInLocation(layout, latestUserMembershipDetails.EndDate, time.UTC)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		now := time.Now().UTC()
		isBetween := now.After(startTime) && now.Before(endTime)

		fmt.Println("Start:", startTime)
		fmt.Println("End:", endTime)
		fmt.Println("Now:", now)

		// Check if current time is between start and end
		if latestUserMembershipDetails.MembershipStatus == "A" && isBetween {
			fmt.Println("✅ Current time is between start and end")
			utils.SendErrorResponse(w, http.StatusBadRequest, "You already have an active plan.")
			return
		} else if latestUserMembershipDetails.MembershipStatus == "P" {
			fmt.Println("✅ Current time is between start and end")
			utils.SendErrorResponse(w, http.StatusBadRequest, "You already have an pending plan request, please wait for the gym staff to make some changes on it.")
			return
		} else {
			fmt.Println("❌ Current time is outside the range")
			insertMembership(w, createdBy, planDetails.PlanId, *planDetails.GymId, planDetails.PlanDuration)

		}
	} else if userMembershipError == sql.ErrNoRows {
		insertMembership(w, createdBy, planDetails.PlanId, *planDetails.GymId, planDetails.PlanDuration)
	} else {
		utils.SendErrorResponse(w, http.StatusBadRequest, userMembershipError.Error())
		return
	}

}

func cancelMembershipRequestByUserToGym(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("cancelMembershipRequestByUserToGym was called")
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.CancelRequestMembershipModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	claims, ok := r.Context().Value(ctxClaimDataKey).(*bussinessAuth.MyCustomClaims)
	if !ok || claims == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	createdBy := claims.UserId
	if createdBy < 0 {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "invalid user id in context")
		return
	}

	membershipDetails, membershipDetailsErr := membershipRepo.GetUserMembershipByMembershipId(req.MembershipId)
	if membershipDetailsErr != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, membershipDetailsErr.Error())
		return
	}

	if membershipDetails == nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "No data found. Please try again later.")
		return
	}

	if membershipDetails.UserId != createdBy {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Access Denied")
		return
	}

	if membershipDetails.IsActive {
		cancelMembershipRequestErr := membershipRepo.CancelMembershipRequest(membershipDetails.MembershipId)
		if cancelMembershipRequestErr != nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, cancelMembershipRequestErr.Error())
			return
		}
		utils.SendSuccessResponse(w, http.StatusAccepted, nil)
	} else {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Request is already cancelled.")
		return
	}

}

func insertMembership(w http.ResponseWriter, createdBy int, gymId int, planId int, planDuration int) {
	umRequest := models.UserMembership{
		UserId:           createdBy,
		PlanId:           planId,
		StartDate:        utils.ConvertTimeToDBDateTime(utils.GetCurrentDateTime()),
		EndDate:          utils.ConvertTimeToDBDateTime(utils.GetCurrentDateTime().Add(time.Hour * 24 * time.Duration(planDuration))),
		IsActive:         true,
		MembershipStatus: "P",
	}
	memberShipId, memberInsertionErr := membershipRepo.InsertMembershipRequest(umRequest)
	if memberInsertionErr != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, memberInsertionErr.Error())
		return
	}
	if memberShipId != nil && *memberShipId > 0 {
		go func(gymId int, userId int) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered in SendNotificationForMembershipRequested: %v", r)
				}
			}()

			notificationRepo.SendNotificationForMembershipRequested(gymId, userId)
		}(gymId, umRequest.UserId)
		utils.SendSuccessResponse(w, http.StatusAccepted, "Request ID: "+strconv.Itoa((*memberShipId))+" has been generated, wait for the gym staff to approve it.")
	} else {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Something went wrong, please try again later.")
	}
}

func userMembershipHistory(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("userMembershipHistory was called")
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

	filterByStr := r.URL.Query().Get("filterBy")
	if filterByStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "filterBy is required")
		return
	}

	if filterByStr != "P" && filterByStr != "A" && filterByStr != "R" && filterByStr != "ALL" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "filterBy can only be [P,R,A,ALL], not "+filterByStr)
		return
	}

	var filterType membershipRepo.MembershipStatus = membershipRepo.All

	switch filterByStr {
	case "P":
		filterType = membershipRepo.Pending
	case "A":
		filterType = membershipRepo.Approved
	case "R":
		filterType = membershipRepo.Rejected
	default:
		filterType = membershipRepo.All
	}

	memberships, membershipHistoryErr := membershipRepo.GetAllMembershipPlansRequestByUserId(currentUserId, filterType)
	if membershipHistoryErr != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, membershipHistoryErr.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, memberships)
}

func allMembershipHistoryForGymOwnerAndManagers(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("allMembershipHistoryForGymOwnerAndManagers was called")
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

	if businessRoles.IsNotAllowedToTakeActionOnMemberships(businessRoleType.RoleType(claims.RoleId)) {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Access Denied")
		return
	}

	filterByStr := r.URL.Query().Get("filterBy")
	if filterByStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "filterBy is required")
		return
	}

	if filterByStr != "P" && filterByStr != "A" && filterByStr != "R" && filterByStr != "ALL" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "filterBy can only be [P,R,A,ALL], not "+filterByStr)
		return
	}

	roleTypePtr := businessRoleType.GetRoleTypeFromInt(claims.RoleId)
	if roleTypePtr == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Access Denied")
		return
	}

	roleType := *roleTypePtr

	if roleType != businessRoleType.RoleOwner && roleType != businessRoleType.RoleManager {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Access Denied.")
		return
	}

	var gymId int
	if roleType == businessRoleType.RoleOwner {
		gymDetails, gymDetailsErr := gymRepo.GetGymWithAdditionalDataByOwnerId(claims.UserId)
		if gymDetailsErr != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, gymDetailsErr.Error())
			return
		}
		gymId = gymDetails.GymId
	} else {
		employeeDetails, employeeDetailsErr := employeesRepo.GetEmployeeByUserId(claims.UserId)
		if employeeDetailsErr != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, employeeDetailsErr.Error())
			return
		}
		gymId = employeeDetails.GymId
	}

	if gymId <= 0 {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Gym Id not found for the user.")
		return
	}

	var filterType membershipRepo.MembershipStatus = membershipRepo.All

	switch filterByStr {
	case "P":
		filterType = membershipRepo.Pending
	case "A":
		filterType = membershipRepo.Approved
	case "R":
		filterType = membershipRepo.Rejected
	default:
		filterType = membershipRepo.All
	}

	memberships, membershipHistoryErr := membershipRepo.GetAllMembershipPlansRequest(gymId, filterType)
	if membershipHistoryErr != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, membershipHistoryErr.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, memberships)
}

func getAllPlansByGymId(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getAllPlansByGymId was called")
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

	gymIdStr := r.URL.Query().Get("gymId")
	if gymIdStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "gymId is required")
		return
	}

	gymId, err := strconv.Atoi(gymIdStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid gymId")
		return
	}

	planDetails, planDetailsErr := membershipRepo.GetPlansByGymId(gymId)
	if planDetailsErr != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, planDetailsErr.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, planDetails)
}

func getPlanDetails(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getPlanDetails was called")
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

	planIdStr := r.URL.Query().Get("planId")
	if planIdStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "planId is required")
		return
	}

	planId, err := strconv.Atoi(planIdStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid planId")
		return
	}

	planDetails, planDetailsErr := membershipRepo.GetPlanById(planId)
	if planDetailsErr != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, planDetailsErr.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, planDetails)
}

func takeActionOnMembership(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("takeActionOnMembership was called")
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.TakeActionOnMembershipRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

	rawUserRoleType := businessRoleType.GetRoleTypeFromInt(claims.RoleId)
	if rawUserRoleType == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Role type was null")
		return
	}

	userRole := *rawUserRoleType
	if businessRoles.IsNotAllowedToTakeActionOnMemberships(userRole) {
		utils.SendErrorResponse(w, http.StatusForbidden, "Access Denied.")
		return
	}

	membershipRequestDetails, membershipRequestDetailsErr := membershipRepo.GetUserMembershipByMembershipId(req.MembershipId)
	if membershipRequestDetailsErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Couldn't get membership data: "+membershipRequestDetailsErr.Error())
		return
	}
	if membershipRequestDetails == nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Couldn't get membership data: data was nil")
		return
	}

	planDetails, planDetailErr := membershipRepo.GetPlanById(membershipRequestDetails.PlanId)
	if planDetailErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Couldn't get plan data: "+planDetailErr.Error())
		return
	}

	if planDetails == nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, "Couldn't get planDetails data: data was nil")
		return
	}

	var action bussinessMembershipRequest.ActionType = bussinessMembershipRequest.Approve
	if req.ActionTaken != "A" && req.ActionTaken != "R" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Action Type can only be A for Accept or R for reject.")
		return
	}

	if req.ActionTaken == "R" {
		action = bussinessMembershipRequest.Reject
	}
	actionTakenId, actionTakingErr := membershipRepo.TakeActionOnMembershipRequest(
		membershipRequestDetails.UserId,
		membershipRequestDetails.MembershipId,
		action,
		currentUserId,
	)

	if actionTakingErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, actionTakingErr.Error())
		return
	}

	updateMembershipOfTheUserErr := membershipRepo.UpdateActionTakenOnUserMembership(
		membershipRequestDetails.MembershipId,
		action,
	)

	if updateMembershipOfTheUserErr != nil {
		utils.SendErrorResponse(w, http.StatusExpectationFailed, updateMembershipOfTheUserErr.Error())
		return
	}

	go func(nById int, nForId int, pd models.PlanRecord, act bussinessMembershipRequest.ActionType) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in SendNotificationForMembershipResponse: %v", r)
			}
		}()
		notificationRepo.SendNotificationForMembershipResponseTakenByGymOwnerManagers(nById, nForId, pd, act)
	}(currentUserId, membershipRequestDetails.UserId, *planDetails, action)

	type finalResponse struct {
		actionTakenId int
	}
	utils.SendSuccessResponse(w, http.StatusOK, finalResponse{actionTakenId: *actionTakenId})
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
