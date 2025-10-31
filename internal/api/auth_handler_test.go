package api

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    bussinessAuth "zymm/internal/business/auth"
    "zymm/internal/models"
)

func setupJWT(t *testing.T) func() {
    t.Helper()
    oldSecret := os.Getenv("ZYMM_BACKEND_AUTH_SECRET_KEY")
    if err := os.Setenv("ZYMM_BACKEND_AUTH_SECRET_KEY", "unit-test-secret"); err != nil {
        t.Fatalf("failed to set env: %v", err)
    }
    bussinessAuth.MY_SECRET_KEY = ""
    return func() {
        if oldSecret == "" {
            os.Unsetenv("ZYMM_BACKEND_AUTH_SECRET_KEY")
        } else {
            os.Setenv("ZYMM_BACKEND_AUTH_SECRET_KEY", oldSecret)
        }
        bussinessAuth.MY_SECRET_KEY = ""
    }
}

func TestGetAuthTokenDataSuccess(t *testing.T) {
    cleanup := setupJWT(t)
    defer cleanup()

    token, err := bussinessAuth.GenerateJWTToken(7, 2)
    if err != nil {
        t.Fatalf("failed to generate token: %v", err)
    }

    req := httptest.NewRequest(http.MethodGet, "/", nil)
    req.Header.Set("Authorization", "Bearer "+*token)
    rr := httptest.NewRecorder()

    getAuthTokenData(rr, req)

    if rr.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", rr.Code)
    }

    var resp models.APIResponse
    if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
        t.Fatalf("unable to parse response: %v", err)
    }
    claims, ok := resp.Data.(map[string]interface{})
    if !ok {
        t.Fatalf("expected claims map, got %T", resp.Data)
    }
    if int(claims["userId"].(float64)) != 7 {
        t.Fatalf("expected userId to be 7, got %v", claims["userId"])
    }
}

func TestGetAuthTokenDataMissingHeader(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    rr := httptest.NewRecorder()

    getAuthTokenData(rr, req)

    if rr.Code != http.StatusUnauthorized {
        t.Fatalf("expected status 401, got %d", rr.Code)
    }
}

func TestAuthMiddleware(t *testing.T) {
    cleanup := setupJWT(t)
    defer cleanup()

    token, err := bussinessAuth.GenerateJWTToken(99, 4)
    if err != nil {
        t.Fatalf("failed to generate token: %v", err)
    }

    called := false
    handler := AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
        called = true
        claims, ok := r.Context().Value(ctxClaimDataKey).(*bussinessAuth.MyCustomClaims)
        if !ok || claims == nil {
            t.Fatalf("expected claims in context")
        }
        if claims.UserId != 99 || claims.RoleId != 4 {
            t.Fatalf("unexpected claims data: %+v", claims)
        }
    })

    req := httptest.NewRequest(http.MethodGet, "/", nil)
    req.Header.Set("Authorization", "Bearer "+*token)
    rr := httptest.NewRecorder()

    handler(rr, req)

    if !called {
        t.Fatalf("expected wrapped handler to be invoked")
    }
    if rr.Code != 0 && rr.Code != http.StatusOK {
        t.Fatalf("expected no error response, got %d", rr.Code)
    }
}

func TestAuthMiddlewareUnauthorized(t *testing.T) {
    handler := AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {})

    req := httptest.NewRequest(http.MethodGet, "/", nil)
    rr := httptest.NewRecorder()

    handler(rr, req)

    if rr.Code != http.StatusUnauthorized {
        t.Fatalf("expected status 401, got %d", rr.Code)
    }
}

