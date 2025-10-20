package api

import (
    "context"
    "net/http"
    "strings"

    bussinessAuth "zymm/internal/business/auth"
    LogService "zymm/internal/service/log_service"
    "zymm/utils"
)

type ctxKey string

const ctxUserIDKey ctxKey = "userId"

// AuthMiddleware validates Authorization: Bearer <token>, extracts userId from JWT
// and attaches it to request context under ctxUserIDKey. It uses the project's
// business auth helpers for parsing/validation so it follows existing claim names.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        LogService.LogMessage("AuthMiddleware: checking Authorization header")

        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            utils.SendErrorResponse(w, http.StatusUnauthorized, "Missing Authorization header")
            return
        }
        if !strings.HasPrefix(authHeader, "Bearer ") {
            utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid Authorization header")
            return
        }

        token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
        claims, err := bussinessAuth.GetDataFromJWTToken(token)
        if err != nil {
            LogService.LogError("AuthMiddleware: token validation failed: ", err)
            utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid or expired auth token")
            return
        }

        // attach user id to context and call next
        ctx := context.WithValue(r.Context(), ctxUserIDKey, claims.UserId)
        next(w, r.WithContext(ctx))
    }
}
