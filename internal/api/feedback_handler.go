package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	bussinessAuth "zymm/internal/business/auth"
	"zymm/internal/models"
	feedbackRepo "zymm/internal/repository/feedback_repo"
	LogService "zymm/internal/service/log_service"
	"zymm/utils"
)

const feedbackApiVersion = "v1"
const feedbackApiPrefix = "feedback"

func feedbackRouteAppended(newRoute string) string {
	return utils.ApiRoute(feedbackApiVersion, feedbackApiPrefix, newRoute)
}

func FeedbackHandlerDelegate(mux *http.ServeMux) {
	mux.HandleFunc(feedbackRouteAppended("create"), AuthMiddleware(createFeedback))
	mux.HandleFunc(feedbackRouteAppended("get"), AuthMiddleware(getFeedback))
	mux.HandleFunc(feedbackRouteAppended("getAllByGymId"), AuthMiddleware(getAllFeedbacksByGymId))
	mux.HandleFunc(feedbackRouteAppended("update"), AuthMiddleware(updateFeedback))
	mux.HandleFunc(feedbackRouteAppended("delete"), AuthMiddleware(deleteFeedback))
}

func createFeedback(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("createFeedback was called")
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.CreateFeedbackRequest
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

	// Validate rating is between 1 and 5 if provided
	if req.Rating != nil && (*req.Rating < 1 || *req.Rating > 5) {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Rating must be between 1 and 5")
		return
	}

	// Validate gymId
	if req.GymId <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid gymId")
		return
	}

	feedback := models.FeedbackRecord{
		Rating:    req.Rating,
		Comments:  req.Comments,
		GymId:     req.GymId,
		CreatedBy: currentUserId,
	}

	insertedId, err := feedbackRepo.InsertFeedback(feedback)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error creating feedback: "+err.Error())
		return
	}

	resp := map[string]interface{}{"feedbackId": insertedId}
	utils.SendSuccessResponse(w, http.StatusOK, resp)
}

func getFeedback(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getFeedback was called")
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	claims, ok := r.Context().Value(ctxClaimDataKey).(*bussinessAuth.MyCustomClaims)
	if !ok || claims == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	feedbackIdStr := r.URL.Query().Get("feedbackId")
	if feedbackIdStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "feedbackId is required")
		return
	}

	feedbackId, err := strconv.Atoi(feedbackIdStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid feedbackId")
		return
	}

	feedback, err := feedbackRepo.GetFeedbackById(feedbackId)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendErrorResponse(w, http.StatusNotFound, "Feedback not found")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching feedback: "+err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, feedback)
}

func getAllFeedbacksByGymId(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getAllFeedbacksByGymId was called")
	if r.Method != http.MethodGet {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	claims, ok := r.Context().Value(ctxClaimDataKey).(*bussinessAuth.MyCustomClaims)
	if !ok || claims == nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
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

	feedbacks, err := feedbackRepo.GetAllFeedbacksByGymId(gymId)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching feedbacks: "+err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, feedbacks)
}

func updateFeedback(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("updateFeedback was called")
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.UpdateFeedbackRequest
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

	// Validate rating is between 1 and 5 if provided
	if req.Rating != nil && (*req.Rating < 1 || *req.Rating > 5) {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Rating must be between 1 and 5")
		return
	}

	// Check if feedback exists and belongs to the user
	existingFeedback, err := feedbackRepo.GetFeedbackById(req.FeedbackId)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendErrorResponse(w, http.StatusNotFound, "Feedback not found")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching feedback: "+err.Error())
		return
	}

	// Check if the feedback belongs to the current user
	if existingFeedback.CreatedBy != currentUserId {
		utils.SendErrorResponse(w, http.StatusForbidden, "You are not authorized to update this feedback")
		return
	}

	var commentStr *string = req.Comments
	if commentStr != nil {
		trimmed := strings.TrimSpace(*commentStr)
		if trimmed == "" {
			commentStr = nil
		}
	}

	// Update the feedback
	feedback := models.FeedbackRecord{
		FeedbackId: req.FeedbackId,
		Rating:     req.Rating,
		Comments:   commentStr,
		GymId:      existingFeedback.GymId,
		CreatedBy:  existingFeedback.CreatedBy,
		CreatedAt:  existingFeedback.CreatedAt,
	}

	err = feedbackRepo.UpdateFeedback(feedback)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error updating feedback: "+err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, map[string]string{"message": "Feedback updated successfully"})
}

func deleteFeedback(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("deleteFeedback was called")
	if r.Method != http.MethodDelete {
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

	feedbackIdStr := r.URL.Query().Get("feedbackId")
	if feedbackIdStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "feedbackId is required")
		return
	}

	feedbackId, err := strconv.Atoi(feedbackIdStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid feedbackId")
		return
	}

	// Check if feedback exists and belongs to the user
	existingFeedback, err := feedbackRepo.GetFeedbackById(feedbackId)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendErrorResponse(w, http.StatusNotFound, "Feedback not found")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching feedback: "+err.Error())
		return
	}

	// Check if the feedback belongs to the current user
	if existingFeedback.CreatedBy != currentUserId {
		utils.SendErrorResponse(w, http.StatusForbidden, "You are not authorized to delete this feedback")
		return
	}

	err = feedbackRepo.DeleteFeedback(feedbackId)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendErrorResponse(w, http.StatusNotFound, "Feedback not found")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error deleting feedback: "+err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, map[string]string{"message": "Feedback deleted successfully"})
}


