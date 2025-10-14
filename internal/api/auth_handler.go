package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"zymm/internal/business/auth"
	"zymm/internal/db"
	"zymm/internal/models"
	"zymm/utils"
)

const authApiVersion = "v1"
const authApiPrefix = "auth"

func authRouteAppended(newRoute string) string {
	return utils.ApiRoute(authApiVersion, authApiPrefix, newRoute)
}

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func GetCommonErrorResponse(status int, errMsg string) ErrorResponse {
	return ErrorResponse{
		Status: status,
		Error:  errMsg,
	}
}

func SendErrorResponse(writer http.ResponseWriter, status int, errMsg string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(GetCommonErrorResponse(status, errMsg))
}

func AuthHandlerDelegate(mux *http.ServeMux) {
	mux.HandleFunc(authRouteAppended("login"), loginHandler)
	mux.HandleFunc(authRouteAppended("register"), registerHandler)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("loginHandler was called")
	if r.Method != http.MethodPost {
		SendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Decode JSON body
	var req models.LoginRequestModel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Check user in DB
	var dbPassword string
	err := db.DB.QueryRow("SELECT userPass FROM users WHERE email = @p1 or mobile = @p2", req.EmailOrMobile, req.EmailOrMobile).Scan(&dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			SendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials: no user found with email or mobile "+req.EmailOrMobile)
			return
		}
		log.Printf("❌ DB error: %v", err)
		SendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	isPassSame, passMatchError := auth.ComparePasswordArgon2id(dbPassword, req.Password)
	if passMatchError != nil {
		SendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials: pass: "+req.Password+" was wrong, correct Pass is "+dbPassword)
		return
	} else if !isPassSame {
		SendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials: pass: "+req.Password+" was wrong, correct Pass is "+dbPassword)
		return
	}
	// Success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "✅ Login successful for " + req.EmailOrMobile})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	pass, pErr := auth.HashPasswordArgon2id("Password@123*")
	if pErr != nil {
		log.Println("Error generating password hash: ", pErr)
	} else {
		log.Println("Password will be - " + pass)
	}
	// Success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "✅ Registration successful for " + pass})
}
