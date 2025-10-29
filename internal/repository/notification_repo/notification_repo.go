package notificationRepo

import (
	"database/sql"

	businessNotificationType "zymm/internal/business/notification/notification_type"
	"zymm/internal/db"
	"zymm/internal/models"
	LogService "zymm/internal/service/log_service"
)

func InsertNotification(userId int, notificationType businessNotificationType.NotificationType, title string, description string, createdBy int) error {
	_, err := db.DB.Exec(`
		INSERT INTO notifications (userId, notificationType, notificationTitle, notificationDesc, createdBy)
		VALUES (@p1, @p2, @p3, @p4, @p5)`,
		userId,
		notificationType,
		title,
		description,
		createdBy,
	)

	if err != nil {
		LogService.LogError("❌ DB error inserting notification: ", err)
		return err
	}

	return nil
}

func MarkNotificationAsRead(notificationId int, userId int) error {
	result, err := db.DB.Exec(`
		UPDATE notifications
		SET isRead = 1
		WHERE notificationId = @p1 AND userId = @p2`, notificationId, userId)
	if err != nil {
		LogService.LogError("❌ DB error updating notification read status: ", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		LogService.LogError("❌ DB rows affected error updating notification read status: ", err)
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func GetNotificationsByUserId(userId int) ([]models.NotificationRecord, error) {
	rows, err := db.DB.Query(`
		SELECT notificationId,
		       notificationType,
		       notificationTitle,
		       notificationDesc,
		       userId,
		       createdBy,
		       createdOn,
		       isRead
		FROM notifications
		WHERE userId = @p1
		ORDER BY createdOn DESC`, userId)
	if err != nil {
		LogService.LogError("❌ DB query error fetching notifications: ", err)
		return nil, err
	}
	defer rows.Close()

	var notifications []models.NotificationRecord
	for rows.Next() {
		var notif models.NotificationRecord
		var isRead sql.NullBool
		if scanErr := rows.Scan(
			&notif.NotificationId,
			&notif.NotificationType,
			&notif.NotificationTitle,
			&notif.NotificationDesc,
			&notif.UserId,
			&notif.CreatedBy,
			&notif.CreatedOn,
			&isRead,
		); scanErr != nil {
			LogService.LogError("❌ DB scan error fetching notifications: ", scanErr)
			return nil, scanErr
		}
		notif.IsRead = isRead.Valid && isRead.Bool
		notifications = append(notifications, notif)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		LogService.LogError("❌ DB rows iteration error fetching notifications: ", rowsErr)
		return nil, rowsErr
	}

	return notifications, nil
}
