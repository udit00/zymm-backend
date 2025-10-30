package utils

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    "zymm/internal/config"
    "zymm/internal/models"
)

func setupConfig(t *testing.T) func() {
    t.Helper()
    oldDebug := os.Getenv("IS_DEBUG")
    oldPort := os.Getenv("ZYMM_PORT")
    if err := os.Setenv("IS_DEBUG", "true"); err != nil {
        t.Fatalf("setenv: %v", err)
    }
    if err := os.Setenv("ZYMM_PORT", "5000"); err != nil {
        t.Fatalf("setenv: %v", err)
    }
    config.Init()
    return func() {
        if oldDebug == "" {
            os.Unsetenv("IS_DEBUG")
        } else {
            os.Setenv("IS_DEBUG", oldDebug)
        }
        if oldPort == "" {
            os.Unsetenv("ZYMM_PORT")
        } else {
            os.Setenv("ZYMM_PORT", oldPort)
        }
    }
}

func TestApiRoute(t *testing.T) {
    cleanup := setupConfig(t)
    defer cleanup()

    got := ApiRoute("v1", "auth", "login")
    want := "/" + config.GetAppName() + "/v1/auth/login"
    if got != want {
        t.Fatalf("expected %s, got %s", want, got)
    }
}

func TestSendErrorResponse(t *testing.T) {
    rr := httptest.NewRecorder()
    SendErrorResponse(rr, http.StatusBadRequest, "bad request")

    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
    }

    var resp models.APIResponse
    if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
        t.Fatalf("failed to unmarshal response: %v", err)
    }

    if resp.Status != http.StatusBadRequest || resp.Error != "bad request" {
        t.Fatalf("unexpected response body: %+v", resp)
    }
}

func TestSendSuccessResponse(t *testing.T) {
    rr := httptest.NewRecorder()
    data := map[string]string{"message": "ok"}
    SendSuccessResponse(rr, http.StatusOK, data)

    if rr.Code != http.StatusOK {
        t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
    }

    var resp models.APIResponse
    if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
        t.Fatalf("failed to unmarshal response: %v", err)
    }

    if resp.Status != http.StatusOK {
        t.Fatalf("unexpected status in response: %+v", resp)
    }
    if resp.Error != "" {
        t.Fatalf("expected empty error, got %s", resp.Error)
    }
    if resp.Data == nil {
        t.Fatalf("expected data to be present")
    }
}

