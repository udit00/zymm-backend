package authRepo

import (
	bussinessAuth "zymm/internal/business/auth"
	"zymm/internal/db"
	"zymm/internal/models"
	LogService "zymm/internal/service/log_service"
)

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

func UpdateUserPassword(userId int, newPass string) error {
	hashedPass, pErr := bussinessAuth.HashPasswordArgon2id(newPass)
	if pErr != nil {
		LogService.LogError("Error generating password hash: ", pErr)
		return pErr
	}
	_, err := db.DB.Exec(`UPDATE users SET userPass = @p1 WHERE userId = @p2`, hashedPass, userId)
	if err != nil {
		LogService.LogError("❌ DB error: ", err)
		return err
	}
	return nil
}

func InActiveUser(userId int) error {
	_, err := db.DB.Exec(`UPDATE users SET isActive = 0 WHERE userId = @p1`, userId)
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
