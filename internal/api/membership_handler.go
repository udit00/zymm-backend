package api

import (
	"encoding/json"
	"net/http"

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

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
