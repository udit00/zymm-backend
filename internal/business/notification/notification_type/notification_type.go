package businessNotificationType

type notification struct {
	Name        string
	Description string
}

type NotificationType int

const (
	NotificationZymm             NotificationType = 1
	NotificationGeneralInfo      NotificationType = 2
	NotificationGymBroadcast     NotificationType = 3
	NotificationMembershipReq    NotificationType = 4
	NotificationMembershipRes    NotificationType = 5
	NotificationFeedbackReceived NotificationType = 6
)

var notificationsData = map[NotificationType]notification{
	NotificationZymm: {
		Name:        "Zymm level broadcast.",
		Description: "Reserved for the app holder to notify all zymm owners/managers.",
	},
	NotificationGeneralInfo: {
		Name:        "general_info",
		Description: "Normal general information notifications.",
	},
	NotificationGymBroadcast: {
		Name:        "gym_broadcast",
		Description: "Gym Level Broadcast (use cases: holiday notifications, new gym trainer notifications, time schedule changes).",
	},
	NotificationMembershipReq: {
		Name:        "membership_request",
		Description: "Membership Request (use cases: from a member to an owner/gym manager).",
	},
	NotificationMembershipRes: {
		Name:        "membership_response",
		Description: "Membership Accept/Reject Notification.",
	},
	NotificationFeedbackReceived: {
		Name:        "feedback_received",
		Description: "A member has submitted a feedback.",
	},
}
