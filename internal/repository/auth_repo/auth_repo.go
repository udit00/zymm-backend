package authRepo

import (
	"database/sql"
	"log"
	bussinessAuth "zymm/internal/business/auth"
	"zymm/internal/db"
	"zymm/internal/models"
)

func InsertRegistrationRecord(record models.UserRecord) (*models.UserRecord, error) {
	var userId int
	// Hash password
	hashedPass, pErr := bussinessAuth.HashPasswordArgon2id(record.UserPass)
	if pErr != nil {
		log.Println("Error generating password hash: ", pErr)
		return nil, pErr
	}

	err := db.DB.QueryRow(`
		INSERT INTO users (userName, userPass, gender, mobile, email, profilePic, roleId)
		OUTPUT INSERTED.userId
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)`,
		record.UserName, hashedPass, record.Gender, record.Mobile, record.Email, record.ProfilePic, 5).Scan(&userId)
	if err != nil {
		log.Printf("❌ DB error: %v", err)
		return nil, err
	}
	record.UserId = userId
	return &record, nil
}

func InsertRegistrationLog(record models.RegistrationLogsRecord) error {
	err := db.DB.QueryRow(`
		INSERT INTO registrationLogs (userId, userName, userPass, gender, mobile, email, profilePic, roleId, appVersion, appPlatform)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10)`,
		record.UserId, record.UserName, record.UserPass, record.Gender, record.Mobile, record.Email, record.ProfilePic, record.RoleId, record.AppVersion, record.AppPlatform).Err()
	if err != nil {
		log.Printf("❌ DB error: %v", err)
		return err
	}
	return nil
}

func GetUserByEmailOrMobile(emailOrMobile string) (*models.UserRecord, error) {
	user := &models.UserRecord{}
	err := db.DB.QueryRow(`
		SELECT userId, userName, userPass, gender, mobile, email, profilePic, roleId, createdAt, updatedAt
		FROM users
		WHERE email = @p1 OR mobile = @p2`,
		emailOrMobile, emailOrMobile).Scan(
		&user.UserId, &user.UserName, &user.UserPass, &user.Gender, &user.Mobile, &user.Email, &user.ProfilePic, &user.RoleId, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		log.Printf("❌ DB error: %v", err)
		return nil, err
	}
	return user, nil
}
