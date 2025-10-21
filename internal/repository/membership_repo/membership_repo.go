package membershipRepo

import (
	"database/sql"
	"zymm/internal/db"
	"zymm/internal/models"
	LogService "zymm/internal/service/log_service"
)

// InsertPlan inserts a plan row and returns the newly created id
func InsertPlan(p models.PlanRecord) (*int, error) {
	var id int
	// Use OUTPUT INSERTED.planId to get the inserted id
	err := db.DB.QueryRow(`
        INSERT INTO plans (planBanner, planName, planDesc, planPrice, planDuration, isActive, createdBy, gymId)
        OUTPUT INSERTED.planId
        VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8)`,
		p.PlanBanner, p.PlanName, p.PlanDesc, p.PlanPrice, p.PlanDuration, p.IsActive, p.CreatedBy, p.GymId,
	).Scan(&id)
	if err != nil {
		LogService.LogError("❌ DB error inserting plan: ", err)
		return nil, err
	}
	return &id, nil
}

// InsertPlanChangeLog inserts a row into planChangesLogs
func InsertPlanChangeLog(log models.PlanChangeLog) error {
	_, err := db.DB.Exec(`
        INSERT INTO planChangesLogs (planId, changedBy, changeType, changeDetails, oldPlanBanner, oldPlanName, oldPlanDesc, oldPlanPrice, oldPlanDuration, oldIsActive, newPlanBanner, newPlanName, newPlanDesc, newPlanPrice, newPlanDuration, newIsActive)
        VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15, @p16)`,
		log.PlanId, log.ChangedBy, log.ChangeType, log.ChangeDetails, log.OldPlanBanner, log.OldPlanName, log.OldPlanDesc, log.OldPlanPrice, log.OldPlanDuration, log.OldIsActive, log.NewPlanBanner, log.NewPlanName, log.NewPlanDesc, log.NewPlanPrice, log.NewPlanDuration, log.NewIsActive,
	)
	if err != nil {
		LogService.LogError("❌ DB error inserting plan change log: ", err)
		return err
	}
	return nil
}

func GetPlanById(planId int) (*models.PlanRecord, error) {
	plan := &models.PlanRecord{}
	err := db.DB.QueryRow(`
		SELECT planId, planBanner, planName, planDesc, planPrice, planDuration, isActive, createdBy, createdAt, gymId
		FROM plans
		WHERE planId = @p1`,
		planId).Scan(
		&plan.PlanId, &plan.PlanBanner, &plan.PlanName, &plan.PlanDesc, &plan.PlanPrice, &plan.PlanDuration, &plan.IsActive, &plan.CreatedBy, &plan.CreatedAt, &plan.GymId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		LogService.LogError("❌ DB error: ", err)
		return nil, err
	}
	return plan, nil
}

func GetUserMembershipByUserId(userId int) (*models.UserMembership, error) {
	membership := &models.UserMembership{}
	err := db.DB.QueryRow(`
		SELECT top 1 membershipId, userId, planId, startDate, endDate, isActive, membershipStatus, createdAt
		FROM userMemberships
		WHERE userId = @p1
		and isActive = 1
		order by createdAt desc`,
		userId).Scan(
		&membership.MembershipId, &membership.UserId, &membership.PlanId, &membership.StartDate, &membership.EndDate, &membership.IsActive, &membership.MembershipStatus, &membership.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		LogService.LogError("❌ DB error: ", err)
		return nil, err
	}
	return membership, nil
}

func InsertMembershipRequest(um models.UserMembership) (*int, error) {
	var id int
	// Use OUTPUT INSERTED.planId to get the inserted id
	err := db.DB.QueryRow(`
        INSERT INTO userMemberships (userId, planId, startDate, endDate, membershipStatus)
        OUTPUT INSERTED.membershipId
        VALUES (@p1, @p2, @p3, @p4, @p5)`,
		um.UserId, um.PlanId, um.StartDate, um.EndDate, um.MembershipStatus,
	).Scan(&id)
	if err != nil {
		LogService.LogError("❌ DB error inserting userMembership: ", err)
		return nil, err
	}
	return &id, nil
}
