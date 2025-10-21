package membershipRepo

import (
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
