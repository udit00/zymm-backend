package notification

import businessNotificationType "zymm/internal/business/notification/notification_type"

func GetNotificationTitleAndDesc(notificationType businessNotificationType.NotificationType, data string) (string, string) {
	switch notificationType {
	case businessNotificationType.NotificationZymm:
		return GetTitleAndDescForZymmLevelNotification(data)
	case businessNotificationType.NotificationGeneralInfo:
		return GetTitleAndDescForGeneralNotification(data)
	case businessNotificationType.NotificationGymBroadcast:
		return GetTitleAndDescForGymNotification(data)
	case businessNotificationType.NotificationMembershipReq:
		return GetTitleAndDescForMembershipRequest(data)
	case businessNotificationType.NotificationMembershipRes:
		return GetTitleAndDescForMembershipResponse(data)
	case businessNotificationType.NotificationFeedbackReceived:
		return GetTitleAndDescForFeedbackReceived()
	case businessNotificationType.NotificationPendingFees:
		return GetTitleAndDescForPendingFees(data)
	default:
		return "", ""
	}

}

func GetTitleAndDescForZymmLevelNotification(data string) (string, string) {
	return "Attention Gym Owners/ Managers", data
}

func GetTitleAndDescForGeneralNotification(data string) (string, string) {
	return "Attention!!!", data
}

func GetTitleAndDescForGymNotification(data string) (string, string) {
	return "Attention Members!!!", data
}

func GetTitleAndDescForMembershipRequest(userName string) (string, string) {
	return "Membership Request!", "High Priority: You have a new membership request from " + userName + ", take action now."
}

func GetTitleAndDescForMembershipResponse(gym string) (string, string) {
	return "Membership Response!", gym + " management has taken action on your membership request."
}

func GetTitleAndDescForFeedbackReceived() (string, string) {
	return "Feedback!", "A new feedback was submitted by a user."
}

func GetTitleAndDescForPendingFees(planName string) (string, string) {
	// Keep message short to fit in 100 char DB limit
	return "Fee Reminder!", "Your " + planName + " membership needs renewal. Please pay to continue."
}
