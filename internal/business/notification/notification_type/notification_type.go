package businessNotificationType

type Notification struct {
	Name        string
	Description string
}

const (
	NotificationGeneralInfo   = 1
	NotificationGymBroadcast  = 2
	NotificationMembershipReq = 3
	NotificationMembershipRes = 4
)

var Notifications = map[int]Notification{
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
}
