package businessMembership

import (
	"database/sql"
	"errors"
	"strings"

	"zymm/internal/models"
	authRepo "zymm/internal/repository/auth_repo"
)

func ValidateCreatePlanRequest(req models.CreatePlanRequest) error {
	missing := []string{}

	if req.PlanName == "" {
		missing = append(missing, "planName")
	}
	if req.PlanDesc == "" {
		missing = append(missing, "planDesc")
	}
	if req.PlanPrice < 0 {
		missing = append(missing, "planPrice (cannot be negative)")
	}
	if req.PlanDuration <= 0 {
		missing = append(missing, "planDuration (must be greater than 0)")
	}
	if req.GymId <= 0 {
		missing = append(missing, "gymId (must be greater than 0)")
	}
	if req.UserAgent == "" {
		missing = append(missing, "userAgent")
	}
	if req.AppVersion == "" {
		missing = append(missing, "appVersion")
	}

	if len(missing) > 0 {
		return errors.New("Missing Required Fields: " + strings.Join(missing, ", "))
	}

	return nil
}

func ValidateCreatePlanWithDBChecks(req models.CreatePlanRequest, userId int) error {
	// First validate the basic request fields
	if err := ValidateCreatePlanRequest(req); err != nil {
		return err
	}

	// Check if user exists
	_, err := authRepo.GetUserByUserId(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("User does not exist")
		}
		return errors.New("Error validating user: " + err.Error())
	}

	// Check if gym exists
	_, err = authRepo.GetGymById(req.GymId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Gym does not exist")
		}
		return errors.New("Error validating gym: " + err.Error())
	}

	return nil
}
