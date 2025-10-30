package userRepo

import (
	"database/sql"
	"zymm/internal/db"
	"zymm/internal/models"
	LogService "zymm/internal/service/log_service"
)

func GetUserDataByEmailOrMobile(emailOrMobile string) (*models.LoginUserDataModel, error) {
	var userData models.LoginUserDataModel
	err := db.DB.QueryRow(`
		SELECT userId,userName,userPass,roleId 
		FROM users 
		WHERE email = @p1 or mobile = @p2
		and isActive = 1
		`,
		emailOrMobile, emailOrMobile).Scan(&userData.UserId, &userData.DisplayName, &userData.Password, &userData.RoleId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		LogService.LogError("❌ DB error: ", err)
		return nil, err
	}
	return &userData, nil
}

// func GetUserByEmailOrMobile(emailOrMobile string) (*models.UserRecord, error) {
// 	user := &models.UserRecord{}
// 	err := db.DB.QueryRow(`
// 		SELECT userId, userName, gender, mobile, email, profilePic, roleId, createdAt, updatedAt
// 		FROM users
// 		WHERE email = @p1 OR mobile = @p2`,
// 		emailOrMobile, emailOrMobile).Scan(
// 		&user.UserId, &user.UserName, &user.Gender, &user.Mobile, &user.Email, &user.ProfilePic, &user.RoleId, &user.CreatedAt, &user.UpdatedAt)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, sql.ErrNoRows
// 		}
// 		LogService.LogError("❌ DB error: ", err)
// 		return nil, err
// 	}
// 	return user, nil
// }

func GetUserByUserId(userId int) (*models.UserRecord, error) {
	user := &models.UserRecord{}
	err := db.DB.QueryRow(`
		SELECT userId, userName, gender, mobile, email, profilePic, roleId, createdAt, updatedAt
		FROM users
		WHERE userId = @p1
		and isActive = 1
		`,
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

func CheckUserExistsByMobile(mobile string) (bool, error) {
	var exists int

	err := db.DB.QueryRow(`
		SELECT 
			CASE 
				WHEN EXISTS (SELECT 1 FROM users WHERE mobile = @p1 and isActive = 1)
					THEN 1 
				ELSE 0 
			END`,
		mobile).Scan(&exists)

	if err != nil {
		LogService.LogError("❌ DB error: ", err)
		return false, err
	}

	return exists == 1, nil
}

func CheckUserExistsByEmail(email string) (bool, error) {
	var exists int

	err := db.DB.QueryRow(`
		SELECT 
			CASE 
				WHEN EXISTS (SELECT 1 FROM users WHERE email = @p1 and isActive = 1)
					THEN 1 
				ELSE 0 
			END`,
		email).Scan(&exists)

	if err != nil {
		LogService.LogError("❌ DB error: ", err)
		return false, err
	}

	return exists == 1, nil
}

func UpdateUserProfilePicture(userId int, profilePicUrl string) error {
	_, err := db.DB.Exec(`
		UPDATE users 
		SET profilePic = @p1, updatedAt = GETDATE()
		WHERE userId = @p2 and isActive = 1
	`, profilePicUrl, userId)

	if err != nil {
		LogService.LogError("❌ DB error updating profile picture: ", err)
		return err
	}

	return nil
}

// DeactivateUser sets isActive to 0 for a user
func DeactivateUser(userId int) error {
	result, err := db.DB.Exec(`
		UPDATE users 
		SET isActive = 0, updatedAt = GETDATE()
		WHERE userId = @p1 and isActive = 1
	`, userId)

	if err != nil {
		LogService.LogError("❌ DB error deactivating user: ", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		LogService.LogError("❌ DB error getting rows affected: ", err)
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// ActivateUser sets isActive to 1 for a user
func ActivateUser(userId int) error {
	result, err := db.DB.Exec(`
		UPDATE users 
		SET isActive = 1, updatedAt = GETDATE()
		WHERE userId = @p1 and isActive = 0
	`, userId)

	if err != nil {
		LogService.LogError("❌ DB error activating user: ", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		LogService.LogError("❌ DB error getting rows affected: ", err)
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
