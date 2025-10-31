package businessMembership

import (
    "strings"
    "testing"
    "time"
    "zymm/internal/models"
)

func TestValidateUpsertPlanRequest(t *testing.T) {
    planId := 0
    valid := models.UpsertPlanRequest{
        PlanId:      &planId,
        PlanName:    "Pro",
        PlanDesc:    "Desc",
        PlanPrice:   1000,
        PlanDuration:30,
        GymId:       10,
        UserAgent:   "ios",
        AppVersion:  "1.0.0",
    }

    if err := ValidateUpsertPlanRequest(valid); err != nil {
        t.Fatalf("expected valid request, got %v", err)
    }

    invalid := models.UpsertPlanRequest{}
    err := ValidateUpsertPlanRequest(invalid)
    if err == nil {
        t.Fatalf("expected error for missing fields")
    }

    required := []string{"planId", "planName", "planDesc", "planPrice", "planDuration", "gymId", "userAgent", "appVersion"}
    for _, field := range required {
        if !strings.Contains(err.Error(), field) {
            t.Errorf("expected error to mention %s", field)
        }
    }

    negativePrice := valid
    negativePrice.PlanPrice = -1
    if err := ValidateUpsertPlanRequest(negativePrice); err == nil || !strings.Contains(err.Error(), "planPrice") {
        t.Fatalf("expected plan price validation error, got %v", err)
    }
}

func TestGetPlanChangeDiffDetails(t *testing.T) {
    now := time.Now().UTC()
    later := now.Add(time.Hour)
    plan1 := models.PlanRecord{
        PlanName:     "Plan",
        PlanDesc:     "Desc",
        PlanPrice:    100,
        PlanDuration: 30,
        PlanBanner:   stringPtr("banner1"),
        IsActive:     boolPtr(true),
        CreatedBy:    intPtr(1),
        CreatedAt:    &now,
        GymId:        intPtr(10),
    }

    plan2 := models.PlanRecord{
        PlanName:     "Plan Updated",
        PlanDesc:     "Desc",
        PlanPrice:    150,
        PlanDuration: 60,
        PlanBanner:   stringPtr("banner2"),
        IsActive:     boolPtr(false),
        CreatedBy:    intPtr(2),
        CreatedAt:    &later,
        GymId:        intPtr(11),
    }

    diff := GetPlanChangeDiffDetails(plan1, plan2)
    expectedFields := []string{"PlanName", "PlanBanner", "PlanPrice", "PlanDuration", "IsActive", "CreatedBy", "CreatedAt", "GymId"}
    for _, field := range expectedFields {
        if !strings.Contains(diff, field) {
            t.Errorf("expected diff to list %s", field)
        }
    }

    if diff == "" {
        t.Fatalf("expected non-empty diff string")
    }

    // identical plans should return empty diff
    if diff2 := GetPlanChangeDiffDetails(plan1, plan1); diff2 != "" {
        t.Fatalf("expected empty diff for identical plans, got %s", diff2)
    }
}

func stringPtr(s string) *string { return &s }
func boolPtr(b bool) *bool       { return &b }
func intPtr(i int) *int          { return &i }

