package gymRepo

import (
	"database/sql"
	"zymm/internal/db"
	"zymm/internal/models"
	gymModels "zymm/internal/models/gym_models"
	LogService "zymm/internal/service/log_service"
)

func GetActiveUserIdsForGym(gymId int) ([]int, error) {
	rows, err := db.DB.Query(`
		select distinct um.userId
		from userMemberships um
		inner join plans p on p.planId = um.planId
		where p.isActive = 1
		and um.isActive = 1
		and um.startDate <= getDate()
		and um.endDate >= getDate()
		and p.gymId = @p1`, gymId)
	if err != nil {
		LogService.LogError("❌ DB query error fetching active user ids: ", err)
		return nil, err
	}
	defer rows.Close()

	var userIds []int
	for rows.Next() {
		var id int
		if scanErr := rows.Scan(&id); scanErr != nil {
			LogService.LogError("❌ DB scan error fetching active user ids: ", scanErr)
			return nil, scanErr
		}
		userIds = append(userIds, id)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		LogService.LogError("❌ DB rows iteration error fetching active user ids: ", rowsErr)
		return nil, rowsErr
	}

	return userIds, nil
}

func GetActiveMembershipDetailsForGym(gymId int) ([]models.ActiveMembershipUser, error) {
	rows, err := db.DB.Query(`
		select distinct um.userId,
			   um.membershipId,
			   um.startDate,
			   um.endDate,
			   um.planId,
			   um.membershipStatus,
			   u.userName,
			   u.email,
			   u.mobile,
			   p.planName,
			   p.planDuration,
			   p.planPrice
		from userMemberships um
		inner join plans p on p.planId = um.planId
		inner join users u on u.userId = um.userId
		where p.isActive = 1
		and um.isActive = 1
		and um.startDate <= getDate()
		and um.endDate >= getDate()
		and p.gymId = @p1`, gymId)
	if err != nil {
		LogService.LogError("❌ DB query error fetching active membership details: ", err)
		return nil, err
	}
	defer rows.Close()

	var members []models.ActiveMembershipUser
	for rows.Next() {
		var member models.ActiveMembershipUser
		if scanErr := rows.Scan(
			&member.UserId,
			&member.MembershipId,
			&member.StartDate,
			&member.EndDate,
			&member.PlanId,
			&member.MembershipStatus,
			&member.UserName,
			&member.Email,
			&member.Mobile,
			&member.PlanName,
			&member.PlanDuration,
			&member.PlanPrice,
		); scanErr != nil {
			LogService.LogError("❌ DB scan error fetching active membership details: ", scanErr)
			return nil, scanErr
		}

		members = append(members, member)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		LogService.LogError("❌ DB rows iteration error fetching active membership details: ", rowsErr)
		return nil, rowsErr
	}

	return members, nil
}

func GetAllGym() ([]gymModels.GymRecord, error) {
	query := `
		SELECT gymId, gymName, state, city, gymAddress, contactNo, officialEmail, createdBy, createdAt, updatedAt, locationLat, locationLong
		FROM gym
		`

	rows, err := db.DB.Query(query)
	if err != nil {
		LogService.LogError("❌ DB query error: ", err)
		return nil, err
	}
	defer rows.Close()

	var allGyms []gymModels.GymRecord

	for rows.Next() {
		var gym gymModels.GymRecord
		err := rows.Scan(
			&gym.GymId,
			&gym.GymName,
			&gym.State,
			&gym.City,
			&gym.GymAddress,
			&gym.ContactNo,
			&gym.OfficialEmail,
			&gym.CreatedBy,
			&gym.CreatedAt,
			&gym.UpdatedAt,
			&gym.LocationLat,
			&gym.LocationLong,
		)
		if err != nil {
			LogService.LogError("❌ DB scan error: ", err)
			return nil, err
		}
		allGyms = append(allGyms, gym)
	}

	if err = rows.Err(); err != nil {
		LogService.LogError("❌ DB rows error: ", err)
		return nil, err
	}

	if len(allGyms) == 0 {
		return []gymModels.GymRecord{}, nil
	}

	return allGyms, nil
}

func SearchGymsByName(searchQuery string) ([]gymModels.GymRecordWithAdditionalData, error) {
	query := `
		SELECT gymId, gymName, state, city, gymAddress, contactNo, officialEmail, createdBy, createdAt, updatedAt, locationLat, locationLong, 
		isnull(ratingData.averageRating, 0) as averageRating, isnull(ratingData.totalFeedbacks, 0) as totalFeedbacks, isnull(empCounts.trainersCount, 0) as trainersCount, 
		isnull(empCounts.staffCount, 0) as staffCount, isnull(empCounts.managersCount, 0) as managersCount, isnull(activePlans.cnt, 0) as activePlansCount
		FROM gym g
		outer apply (SELECT AVG(rating) AS averageRating, COUNT(*) AS totalFeedbacks FROM feedback where gymId = g.gymId) ratingData
		OUTER APPLY (
			SELECT 
				SUM(CASE WHEN u.roleId = 4 THEN 1 ELSE 0 END) AS trainersCount,
				SUM(CASE WHEN u.roleId = 3 THEN 1 ELSE 0 END) AS staffCount,
				SUM(CASE WHEN u.roleId = 2 THEN 1 ELSE 0 END) AS managersCount
			FROM employees e
			INNER JOIN users u ON u.userId = e.userId
			WHERE e.gymId = g.gymId
		) empCounts
		outer apply (Select count(*) as cnt from plans where gymId = g.gymId and isActive = 1) activePlans
		WHERE LOWER(g.gymName) LIKE '%' + LOWER(@p1) + '%' 
		   OR LOWER(g.city) LIKE '%' + LOWER(@p1) + '%'
		   OR LOWER(g.state) LIKE '%' + LOWER(@p1) + '%'
		ORDER BY g.gymName
		`

	rows, err := db.DB.Query(query, searchQuery)
	if err != nil {
		LogService.LogError("❌ DB query error: ", err)
		return nil, err
	}
	defer rows.Close()

	var gyms []gymModels.GymRecordWithAdditionalData

	for rows.Next() {
		var gym gymModels.GymRecordWithAdditionalData
		err := rows.Scan(
			&gym.GymId,
			&gym.GymName,
			&gym.State,
			&gym.City,
			&gym.GymAddress,
			&gym.ContactNo,
			&gym.OfficialEmail,
			&gym.CreatedBy,
			&gym.CreatedAt,
			&gym.UpdatedAt,
			&gym.LocationLat,
			&gym.LocationLong,
			&gym.AverageRating,
			&gym.TotalFeedback,
			&gym.TrainersCount,
			&gym.StaffCount,
			&gym.ManagersCount,
			&gym.ActivePlans,
		)
		if err != nil {
			LogService.LogError("❌ DB scan error: ", err)
			return nil, err
		}
		gyms = append(gyms, gym)
	}

	if err = rows.Err(); err != nil {
		LogService.LogError("❌ DB rows error: ", err)
		return nil, err
	}

	if len(gyms) == 0 {
		return []gymModels.GymRecordWithAdditionalData{}, nil
	}

	return gyms, nil
}

func GetAllGymWithAdditionalData() ([]gymModels.GymRecordWithAdditionalData, error) {
	query := `
		SELECT gymId, gymName, state, city, gymAddress, contactNo, officialEmail, createdBy, createdAt, updatedAt, locationLat, locationLong, 
		isnull(ratingData.averageRating, 0) as averageRating, isnull(ratingData.totalFeedbacks, 0) as totalFeedbacks, isnull(empCounts.trainersCount, 0) as trainersCount, 
		isnull(empCounts.staffCount, 0) as staffCount, isnull(empCounts.managersCount, 0) as managersCount, isnull(activePlans.cnt, 0) as activePlansCount
		FROM gym g
		outer apply (SELECT AVG(rating) AS averageRating, COUNT(*) AS totalFeedbacks FROM feedback where gymId = g.gymId) ratingData
		OUTER APPLY (
			SELECT 
				SUM(CASE WHEN u.roleId = 4 THEN 1 ELSE 0 END) AS trainersCount,
				SUM(CASE WHEN u.roleId = 3 THEN 1 ELSE 0 END) AS staffCount,
				SUM(CASE WHEN u.roleId = 2 THEN 1 ELSE 0 END) AS managersCount
			FROM employees e
			INNER JOIN users u ON u.userId = e.userId
			WHERE e.gymId = g.gymId
		) empCounts
		outer apply (Select count(*) as cnt from plans where gymId = g.gymId and isActive = 1) activePlans
		`

	rows, err := db.DB.Query(query)
	if err != nil {
		LogService.LogError("❌ DB query error: ", err)
		return nil, err
	}
	defer rows.Close()

	var allGyms []gymModels.GymRecordWithAdditionalData

	for rows.Next() {
		var gym gymModels.GymRecordWithAdditionalData
		err := rows.Scan(
			&gym.GymId,
			&gym.GymName,
			&gym.State,
			&gym.City,
			&gym.GymAddress,
			&gym.ContactNo,
			&gym.OfficialEmail,
			&gym.CreatedBy,
			&gym.CreatedAt,
			&gym.UpdatedAt,
			&gym.LocationLat,
			&gym.LocationLong,
			&gym.AverageRating,
			&gym.TotalFeedback,
			&gym.TrainersCount,
			&gym.StaffCount,
			&gym.ManagersCount,
			&gym.ActivePlans,
		)
		if err != nil {
			LogService.LogError("❌ DB scan error: ", err)
			return nil, err
		}
		allGyms = append(allGyms, gym)
	}

	if err = rows.Err(); err != nil {
		LogService.LogError("❌ DB rows error: ", err)
		return nil, err
	}

	if len(allGyms) == 0 {
		return []gymModels.GymRecordWithAdditionalData{}, nil
	}

	return allGyms, nil
}

func InsertGym(db *sql.DB, gym gymModels.GymRecord) (*int, error) {
	var gymId *int
	err := db.QueryRow(`
		INSERT INTO gym (gymName, state, city, gymAddress, contactNo, officialEmail, createdBy, locationLat, locationLong)
		OUTPUT INSERTED.gymId
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9)`,
		gym.GymName, gym.State, gym.City, gym.GymAddress, gym.ContactNo, gym.OfficialEmail, gym.CreatedBy, gym.LocationLat, gym.LocationLong,
	).Scan(&gymId)
	if err != nil {
		return nil, err
	}
	return gymId, nil
}

func GetGymById(gymId int) (*gymModels.GymRecord, error) {
	gym := &gymModels.GymRecord{}
	err := db.DB.QueryRow(`
		SELECT gymId, gymName, state, city, gymAddress, contactNo, officialEmail, createdBy, createdAt, updatedAt, locationLat, locationLong
		FROM gym
		WHERE gymId = @p1`,
		gymId).Scan(
		&gym.GymId, &gym.GymName, &gym.State, &gym.City, &gym.GymAddress, &gym.ContactNo, &gym.OfficialEmail, &gym.CreatedBy, &gym.CreatedAt, &gym.UpdatedAt, &gym.LocationLat, &gym.LocationLong)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		LogService.LogError("❌ DB error: ", err)
		return nil, err
	}
	return gym, nil
}

func GetGymWithAdditionalDataByGymId(gymId int) (*gymModels.GymRecordWithAdditionalData, error) {
	gym := &gymModels.GymRecordWithAdditionalData{}
	err := db.DB.QueryRow(`
		SELECT gymId, gymName, state, city, gymAddress, contactNo, officialEmail, createdBy, createdAt, updatedAt, locationLat, locationLong, 
		isnull(ratingData.averageRating, 0) as averageRating, isnull(ratingData.totalFeedbacks, 0) as totalFeedbacks, isnull(empCounts.trainersCount, 0) as trainersCount, 
		isnull(empCounts.staffCount, 0) as staffCount, isnull(empCounts.managersCount, 0) as managersCount, isnull(activePlans.cnt, 0) as activePlansCount
		FROM gym g
		outer apply (SELECT AVG(rating) AS averageRating, COUNT(*) AS totalFeedbacks FROM feedback where gymId = g.gymId) ratingData
		OUTER APPLY (
			SELECT 
				SUM(CASE WHEN u.roleId = 4 THEN 1 ELSE 0 END) AS trainersCount,
				SUM(CASE WHEN u.roleId = 3 THEN 1 ELSE 0 END) AS staffCount,
				SUM(CASE WHEN u.roleId = 2 THEN 1 ELSE 0 END) AS managersCount
			FROM employees e
			INNER JOIN users u ON u.userId = e.userId
			WHERE e.gymId = g.gymId
		) empCounts
		outer apply (Select count(*) as cnt from plans where gymId = g.gymId and isActive = 1) activePlans
		where g.gymId = @p1`, gymId).Scan(
		&gym.GymId,
		&gym.GymName,
		&gym.State,
		&gym.City,
		&gym.GymAddress,
		&gym.ContactNo,
		&gym.OfficialEmail,
		&gym.CreatedBy,
		&gym.CreatedAt,
		&gym.UpdatedAt,
		&gym.LocationLat,
		&gym.LocationLong,
		&gym.AverageRating,
		&gym.TotalFeedback,
		&gym.TrainersCount,
		&gym.StaffCount,
		&gym.ManagersCount,
		&gym.ActivePlans,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		LogService.LogError("❌ DB error: ", err)
		return nil, err
	}
	return gym, nil
}

func GetGymWithAdditionalDataByOwnerId(userId int) (*gymModels.GymRecordWithAdditionalData, error) {
	gym := &gymModels.GymRecordWithAdditionalData{}
	err := db.DB.QueryRow(`
		SELECT gymId, gymName, state, city, gymAddress, contactNo, officialEmail, createdBy, createdAt, updatedAt, locationLat, locationLong, 
		isnull(ratingData.averageRating, 0) as averageRating, isnull(ratingData.totalFeedbacks, 0) as totalFeedbacks, isnull(empCounts.trainersCount, 0) as trainersCount, 
		isnull(empCounts.staffCount, 0) as staffCount, isnull(empCounts.managersCount, 0) as managersCount, isnull(activePlans.cnt, 0) as activePlansCount
		FROM gym g
		outer apply (SELECT AVG(rating) AS averageRating, COUNT(*) AS totalFeedbacks FROM feedback where gymId = g.gymId) ratingData
		OUTER APPLY (
			SELECT 
				SUM(CASE WHEN u.roleId = 4 THEN 1 ELSE 0 END) AS trainersCount,
				SUM(CASE WHEN u.roleId = 3 THEN 1 ELSE 0 END) AS staffCount,
				SUM(CASE WHEN u.roleId = 2 THEN 1 ELSE 0 END) AS managersCount
			FROM employees e
			INNER JOIN users u ON u.userId = e.userId
			WHERE e.gymId = g.gymId
		) empCounts
		outer apply (Select count(*) as cnt from plans where gymId = g.gymId and isActive = 1) activePlans
		where g.createdBy = @p1`, userId).Scan(
		&gym.GymId,
		&gym.GymName,
		&gym.State,
		&gym.City,
		&gym.GymAddress,
		&gym.ContactNo,
		&gym.OfficialEmail,
		&gym.CreatedBy,
		&gym.CreatedAt,
		&gym.UpdatedAt,
		&gym.LocationLat,
		&gym.LocationLong,
		&gym.AverageRating,
		&gym.TotalFeedback,
		&gym.TrainersCount,
		&gym.StaffCount,
		&gym.ManagersCount,
		&gym.ActivePlans,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		LogService.LogError("❌ DB error: ", err)
		return nil, err
	}
	return gym, nil
}

func GetGymOwnerAndManagers(gymId int) ([]int, error) {

	query := `select u.userId
		from users u 
		inner join employees e on e.userId = u.userId
		where e.gymId = @p1
		and u.roleId = 2
		union 
		select u.userId
		from gym g 
		inner join users u on u.userId = g.createdBy
		where g.gymId = @p1`
	rows, err := db.DB.Query(query, gymId)
	if err != nil {
		LogService.LogError("❌ DB query error: ", err)
		return nil, err
	}
	defer rows.Close()

	var allUsers []int

	for rows.Next() {
		var currentUser int
		err := rows.Scan(
			&currentUser,
		)
		if err != nil {
			LogService.LogError("❌ DB scan error: ", err)
			return nil, err
		}
		allUsers = append(allUsers, currentUser)
	}

	if err = rows.Err(); err != nil {
		LogService.LogError("❌ DB rows error: ", err)
		return nil, err
	}

	if len(allUsers) == 0 {
		return []int{}, nil
	}

	return allUsers, nil
}

// GetGymMembersWithDetails returns all gym members with their active plans and feedback ratings
func GetGymMembersWithDetails(gymId int) ([]models.GymMemberWithDetails, error) {
	rows, err := db.DB.Query(`
		SELECT 
			u.userId,
			u.userName,
			u.mobile,
			u.email,
			u.gender,
			u.profilePic,
			um.membershipId,
			p.planId,
			p.planName,
			p.planPrice,
			p.planDuration,
			um.startDate,
			um.endDate,
			um.membershipStatus,
			feedback.feedbackId,
			feedback.rating,
			feedback.comments,
			feedback.createdAt as feedbackDate
		FROM userMemberships um
		INNER JOIN users u ON u.userId = um.userId
		INNER JOIN plans p ON p.planId = um.planId
		outer apply (select top 1 feedbackId, rating, comments, createdAt from feedback where createdBy = u.userId AND gymId = p.gymId order by feedback.createdAt desc) feedback
		WHERE p.gymId = @p1
		AND um.isActive = 1
		AND um.membershipStatus = 'A'
		AND um.startDate <= GETDATE()
		AND um.endDate >= GETDATE()
		AND u.roleId = 5
		ORDER BY u.userName ASC
	`, gymId)

	if err != nil {
		LogService.LogError("❌ DB query error fetching gym members: ", err)
		return nil, err
	}
	defer rows.Close()

	members := []models.GymMemberWithDetails{}

	for rows.Next() {
		member := models.GymMemberWithDetails{}
		err := rows.Scan(
			&member.UserId,
			&member.UserName,
			&member.Mobile,
			&member.Email,
			&member.Gender,
			&member.ProfilePic,
			&member.MembershipId,
			&member.PlanId,
			&member.PlanName,
			&member.PlanPrice,
			&member.PlanDuration,
			&member.StartDate,
			&member.EndDate,
			&member.MembershipStatus,
			&member.FeedbackId,
			&member.Rating,
			&member.Comments,
			&member.FeedbackDate,
		)

		if err != nil {
			LogService.LogError("❌ DB scan error fetching gym members: ", err)
			return nil, err
		}

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		LogService.LogError("❌ DB rows iteration error fetching gym members: ", err)
		return nil, err
	}

	return members, nil
}

// UpdateGym updates gym details (only gym information, not owner)
func UpdateGym(req gymModels.UpdateGymRequest) error {
	result, err := db.DB.Exec(`
		UPDATE gym
		SET 
			gymName = @p1,
			state = @p2,
			city = @p3,
			gymAddress = @p4,
			contactNo = @p5,
			officialEmail = @p6,
			locationLat = @p7,
			locationLong = @p8,
			updatedAt = GETDATE()
		WHERE gymId = @p9
	`, req.GymName, req.State, req.City, req.GymAddress, req.ContactNo, req.OfficialEmail, req.LocationLat, req.LocationLong, req.GymId)

	if err != nil {
		LogService.LogError("❌ DB error updating gym: ", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		LogService.LogError("❌ Error getting rows affected: ", err)
		return err
	}

	if rowsAffected == 0 {
		LogService.LogMessage("⚠️ No gym was updated (gymId might not exist)")
		return sql.ErrNoRows
	}

	LogService.LogMessage("✅ Gym updated successfully")
	return nil
}
