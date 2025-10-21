package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	businessMembership "zymm/internal/business/membership"
	"zymm/internal/models"
	membershipRepo "zymm/internal/repository/membership_repo"
	LogService "zymm/internal/service/log_service"
	"zymm/utils"
)

const membershipApiVersion = "v1"
const membershipApiPrefix = "membership"

func membershipRouteAppended(newRoute string) string {
	return utils.ApiRoute(membershipApiVersion, membershipApiPrefix, newRoute)
}

func MembershipHandlerDelegate(mux *http.ServeMux) {
	mux.HandleFunc(membershipRouteAppended("createPlan"), AuthMiddleware(createMembershipPlan))
	mux.HandleFunc(membershipRouteAppended("requestPlan"), AuthMiddleware(requestMembershipByUserToGym))
}

func createMembershipPlan(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("createMembershipPlan was called")
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.CreatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// extract user id from context
	uid := r.Context().Value(ctxUserIDKey)
	if uid == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	createdBy, ok := uid.(int)
	if !ok {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "invalid user id in context")
		return
	}

	// Validate request and check if user and gym exist in database
	validationErr := businessMembership.ValidateCreatePlanWithDBChecks(req, createdBy)
	if validationErr != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	isActive := true
	// build plan record
	plan := models.PlanRecord{
		PlanBanner:   req.PlanBanner,
		PlanName:     req.PlanName,
		PlanDesc:     req.PlanDesc,
		PlanPrice:    req.PlanPrice,
		PlanDuration: req.PlanDuration,
		IsActive:     &isActive,
		CreatedBy:    &createdBy,
		GymId:        &req.GymId,
	}

	insertedId, err := membershipRepo.InsertPlan(plan)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error creating plan: "+err.Error())
		return
	}

	// insert plan change log - for creation, old fields can be empty/defaults
	planBannerStr := derefString(req.PlanBanner)
	changeLog := models.PlanChangeLog{
		PlanId:          insertedId,
		ChangedBy:       &createdBy,
		ChangeType:      "create",
		ChangeDetails:   "Created",
		OldPlanBanner:   req.PlanBanner,
		OldPlanName:     req.PlanName,
		OldPlanDesc:     req.PlanDesc,
		OldPlanPrice:    req.PlanPrice,
		OldPlanDuration: req.PlanDuration,
		OldIsActive:     false,
		NewPlanBanner:   planBannerStr,
		NewPlanName:     req.PlanName,
		NewPlanDesc:     req.PlanDesc,
		NewPlanPrice:    &req.PlanPrice,
		NewPlanDuration: req.PlanDuration,
		NewIsActive:     true,
	}

	if err := membershipRepo.InsertPlanChangeLog(changeLog); err != nil {
		// log error but do not fail the creation (or choose to fail depending on desired semantics)
		LogService.LogError("Failed to insert plan change log: ", err)
	}

	resp := map[string]interface{}{"planId": insertedId}
	utils.SendSuccessResponse(w, http.StatusOK, resp)
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

	// extract user id from context
	uid := r.Context().Value(ctxUserIDKey)
	if uid == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	createdBy, ok := uid.(int)
	if !ok {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "invalid user id in context")
		return
	}

	planDetails, planDetailErr := membershipRepo.GetPlanById(req.PlanId)
	if planDetailErr != nil || planDetails == nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, planDetailErr.Error())
		return
	}

	latestUserMembershipDetails, userMembershipError := membershipRepo.GetUserMembershipByUserId(createdBy)
	if latestUserMembershipDetails != nil {
		startTime, err := time.Parse(time.RFC3339, latestUserMembershipDetails.StartDate)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		endTime, err := time.Parse(time.RFC3339, latestUserMembershipDetails.EndDate)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		now := time.Now().UTC() // use UTC for fair comparison
		fmt.Println("Current time (UTC):", now)

		// Check if current time is between start and end
		if now.After(startTime) && now.Before(endTime) {
			fmt.Println("✅ Current time is between start and end")
			utils.SendErrorResponse(w, http.StatusBadRequest, "You already have an active plan.")
			return
		} else {
			fmt.Println("❌ Current time is outside the range")
			insertMembership(w, createdBy, planDetails.PlanId, planDetails.PlanDuration)
		}
	} else if userMembershipError == sql.ErrNoRows {
		insertMembership(w, createdBy, planDetails.PlanId, planDetails.PlanDuration)
	} else {
		utils.SendErrorResponse(w, http.StatusBadRequest, userMembershipError.Error())
		return
	}

}

func insertMembership(w http.ResponseWriter, createdBy int, planId int, planDuration int) {
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
		utils.SendSuccessResponse(w, http.StatusAccepted, "Request ID: "+string(*memberShipId)+"has been generated, wait for the gym staff to approve it.")
	} else {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Something went wrong, please try again later.")
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
