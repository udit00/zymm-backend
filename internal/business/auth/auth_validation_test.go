package bussinessAuth

import (
    "strings"
    "testing"
    "zymm/internal/models"
)

func TestValidateRegistrationRequest(t *testing.T) {
    validReq := models.RegistrationApiRequestModel{
        DisplayName: "John Doe",
        Mobile:      "1234567890",
        Password:    "strongpass",
        Gender:      "M",
        AppVersion:  "1.0.0",
        UserAgent:   "ios",
        LocationLat: "12.34",
        LocationLong:"56.78",
        IpAddress:   "127.0.0.1",
    }

    if err := ValidateRegistrationRequest(validReq); err != nil {
        t.Fatalf("expected no error for valid request, got %v", err)
    }

    invalidReq := models.RegistrationApiRequestModel{}
    err := ValidateRegistrationRequest(invalidReq)
    if err == nil {
        t.Fatalf("expected error for missing fields, got nil")
    }
    missingFields := []string{"displayName", "mobile", "password", "gender", "appVersion", "userAgent", "locationLat", "locationLong", "ipAddress"}
    for _, field := range missingFields {
        if !strings.Contains(err.Error(), field) {
            t.Errorf("expected error to mention %s", field)
        }
    }

    shortPass := validReq
    shortPass.Password = "123"
    if err := ValidateRegistrationRequest(shortPass); err == nil || !strings.Contains(err.Error(), "Password must be at least 6 characters") {
        t.Fatalf("expected password length error, got %v", err)
    }
}

func TestValidateEmployeeRegistrationRequest(t *testing.T) {
    validReq := models.EmployeeRegistrationApiRequestModel{
        DisplayName: "Jane",
        Mobile:      "1234567890",
        Password:    "strongpass",
        Gender:      "F",
        AppVersion:  "1.0.0",
        UserAgent:   "android",
        LocationLat: "12.34",
        LocationLong:"56.78",
        IpAddress:   "127.0.0.1",
        RoleId:      3,
    }

    if err := ValidateEmployeeRegistrationRequest(validReq); err != nil {
        t.Fatalf("expected no error for valid request, got %v", err)
    }

    invalidRole := validReq
    invalidRole.RoleId = 0
    err := ValidateEmployeeRegistrationRequest(invalidRole)
    if err == nil || !strings.Contains(err.Error(), "roleId") {
        t.Fatalf("expected error for invalid roleId, got %v", err)
    }

    shortPass := validReq
    shortPass.Password = "123"
    if err := ValidateEmployeeRegistrationRequest(shortPass); err == nil || !strings.Contains(err.Error(), "Password must be at least 6 characters") {
        t.Fatalf("expected password length error, got %v", err)
    }
}

func TestValidateOwnerRegistrationRequest(t *testing.T) {
    validReq := models.RegistrationOwnerApiRequestModel{
        DisplayName:  "Owner",
        Mobile:       "1234567890",
        Password:     "strongpass",
        Gender:       "M",
        AppVersion:   "1.0.0",
        UserAgent:    "web",
        IpAddress:    "127.0.0.1",
        GymName:      "Gym",
        State:        "State",
        City:         "City",
        GymAddress:   "Address",
        ContactNo:    "9876543210",
        OfficialEmail:"gym@example.com",
        LocationLat:  "12.34",
        LocationLong: "56.78",
    }

    if err := ValidateOwnerRegistrationRequest(validReq); err != nil {
        t.Fatalf("expected no error for valid request, got %v", err)
    }

    invalidReq := models.RegistrationOwnerApiRequestModel{}
    err := ValidateOwnerRegistrationRequest(invalidReq)
    if err == nil {
        t.Fatalf("expected error for missing fields, got nil")
    }
    requiredFields := []string{"displayName", "mobile", "password", "gender", "appVersion", "userAgent", "ipAddress", "gymName", "state", "city", "gymAddress", "gymOfficialContactNo", "gymOfficialEmail", "gymOfficialLocationLat", "gymOfficialLocationLong"}
    for _, field := range requiredFields {
        if !strings.Contains(err.Error(), field) {
            t.Errorf("expected error to mention %s", field)
        }
    }

    shortPass := validReq
    shortPass.Password = "123"
    if err := ValidateOwnerRegistrationRequest(shortPass); err == nil || !strings.Contains(err.Error(), "Password must be at least 6 characters") {
        t.Fatalf("expected password length error, got %v", err)
    }
}

