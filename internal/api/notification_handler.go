package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	bussinessAuth "zymm/internal/business/auth"
	businessNotificationType "zymm/internal/business/notification/notification_type"
	"zymm/internal/models"
	gymrepo "zymm/internal/repository/gym_repo"
	notificationRepo "zymm/internal/repository/notification_repo"
	LogService "zymm/internal/service/log_service"
	"zymm/utils"
)

const notificationApiVersion = "v1"
const notificationApiPrefix = "notification"

func notificationRouteAppended(newRoute string) string {
	return utils.ApiRoute(notificationApiVersion, notificationApiPrefix, newRoute)
}

func NotificationHandlerDelegate(mux *http.ServeMux) {
	mux.HandleFunc(notificationRouteAppended("createGeneralNotification"), AuthMiddleware(createGeneralNotification))
	mux.HandleFunc(notificationRouteAppended("markAsRead"), AuthMiddleware(markNotificationAsRead))
	mux.HandleFunc(notificationRouteAppended("getMyNotifications"), AuthMiddleware(getMyNotifications))
}

func createGeneralNotification(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("createGeneralNotification was called")
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.CreateGeneralNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	title := strings.TrimSpace(req.Title)
	description := strings.TrimSpace(req.Description)
	if title == "" || description == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Title and description are required")
		return
	}

	userIdsRaw := strings.TrimSpace(req.UserIds)
	var userIds []int
	if strings.EqualFold(userIdsRaw, "ALL") {
		if req.GymId == nil || *req.GymId <= 0 {
			utils.SendErrorResponse(w, http.StatusBadRequest, "gymId is required when userIds is ALL")
			return
		}

		gymUserIds, err := gymrepo.GetActiveUserIdsForGym(*req.GymId)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching users for gym: "+err.Error())
			return
		}
		userIds = gymUserIds
	} else {
		if userIdsRaw == "" {
			utils.SendErrorResponse(w, http.StatusBadRequest, "userIds is required")
			return
		}

		rawUserIds := strings.Split(userIdsRaw, ",")
		seen := make(map[int]struct{})
		for _, rawId := range rawUserIds {
			trimmed := strings.TrimSpace(rawId)
			if trimmed == "" {
				continue
			}

			id, err := strconv.Atoi(trimmed)
			if err != nil || id <= 0 {
				utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid userIds format")
				return
			}

			if _, exists := seen[id]; exists {
				continue
			}
			seen[id] = struct{}{}
			userIds = append(userIds, id)
		}
	}

	if len(userIds) == 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "At least one valid userId is required")
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

	for _, userId := range userIds {
		if err := notificationRepo.InsertNotification(userId, businessNotificationType.NotificationGeneralInfo, title, description, claims.UserId); err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, "Error creating notifications: "+err.Error())
			return
		}
	}

	response := map[string]interface{}{
		"notificationsCreated": len(userIds),
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func markNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("markNotificationAsRead was called")
	if r.Method != http.MethodPost {
		utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.MarkNotificationAsReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.NotificationId <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "notificationId is required")
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

	err := notificationRepo.MarkNotificationAsRead(req.NotificationId, claims.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.SendErrorResponse(w, http.StatusNotFound, "Notification not found")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error updating notification: "+err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, map[string]string{"message": "Notification marked as read"})
}

func getMyNotifications(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getMyNotifications was called")
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

	notifications, err := notificationRepo.GetNotificationsByUserId(claims.UserId)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching notifications: "+err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, notifications)
}
