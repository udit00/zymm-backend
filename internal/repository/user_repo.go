package repository

import (
	"log"
	"zymm/internal/db"
	"zymm/internal/models"
)

// GetAllUsers → fetch all users from DB
func GetAllUsers() ([]models.User, error) {
	rows, err := db.DB.Query("SELECT user_id, user_name, user_pass FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Pass); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		log.Println("Row iteration error:", err)
		return nil, err
	}

	return users, nil
}
