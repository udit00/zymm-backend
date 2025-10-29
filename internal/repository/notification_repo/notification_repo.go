package notificationRepo

import (
	"database/sql"
	"fmt"

	bussinessMembershipRequest "zymm/internal/business/membership/membership_request_action_type"
	"zymm/internal/business/notification"
	businessNotificationType "zymm/internal/business/notification/notification_type"
	"zymm/internal/db"
	"zymm/internal/models"
	gymRepo "zymm/internal/repository/gym_repo"
	"zymm/internal/repository/userRepo"
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

func SendNotificationForMembershipRequested(gymId int, userId int) {
	userDetails, userDetailsErr := userRepo.GetUserByUserId(userId)
	if userDetailsErr != nil {
		fmt.Println("Non fatal error, couldn't get user details for notification.")
	}
	allUsersToSendNotification, allUsersToSendNotificationErr := gymRepo.GetGymOwnerAndManagers(gymId)
	if allUsersToSendNotificationErr != nil {
		fmt.Println("Non fatal error, couldn't get owner and manager userId's.")
	}
	notificationTitle, notificationDescription := notification.GetNotificationTitleAndDesc(businessNotificationType.NotificationMembershipReq, userDetails.UserName)
	for i := 0; i < len(allUsersToSendNotification); i++ {
		currentUserNotificationFor := allUsersToSendNotification[i]
		InsertNotification(currentUserNotificationFor, businessNotificationType.NotificationMembershipReq, notificationTitle, notificationDescription, userId)
	}

}

func SendNotificationForMembershipResponseTakenByGymOwnerManagers(actionTakenBy int, notificationFor int, planDetail models.PlanRecord, responseTaken bussinessMembershipRequest.ActionType) {
	notificationTitle := ""
	notificationDescription := ""
	if responseTaken == bussinessMembershipRequest.Approve {
		notificationTitle = "Acceptance!!!"
		notificationDescription = "Your request for " + planDetail.PlanName + " has been Approved. You can start today."
	} else if responseTaken == bussinessMembershipRequest.Reject {
		notificationTitle = "Rejection!!!"
		notificationDescription = "Your request for " + planDetail.PlanName + " has been Declined. Select another plan."
	} else {
		return
	}

	InsertNotification(notificationFor, businessNotificationType.NotificationMembershipRes, notificationTitle, notificationDescription, actionTakenBy)
}

func NotifyGymOwnerForMemberJoined(actionTakenBy int, notificationFor int, planDetail models.PlanRecord, responseTaken bussinessMembershipRequest.ActionType) {
	actionTakenByDetails, actionTakenByDetailsErr := userRepo.GetUserByUserId(actionTakenBy)
	if actionTakenByDetailsErr != nil {
		fmt.Println("Non fatal error, couldn't get manager's user details.")
		return
	}
	notificationForDetails, notificationForDetailsErr := userRepo.GetUserByUserId(notificationFor)
	if notificationForDetailsErr != nil {
		fmt.Println("Non fatal error, couldn't get user's details.")
		return
	}

	notificationTitle := ""
	notificationDescription := ""
	if responseTaken == bussinessMembershipRequest.Approve {
		notificationTitle = "New Member!"
		notificationDescription = "A new member " + notificationForDetails.UserName + " has joined for " + planDetail.PlanName + ". The request was Approved by " + actionTakenByDetails.UserName + "."
	} else if responseTaken == bussinessMembershipRequest.Reject {
		notificationTitle = "Rejection!!!"
		notificationDescription = "A member " + notificationForDetails.UserName + " was rejected for " + planDetail.PlanName + ". The request was Rejected by " + actionTakenByDetails.UserName + "."
	} else {
		return
	}

	InsertNotification(notificationFor, businessNotificationType.NotificationMembershipRes, notificationTitle, notificationDescription, actionTakenBy)
}
