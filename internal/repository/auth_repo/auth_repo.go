package authRepo

import (
	"database/sql"
	bussinessAuth "zymm/internal/business/auth"
	"zymm/internal/db"
	"zymm/internal/models"
	gymModels "zymm/internal/models/gym_models"
	LogService "zymm/internal/service/log_service"
)

func GetUserDataByEmailOrMobile(emailOrMobile string) (*models.LoginUserDataModel, error) {
	var userData models.LoginUserDataModel
	err := db.DB.QueryRow("SELECT userId,userName,userPass FROM users WHERE email = @p1 or mobile = @p2", emailOrMobile, emailOrMobile).Scan(&userData.UserId, &userData.DisplayName, &userData.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		LogService.LogError("❌ DB error: ", err)
		return nil, err
	}
	return &userData, nil
}

func InsertLoginLog(loginLogsModel models.LoginLogsRecord) error {
	_, err := db.DB.Exec(`INSERT INTO loginLogs (userId, appVersion, userAgent, locationLat, locationLong, ipAddress) VALUES (@p1, @p2, @p3, @p4, @p5, @p6)`, loginLogsModel.UserId, loginLogsModel.AppVersion, loginLogsModel.UserAgent, loginLogsModel.LocationLat, loginLogsModel.LocationLong, loginLogsModel.IpAddress)
	if err != nil {
		LogService.LogError("❌ DB error: ", err)
		return err
	}
	return nil
}

func UpdateLoginAuthToken(userId int, authToken string) error {
	_, err := db.DB.Exec(`UPDATE users SET authToken = @p1 WHERE userId = @p2`, authToken, userId)
	if err != nil {
		LogService.LogError("❌ DB error: ", err)
		return err
	}
	return nil
}

func InsertRegistrationRecord(record models.UserRecord) (*models.UserRecord, error) {
	var userId int
	// Hash password
	hashedPass, pErr := bussinessAuth.HashPasswordArgon2id(record.UserPass)
	if pErr != nil {
		LogService.LogError("Error generating password hash: ", pErr)
		return nil, pErr
	}

	err := db.DB.QueryRow(`
		INSERT INTO users (userName, userPass, gender, mobile, email, profilePic, roleId)
		OUTPUT INSERTED.userId
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)`,
		record.UserName, hashedPass, record.Gender, record.Mobile, record.Email, record.ProfilePic, record.RoleId).Scan(&userId)
	if err != nil {
		LogService.LogError("❌ DB error: ", err)
		return nil, err
	}
	record.UserId = userId
	return &record, nil
}

func InsertRegistrationLog(record models.RegistrationLogsRecord) error {
	err := db.DB.QueryRow(`
		INSERT INTO registrationLogs (userId, userName, userPass, gender, mobile, email, profilePic, roleId, appVersion, userAgent)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10)`,
		record.UserId, record.UserName, record.UserPass, record.Gender, record.Mobile, record.Email, record.ProfilePic, record.RoleId, record.AppVersion, record.UserAgent).Err()
	if err != nil {
		LogService.LogError("❌ DB error: ", err)
		return err
	}
	return nil
}

func GetUserByEmailOrMobile(emailOrMobile string) (*models.UserRecord, error) {
	user := &models.UserRecord{}
	err := db.DB.QueryRow(`
		SELECT userId, userName, gender, mobile, email, profilePic, roleId, createdAt, updatedAt
		FROM users
		WHERE email = @p1 OR mobile = @p2`,
		emailOrMobile, emailOrMobile).Scan(
		&user.UserId, &user.UserName, &user.Gender, &user.Mobile, &user.Email, &user.ProfilePic, &user.RoleId, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		LogService.LogError("❌ DB error: ", err)
		return nil, err
	}
	return user, nil
}

func GetUserByUserId(userId int) (*models.UserRecord, error) {
	user := &models.UserRecord{}
	err := db.DB.QueryRow(`
		SELECT userId, userName, gender, mobile, email, profilePic, roleId, createdAt, updatedAt
		FROM users
		WHERE userId = @p1`,
		userId).Scan(
		&user.UserId, &user.UserName, &user.Gender, &user.Mobile, &user.Email, &user.ProfilePic, &user.RoleId, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		LogService.LogError("❌ DB error: ", err)
		return nil, err
	}
	return user, nil
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
