package membershipRepo

import (
	"database/sql"
	bussinessMembershipRequest "zymm/internal/business/membership/membership_request_action_type"
	"zymm/internal/db"
	"zymm/internal/models"
	LogService "zymm/internal/service/log_service"
)

type MembershipStatus int

const (
	Pending MembershipStatus = iota
	Rejected
	Approved
	All
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

func UpdatePlan(newPlanDetails models.PlanRecord) error {
	_, err := db.DB.Exec(`
		UPDATE plans SET 
		planBanner = @p1, 
		planName = @p2,
		planDesc = @p3,
		planPrice = @p4,
		planDuration = @p5,
		isActive = @p6,		
		gymId	= @p7
		WHERE planId = @p8`,
		newPlanDetails.PlanBanner,
		newPlanDetails.PlanName,
		newPlanDetails.PlanDesc,
		newPlanDetails.PlanPrice,
		newPlanDetails.PlanDuration,
		newPlanDetails.IsActive,
		newPlanDetails.GymId,
		newPlanDetails.PlanId)
	if err != nil {
		LogService.LogError("❌ DB error: ", err)
		return err
	}
	return nil
}

// InsertPlanChangeLog inserts a row into planChangesLogs
func InsertPlanChangeLog(log models.PlanChangeLog) (*int, error) {
	var id int
	err := db.DB.QueryRow(`
        INSERT INTO planChangesLogs (planId, changedBy, changeType, changeDetails, oldPlanBanner, oldPlanName, oldPlanDesc, oldPlanPrice, oldPlanDuration, oldIsActive, newPlanBanner, newPlanName, newPlanDesc, newPlanPrice, newPlanDuration, newIsActive)
		OUTPUT INSERTED.logId
        VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10, @p11, @p12, @p13, @p14, @p15, @p16)`,
		log.PlanId, log.ChangedBy, log.ChangeType, log.ChangeDetails, log.OldPlanBanner, log.OldPlanName, log.OldPlanDesc, log.OldPlanPrice, log.OldPlanDuration, log.OldIsActive, log.NewPlanBanner, log.NewPlanName, log.NewPlanDesc, log.NewPlanPrice, log.NewPlanDuration, log.NewIsActive,
	).Scan(&id)
	if err != nil {
		LogService.LogError("❌ DB error inserting plan change log: ", err)
		return nil, err
	}
	return &id, nil
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
		and membershipStatus in ('A', 'P')
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

func GetUserMembershipByMembershipId(membershipId int) (*models.UserMembership, error) {
	membership := &models.UserMembership{}
	err := db.DB.QueryRow(`
		SELECT membershipId, userId, planId, startDate, endDate, isActive, membershipStatus, createdAt
		FROM userMemberships
		WHERE membershipId = @p1`,
		membershipId).Scan(
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

func GetAllMembershipPlansRequestByUserId(userId int, filterBy MembershipStatus) ([]models.UserMembership, error) {
	var query string

	if filterBy == All {
		query = `SELECT membershipId, userId, planId, startDate, endDate, isActive, membershipStatus, createdAt
			FROM userMemberships
			WHERE userId = @p1
			AND isActive = 1
			ORDER BY createdAt DESC`
	} else {
		var filterByTypeChar string
		switch filterBy {
		case Pending:
			{
				filterByTypeChar = "P"
			}
		case Rejected:
			{
				filterByTypeChar = "R"
			}
		case Approved:
			{
				filterByTypeChar = "A"
			}
		}
		query = `SELECT membershipId, userId, planId, startDate, endDate, isActive, membershipStatus, createdAt
			FROM userMemberships
			WHERE userId = @p1
			AND isActive = 1
			and membershipStatus = ` + filterByTypeChar + ` 
			ORDER BY createdAt DESC`
	}
	rows, err := db.DB.Query(query, userId)
	if err != nil {
		LogService.LogError("❌ DB query error: ", err)
		return nil, err
	}
	defer rows.Close()

	var memberships []models.UserMembership

	for rows.Next() {
		var membership models.UserMembership
		err := rows.Scan(
			&membership.MembershipId,
			&membership.UserId,
			&membership.PlanId,
			&membership.StartDate,
			&membership.EndDate,
			&membership.IsActive,
			&membership.MembershipStatus,
			&membership.CreatedAt,
		)
		if err != nil {
			LogService.LogError("❌ DB scan error: ", err)
			return nil, err
		}
		memberships = append(memberships, membership)
	}

	if err = rows.Err(); err != nil {
		LogService.LogError("❌ DB rows error: ", err)
		return nil, err
	}

	if len(memberships) == 0 {
		return nil, sql.ErrNoRows
	}

	return memberships, nil
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

func TakeActionOnMembershipRequest(memberId int, membershipId int, actionTaken bussinessMembershipRequest.ActionType, userId int) (*int, error) {
	var id int
	// possible values "A" or "R"
	var actionTakenString string = "A"
	if actionTaken == bussinessMembershipRequest.Reject {
		actionTakenString = "R"
	}
	err := db.DB.QueryRow(`
		INSERT INTO userMembershipsChangesLogs (memberId,  membershipId, changesMade, changedBy, loggedAt)
		OUTPUT INSERTED.logId
		VALUES (@p1, @p2, @p3, @p4, getDate())`,
		memberId, membershipId, actionTakenString, userId,
	).Scan(&id)
	if err != nil {
		LogService.LogError("❌ DB error inserting userMembershipsChangesLogs: ", err)
		return nil, err
	}
	return &id, nil
}

func UpdateActionTakenOnUserMembership(membershipId int, actionTaken bussinessMembershipRequest.ActionType) error {
	var actionTakenString string = "A"
	if actionTaken == bussinessMembershipRequest.Reject {
		actionTakenString = "R"
	}
	_, err := db.DB.Exec(`
		UPDATE userMemberships SET membershipStatus = @p1
		WHERE membershipId = @p2`,
		actionTakenString,
		membershipId,
	)
	if err != nil {
		LogService.LogError("❌ DB error: ", err)
		return err
	}
	return nil
}

func GetPlanCountByGymId(gymId int) (*int, error) {
	count := 0
	err := db.DB.QueryRow(`
		SELECT count(*)
		FROM plans
		WHERE gymId = @p1
		and isActive = 1
		`,
		gymId).Scan(
		&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		LogService.LogError("❌ DB error: ", err)
		return nil, err
	}
	return &count, nil
}
