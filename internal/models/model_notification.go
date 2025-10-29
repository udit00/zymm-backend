package models

import "time"

type CreateGeneralNotificationRequest struct {
	UserIds     string `json:"userIds"`
	Title       string `json:"title"`
	Description string `json:"description"`
	GymId       *int   `json:"gymId,omitempty"`
}

type MarkNotificationAsReadRequest struct {
	NotificationId int `json:"notificationId"`
}

type NotificationRecord struct {
	NotificationId   int       `json:"notificationId"`
	NotificationType int       `json:"notificationType"`
	NotificationTitle string    `json:"notificationTitle"`
	NotificationDesc string    `json:"notificationDesc"`
	UserId           int       `json:"userId"`
	CreatedBy        int       `json:"createdBy"`
	CreatedOn        time.Time `json:"createdOn"`
	IsRead           bool      `json:"isRead"`
}


