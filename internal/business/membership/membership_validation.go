package businessMembership

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	businessRoles "zymm/internal/business/roles"
	rolestype "zymm/internal/business/roles/roles_type"
	"zymm/internal/models"
	gymRepo "zymm/internal/repository/gym_repo"
	membershipRepo "zymm/internal/repository/membership_repo"
	"zymm/internal/repository/userRepo"
)

func ValidateUpsertPlanRequest(req models.UpsertPlanRequest) error {
	missing := []string{}
	var checkPlanId int = 0
	if req.PlanId == nil {
		checkPlanId = -1
	} else {
		checkPlanId = *req.PlanId
	}
	if checkPlanId < 0 {
		missing = append(missing, "planId")
	}
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

func ValidateUpsertPlanWithDBChecks(req models.UpsertPlanRequest, userId int) error {
	// First validate the basic request fields
	if err := ValidateUpsertPlanRequest(req); err != nil {
		return err
	}

	// Check if user exists
	userDetails, err := userRepo.GetUserByUserId(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("User does not exist")
		}
		return errors.New("Error validating user: " + err.Error())
	}

	// 1 is owner and 2 is manager, anyone else shouldn't edit or create plans
	roleType := rolestype.GetRoleTypeFromInt(userDetails.RoleId)
	if roleType == nil {
		return errors.New("Could not resolve your role, please login again")
	}

	if businessRoles.IsNotAllowedToManagePlans(*roleType) {
		return errors.New("Only Owner or Manager of the gym can create or edit plans.")
	}

	// Check if gym exists
	_, err = gymRepo.GetGymById(req.GymId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Gym does not exist")
		}
		return errors.New("Error validating gym: " + err.Error())
	}

	if *req.PlanId > 0 {
		planDetails, planDetailsErr := membershipRepo.GetPlanById(*req.PlanId)
		if planDetailsErr != nil {
			if planDetailsErr == sql.ErrNoRows {
				return errors.New("Plan does not exist")
			}
			return errors.New("Error validating Plan: " + planDetailsErr.Error())
		}

		if *planDetails.GymId != req.GymId {
			return errors.New("Gym does not match in the plan and your gym.")
		}
	}

	return nil
}

func GetPlanChangeDiffDetails(p1 models.PlanRecord, p2 models.PlanRecord) string {
	differences := []string{}

	// Helper to compare pointer values safely
	comparePtr := func(a, b interface{}) bool {
		switch v1 := a.(type) {
		case *string:
			v2 := b.(*string)
			if v1 == nil && v2 == nil {
				return false
			}
			if v1 == nil || v2 == nil {
				return true
			}
			return *v1 != *v2

		case *bool:
			v2 := b.(*bool)
			if v1 == nil && v2 == nil {
				return false
			}
			if v1 == nil || v2 == nil {
				return true
			}
			return *v1 != *v2

		case *int:
			v2 := b.(*int)
			if v1 == nil && v2 == nil {
				return false
			}
			if v1 == nil || v2 == nil {
				return true
			}
			return *v1 != *v2

		case *time.Time:
			v2 := b.(*time.Time)
			if v1 == nil && v2 == nil {
				return false
			}
			if v1 == nil || v2 == nil {
				return true
			}
			return !v1.Equal(*v2)
		}
		return false
	}

	// Compare normal (non-pointer) fields
	if p1.PlanName != p2.PlanName {
		differences = append(differences, "PlanName")
	}
	if comparePtr(p1.PlanBanner, p2.PlanBanner) {
		differences = append(differences, "PlanBanner")
	}
	if p1.PlanDesc != p2.PlanDesc {
		differences = append(differences, "PlanDesc")
	}
	if p1.PlanPrice != p2.PlanPrice {
		differences = append(differences, "PlanPrice")
	}
	if p1.PlanDuration != p2.PlanDuration {
		differences = append(differences, "PlanDuration")
	}
	if comparePtr(p1.IsActive, p2.IsActive) {
		differences = append(differences, "IsActive")
	}
	if comparePtr(p1.CreatedBy, p2.CreatedBy) {
		differences = append(differences, "CreatedBy")
	}
	if comparePtr(p1.CreatedAt, p2.CreatedAt) {
		differences = append(differences, "CreatedAt")
	}
	if comparePtr(p1.GymId, p2.GymId) {
		differences = append(differences, "GymId")
	}

	// Return result
	if len(differences) > 0 {
		return strings.Join(differences, ", ")
	}
	return ""
}
