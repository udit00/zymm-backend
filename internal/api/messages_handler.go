package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	bussinessAuth "zymm/internal/business/auth"
	messagesRepo "zymm/internal/repository/messages_repo"
	"zymm/internal/models"
	LogService "zymm/internal/service/log_service"
	"zymm/utils"
)

const messagesApiVersion = "v1"
const messagesApiPrefix = "messages"

func messagesRouteAppended(newRoute string) string {
	return utils.ApiRoute(messagesApiVersion, messagesApiPrefix, newRoute)
}

func MessagesHandlerDelegate(mux *http.ServeMux) {
	mux.HandleFunc(messagesRouteAppended("create"), AuthMiddleware(createMessage))
	mux.HandleFunc(messagesRouteAppended("update"), AuthMiddleware(updateMessage))
	mux.HandleFunc(messagesRouteAppended("delete"), AuthMiddleware(deleteMessage))
	mux.HandleFunc(messagesRouteAppended("getChatParticipants"), AuthMiddleware(getChatParticipants))
	mux.HandleFunc(messagesRouteAppended("getChatMessages"), AuthMiddleware(getChatMessages))
	mux.HandleFunc(messagesRouteAppended("getUnreadCount"), AuthMiddleware(getUnreadCount))
	mux.HandleFunc(messagesRouteAppended("getAvailableChatUsers"), AuthMiddleware(getAvailableChatUsers))
}

// createMessage creates a new message
func createMessage(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("createMessage was called")

	if r.Method != http.MethodPost {
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

	// Parse request body
	var req models.CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Validate
	if req.MessageForUserId <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid recipient user ID")
		return
	}

	if len(req.Comment) == 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Message content cannot be empty")
		return
	}

	// Cannot send message to yourself
	if req.MessageForUserId == currentUserId {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Cannot send message to yourself")
		return
	}

	// Create message
	message, err := messagesRepo.CreateMessage(req, currentUserId)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error creating message: "+err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusCreated, message)
}

// updateMessage updates an existing message
func updateMessage(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("updateMessage was called")

	if r.Method != http.MethodPost {
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

	// Parse request body
	var req models.UpdateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Validate
	if req.MessageId <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	if len(req.Comment) == 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Message content cannot be empty")
		return
	}

	// Update message
	err := messagesRepo.UpdateMessage(req, currentUserId)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendErrorResponse(w, http.StatusNotFound, "Message not found or you don't have permission to update it")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error updating message: "+err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Message updated successfully")
}

// deleteMessage soft deletes a message
func deleteMessage(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("deleteMessage was called")

	if r.Method != http.MethodPost {
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

	// Parse request body
	var req models.DeleteMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Validate
	if req.MessageId <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	// Delete message
	err := messagesRepo.SoftDeleteMessage(req.MessageId, currentUserId)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendErrorResponse(w, http.StatusNotFound, "Message not found or you don't have permission to delete it")
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error deleting message: "+err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, "Message deleted successfully")
}

// getChatParticipants returns all users that the current user has chatted with
// This is useful for showing a list of conversations
func getChatParticipants(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getChatParticipants was called")

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

	// Get chat participants
	participants, err := messagesRepo.GetChatParticipants(currentUserId)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching chat participants: "+err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, participants)
}

// getChatMessages retrieves all messages between the current user and another user
// Automatically marks messages from the other user as read
func getChatMessages(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getChatMessages was called")

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

	// Get other user ID from query parameter
	otherUserIdStr := r.URL.Query().Get("userId")
	if otherUserIdStr == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "userId parameter is required")
		return
	}

	otherUserId, err := strconv.Atoi(otherUserIdStr)
	if err != nil || otherUserId <= 0 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid userId parameter")
		return
	}

	// Cannot get chat with yourself
	if otherUserId == currentUserId {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Cannot get chat with yourself")
		return
	}

	// Get chat messages (this also marks messages as read)
	messages, err := messagesRepo.GetChatMessages(currentUserId, otherUserId)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching chat messages: "+err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, messages)
}

// getUnreadCount returns the total unread message count for the current user
func getUnreadCount(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getUnreadCount was called")

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

	// Get unread count
	count, err := messagesRepo.GetUnreadMessageCount(currentUserId)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching unread count: "+err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, map[string]int{"unreadCount": count})
}

// getAvailableChatUsers returns available users to chat with
// - For trainers: returns all gym members
// - For members: returns all gym trainers
func getAvailableChatUsers(w http.ResponseWriter, r *http.Request) {
	LogService.LogMessage("getAvailableChatUsers was called")

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
	roleId := claims.RoleId

	// Get available chat users based on role
	users, err := messagesRepo.GetAvailableChatUsers(currentUserId, roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return empty list if no users found
			utils.SendSuccessResponse(w, http.StatusOK, []models.AvailableChatUser{})
			return
		}
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Error fetching available chat users: "+err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, users)
}

